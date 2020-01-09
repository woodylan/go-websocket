package servers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go-websocket/clientvar"
	"go-websocket/pkg/redis"
	"go-websocket/tools/util"
	"log"
	"time"
)

// 心跳间隔
var heartbeatInterval = 10 * time.Second

//channel通道
var ToClientChan chan [2]string

func init() {
	ToClientChan = make(chan [2]string, 10)
}

//通过本服务器发送信息
func SendMessage2LocalClient(clientId, message string) {
	ToClientChan <- [2]string{clientId, message}
}

type publishMessage struct {
	MsgType  int    `json:"type"`     //消息类型 1.指定客户端 2.指定分组
	ObjectId string `json:"objectId"` //对象ID，如果是client为clientId，如果是分组则为groupId
	Message  string `json:"message"`  //消息内容SendRpcBindGroup
}

//添加分组到本地
func AddClient2LocalGroup(groupName, clientId string) {
	if util.IsCluster() {
		_, err := redis.SetAdd(util.GetGroupKey(groupName), clientId)
		if err != nil {
			panic(err)
		}
	} else {
		clientvar.ClientGroupsMap[clientId] = append(clientvar.ClientGroupsMap[clientId], groupName)
	}
	//todo
	//b.groupClientIds[groupName] = append(b.groupClientIds[groupName], clientId)
}

//添加客户端到分组
func AddClient2Group(groupName, clientId string) {

	//如果是集群则用redis共享数据
	if util.IsCluster() {
		//判断key是否存在
		addr, _, _, isLocal, err := util.GetAddrInfoAndIsLocal(clientId)
		if err != nil {
			_ = fmt.Errorf("%s", err)
			return
		}

		if isLocal {
			//判断是否已经存在
			if _, isAlive := clientvar.IsAlive(clientId); !isAlive {
				return
			}
			//添加到本地
			AddClient2LocalGroup(groupName, clientId)
		} else {
			//发送到指定的机器
			SendRpcBindGroup(addr, groupName, clientId)
		}
	} else {
		//如果是单机，就直接添加到本地group了
		AddClient2LocalGroup(groupName, clientId)
	}
}

//获取分组客户端列表
func GetGroupClientList(groupName string) ([]string) {
	if util.IsCluster() {
		groupList, err := redis.SMEMBERS(util.GetGroupKey(groupName))
		if err != nil {
			panic(err)
		}
		return groupList
	}

	return clientvar.GetGroupClientIds(groupName)
}

//发送信息到指定客户端
func SendMessage2Client(clientId, message string) {
	if util.IsCluster() {
		addr, _, _, isLocal, err := util.GetAddrInfoAndIsLocal(clientId)
		if err != nil {
			_ = fmt.Errorf("%s", err)
			return
		}

		//如果是本机则发送到本机
		if isLocal {
			go fmt.Println("发送到本机客户端：" + clientId + " 消息：" + message)
			SendMessage2LocalClient(clientId, message)
		} else {
			//发送到指定机器
			go fmt.Println("发送到服务器：" + addr + " 客户端：" + clientId + " 消息：" + message)
			SendRpc2Client(addr, clientId, message)
		}

		//发送到RabbitMQ
		//b.Send2RabbitMQ(define.MESSAGE_TYPE_CLIENT, clientId, message)
	} else {
		//如果是单机服务，则只发送到本机
		SendMessage2LocalClient(clientId, message)
	}
}

//发送到本机分组
func SendMessage2LocalGroup(groupName, message string) {
	if len(groupName) > 0 {
		clientList := GetGroupClientList(groupName)
		if len(clientList) > 0 {
			for _, clientId := range clientList {
				//发送信息
				//todo key的销毁
				SendMessage2Client(clientId, message)
			}
		}
	}
}

//发送信息到指定分组
func SendMessage2Group(groupName, message string) {
	if util.IsCluster() {
		//发送到RabbitMQ
		Send2RabbitMQ(groupName, message)
	} else {
		//如果是单机服务，则只发送到本机
		SendMessage2LocalClient(groupName, message)
	}
}

//发送到RabbitMQ，方便同步到其他机器
func Send2RabbitMQ(objectId, message string) {
	if rabbitMQ == nil {
		panic("rabbitMQ连接失败")
	}

	publishMessage := publishMessage{
		ObjectId: objectId,
		Message:  message,
	}

	messageByte, _ := json.Marshal(publishMessage)

	rabbitMQ.PublishPub(string(messageByte))
}

//发送心跳数据
func SendJump(conn *websocket.Conn) {
	go func() {
		for {
			time.Sleep(heartbeatInterval)
			if err := conn.WriteJSON("heartbeat"); err != nil {
				//todo 删除客户端
				fmt.Printf("删除客户端")
				return
			}
		}

	}()
}

func WriteMessage() {
	for {
		select {
		case clientInfo := <-ToClientChan:
			toConn, ok := clientvar.IsAlive(clientInfo[0]);
			if ok {
				err := toConn.WriteJSON(clientInfo[1]);
				if err != nil {
					go clientvar.DelClient(clientInfo[0])
					log.Println(err)
				} else {
					//todo 给redis续命
				}
			} else {
				go clientvar.DelClient(clientInfo[0])
			}
		}
	}
}
