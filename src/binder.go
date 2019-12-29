package src

import (
	"encoding/json"
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
	//添加客户端到分组
	AddClient2Group(groupName, clientId string)
	//获取分组客户端列表
	GetGroupClientList(groupName string) ([]string)
	//发送信息到指定客户端
	SendMessage2Client(clientId, message string)
	//发送到本机分组
	SendMessage2LocalGroup(groupName, message string)
	//发送信息到指定分组
	SendMessage2Group(groupName, message string)
	//发送到RabbitMQ，方便同步到其他机器
	Send2RabbitMQ(msgType int, objectId, message string)
}

func NewBinder() *binder {
	return &binder{
		clintId2ConnMap: make(map[string]*Conn),
		clientGroupsMap: make(map[string][]string, 0),
		groupClientIds:  make(map[string][]string, 0),
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
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clintId2ConnMap, clientId)
	delete(b.clientGroupsMap, clientId)
}

//客户端数量
func (b *binder) ClientNumber() int {
	return len(b.clintId2ConnMap)
}

//添加客户端到分组
func (b *binder) AddClient2Group(groupName, clientId string) {
	//如果是集群则用redis共享数据
	if util.IsCluster() {
		_, err := redis.SetAdd(define.REDIS_KEY_GROUP+groupName, clientId)
		if err != nil {
			panic(err)
		}
	} else {
		//如果是单机，就没比用redis了
		b.groupClientIds[groupName] = append(b.groupClientIds[groupName], clientId)
	}
}

//获取分组客户端列表
func (b *binder) GetGroupClientList(groupName string) ([]string) {
	if util.IsCluster() {
		groupList, err := redis.SetMembers(define.REDIS_KEY_GROUP + groupName)
		if err != nil {
			panic(err)
		}
		return groupList
	}

	return b.groupClientIds[groupName]
}

//发送信息到指定客户端
func (b *binder) SendMessage2Client(clientId, message string) {
	if util.IsCluster() {
		//发送到RabbitMQ
		b.Send2RabbitMQ(define.MESSAGE_TYPE_CLIENT, clientId, message)
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
				SendMessage2LocalClient(clientId, message)
			}
		}
	}
}

//发送信息到指定分组
func (b *binder) SendMessage2Group(groupName, message string) {
	if util.IsCluster() {
		//发送到RabbitMQ
		b.Send2RabbitMQ(define.MESSAGE_TYPE_GROUP, groupName, message)
	} else {
		//如果是单机服务，则只发送到本机
		b.SendMessage2LocalGroup(groupName, message)
	}
}

//发送到RabbitMQ，方便同步到其他机器
func (b *binder) Send2RabbitMQ(msgType int, objectId, message string) {
	if rabbitMQ == nil {
		panic("rabbitMQ连接失败")
	}

	publishMessage := publishMessage{
		MsgType:  msgType,
		ObjectId: objectId,
		Message:  message,
	}

	messageByte, _ := json.Marshal(publishMessage)

	rabbitMQ.PublishPub(string(messageByte))
}
