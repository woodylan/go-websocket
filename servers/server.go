package servers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"go-websocket/api"
	"go-websocket/define"
	"go-websocket/pkg/redis"
	"go-websocket/tools/util"
	"log"
	"net/http"
	"time"
)

//channel通道
var ToClientChan chan clientInfo

//channel通道结构体
type clientInfo struct {
	ClientId *string
	Code     int
	Msg      string
	Data     *interface{}
}

// 心跳间隔
var heartbeatInterval = 25 * time.Second

type publishMessage struct {
	Type      int         `json:"type"`
	SystemId  string      `json:"systemId"`
	GroupName string      `json:"groupName"`
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
}

func init() {
	ToClientChan = make(chan clientInfo, 1000)
}

var Manager = NewClientManager() // 管理者

func StartWebSocket() {
	websocketHandler := &Controller{}
	http.HandleFunc("/ws", websocketHandler.Run)

	go Manager.Start()
}

//发送信息到指定客户端
func SendMessage2Client(clientId *string, code int, msg string, data *interface{}) {
	if util.IsCluster() {
		addr, _, _, isLocal, err := util.GetAddrInfoAndIsLocal(*clientId)
		if err != nil {
			_ = fmt.Errorf("%s", err)
			return
		}

		//如果是本机则发送到本机
		if isLocal {
			SendMessage2LocalClient(clientId, code, msg, data)
		} else {
			//发送到指定机器
			SendRpc2Client(addr, clientId, msg, data)
		}
	} else {
		//如果是单机服务，则只发送到本机
		SendMessage2LocalClient(clientId, code, msg, data)
	}
}

//添加客户端到分组
func AddClient2Group(systemId *string, groupName *string, clientId string) {
	//如果是集群则用redis共享数据
	if util.IsCluster() {
		//判断key是否存在
		addr, _, _, isLocal, err := util.GetAddrInfoAndIsLocal(clientId)
		if err != nil {
			_ = fmt.Errorf("%s", err)
			return
		}

		if isLocal {
			if client, err := Manager.GetByClientId(clientId); err == nil {
				//添加到本地
				Manager.AddClient2LocalGroup(util.GenGroupKey(*systemId, *groupName), client)
			} else {
				fmt.Println(err)
			}
		} else {
			//发送到指定的机器
			SendRpcBindGroup(&addr, systemId, groupName, &clientId)
		}
	} else {
		if client, err := Manager.GetByClientId(clientId); err == nil {
			//如果是单机，就直接添加到本地group了
			Manager.AddClient2LocalGroup(util.GenGroupKey(*systemId, *groupName), client)
		};
	}
}

//发送信息到指定分组
func SendMessage2Group(systemId, groupName *string, code int, msg string, data *interface{}) {
	if util.IsCluster() {
		//发送到RabbitMQ
		_ = SendGroupMessage2RabbitMQ(systemId, groupName, code, msg, data)
	} else {
		//如果是单机服务，则只发送到本机
		Manager.SendMessage2LocalGroup(systemId, groupName, code, msg, data)
	}
}

//发送信息到指定系统
func SendMessage2System(systemId *string, code int, msg string, data interface{}) {
	if util.IsCluster() {
		//发送到RabbitMQ
		_ = SendSystemMessage2RabbitMQ(systemId, code, msg, &data)
	} else {
		//如果是单机服务，则只发送到本机
		Manager.SendMessage2LocalSystem(systemId, code, msg, &data)
	}
}

//发送到RabbitMQ，方便同步到其他机器
func SendGroupMessage2RabbitMQ(systemId, GroupName *string, code int, msg string, data *interface{}) error {
	publishMessage := publishMessage{
		Type:      define.RPC_MESSAGE_TYPE_GROUP,
		SystemId:  *systemId,
		GroupName: *GroupName,
		Code:      code,
		Msg:       msg,
		Data:      data,
	}

	return send2RabbitMQ(&publishMessage)
}

//发送系统消息到RabbitMQ
func SendSystemMessage2RabbitMQ(systemId *string, code int, msg string, data *interface{}) error {
	publishMessage := publishMessage{
		Type:      define.RPC_MESSAGE_TYPE_SYSTEM,
		SystemId:  *systemId,
		GroupName: "",
		Code:      code,
		Msg:       msg,
		Data:      data,
	}

	return send2RabbitMQ(&publishMessage)
}

func send2RabbitMQ(publishMessage *publishMessage) error {
	if rabbitMQ == nil {
		log.Fatal("rabbitMQ连接失败")
		return errors.New("rabbitMQ连接失败")
	}

	messageByte, _ := json.Marshal(publishMessage)

	rabbitMQ.PublishPub(string(messageByte))
	return nil
}

//通过本服务器发送信息
func SendMessage2LocalClient(clientId *string, code int, msg string, data *interface{}) {
	ToClientChan <- clientInfo{ClientId: clientId, Code: code, Msg: msg, Data: data}
}

//监听并发送给客户端信息
func WriteMessage() {
	for {
		select {
		case clientInfo := <-ToClientChan:
			fmt.Println("发送到本机客户端：" + *clientInfo.ClientId)
			if conn, err := Manager.GetByClientId(*clientInfo.ClientId); err == nil && conn != nil {
				if err := Render(conn.Socket, clientInfo.Code, clientInfo.Msg, clientInfo.Data); err != nil {
					_ = conn.Socket.Close()
					log.Println(err)
					return
				} else {
					//延长key过期时间
					_, err := redis.SetSurvivalTime(define.REDIS_CLIENT_ID_PREFIX+*clientInfo.ClientId, define.REDIS_KEY_SURVIVAL_SECONDS)
					if (err != nil) {
						log.Println(err)
					}
				}
			}
		}
	}
}

func Render(conn *websocket.Conn, code int, message string, data interface{}) error {
	return conn.WriteJSON(api.RetData{
		Code: code,
		Msg:  message,
		Data: data,
	})
}

//启动定时器进行心跳检测
func PingTimer() {
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				//发送心跳
				for clientId, conn := range Manager.AllClient() {
					if err := conn.Socket.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second)); err != nil {
						_ = conn.Socket.Close()
						Manager.DelClient(conn)
						log.Printf("发送心跳失败: %s 总连接数：%d", clientId, Manager.Count())
					}
				}
			}
		}
	}()
}
