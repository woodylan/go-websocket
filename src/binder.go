package src

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go-websocket/define"
	"go-websocket/pkg/redis"
	"go-websocket/tools/util"
	"sync"
)

type binder struct {
	mu sync.RWMutex

	clintId2ConnMap map[string]*Conn
	clientGroupsMap map[string][]string
	groupClientIds  map[string][]string
}

type publishMessage struct {
	MsgType  int    `json:"type"`     //消息类型 1.指定客户端 2.指定分组
	ObjectId string `json:"objectId"` //对象ID，如果是client为clientId，如果是分组则为groupId
	Message  string `json:"message"`  //消息内容
}

type IBinder interface {
	//添加客户端
	AddClient(clientId string, conn *websocket.Conn)
	//删除客户端
	DelClient(clientId string)
	//客户端数量
	ClientNumber() int
	//客户端是否存在
	IsAlive(clientId string) (conn *Conn, ok bool)
	//添加客户端到分组
	AddClient2Group(groupName, clientId string)
	//获取分组客户端列表
	GetGroupClientList(groupName string) ([]string)
	//发送到本机分组
	SendMessage2LocalGroup(groupName, message string)
	//发送信息到指定分组
	SendMessage2Group(groupName, message string)
	//发送到RabbitMQ，方便同步到其他机器
	Send2RabbitMQ(objectId, message string)
}

func NewBinder() *binder {
	define.ClientGroupsMap = make(map[string][]string, 0);
	return &binder{
		clintId2ConnMap: make(map[string]*Conn),
		//clientGroupsMap: make(map[string][]string, 0),
		groupClientIds: make(map[string][]string, 0),
	}
}

//给客户端绑定ID
func (b *binder) AddClient(clientId string, conn *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clintId2ConnMap[clientId] = &Conn{Conn: conn}
}

//删除客户端
func (b *binder) DelClient(clientId string) {
	delete(b.clintId2ConnMap, clientId)
	if util.IsCluster() {
		//todo 删除redis
		//todo 删除集群里的分组信息
	} else {
		//删除单机里的分组
		define.ClientGroupsMapMu.Lock()
		defer define.ClientGroupsMapMu.Unlock()
		delete(define.ClientGroupsMap, clientId)
	}

}

//客户端数量
func (b *binder) ClientNumber() int {
	return len(b.clintId2ConnMap)
}

//客户端是否存在
func (b *binder) IsAlive(clientId string) (conn *Conn, ok bool) {
	conn, ok = b.clintId2ConnMap[clientId];
	return
}

//添加分组到本地
func AddClient2LocalGroup(groupName, clientId string) {
	if util.IsCluster() {
		_, err := redis.SetAdd(util.GetGroupKey(groupName), clientId)
		if err != nil {
			panic(err)
		}
	} else {
		define.ClientGroupsMap[clientId] = append(define.ClientGroupsMap[clientId], groupName)
	}
	//todo
	//b.groupClientIds[groupName] = append(b.groupClientIds[groupName], clientId)
}

//添加客户端到分组
func (b *binder) AddClient2Group(groupName, clientId string) {

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
			if _, isAlive := b.IsAlive(clientId); !isAlive {
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
func (b *binder) GetGroupClientList(groupName string) ([]string) {
	if util.IsCluster() {
		groupList, err := redis.SMEMBERS(util.GetGroupKey(groupName))
		if err != nil {
			panic(err)
		}
		return groupList
	}

	return b.groupClientIds[groupName]
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
func (b *binder) SendMessage2LocalGroup(groupName, message string) {
	if len(groupName) > 0 {
		clientList := b.GetGroupClientList(groupName)
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
func (b *binder) SendMessage2Group(groupName, message string) {
	if util.IsCluster() {
		//发送到RabbitMQ
		b.Send2RabbitMQ(groupName, message)
	} else {
		//如果是单机服务，则只发送到本机
		SendMessage2LocalClient(groupName, message)
	}
}

//发送到RabbitMQ，方便同步到其他机器
func (b *binder) Send2RabbitMQ(objectId, message string) {
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
