package src

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go-websocket/define"
	"go-websocket/pkg/redis"
	"go-websocket/tools/readconfig"
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

//给客户端绑定ID
func (b *binder) BindToMap(clientId string, conn *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clintId2ConnMap[clientId] = &Conn{Conn: conn}
}

//删除客户端
func (b *binder) DelMap(clientId string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clintId2ConnMap, clientId)
	delete(b.clientGroupsMap, clientId)
}

//客户端数量
func (b *binder) ClientNumber() int {
	return len(b.clintId2ConnMap)
}

//是否集群
func isCluster() bool {
	cluster, _ := readconfig.ConfigData.Bool("common::cluster")

	return cluster
}

//发送信息到指定客户端
func (b *binder) SendMessage2Client(clientId, message string) {
	if isCluster() {
		//发送到RabbitMQ
		Send2RabbitMQ(define.MESSAGE_TYPE_CLIENT, clientId, message)
	} else {
		//如果是单机服务，则只发送到本机
		SendMessage2Client(clientId, message)
	}
}

//发送信息到指定分组
func (b *binder) SendMessage2Group(groupName, message string) {
	if isCluster() {
		//发送到RabbitMQ
		Send2RabbitMQ(define.MESSAGE_TYPE_GROUP, groupName, message)
	} else {
		//如果是单机服务，则只发送到本机
		b.SendMessage2LocalGroup(groupName, message)
	}
}

func SendMessage2Client(clientId, message string) {
	toClientChan <- [2]string{clientId, message}
	fmt.Println(clientId, message)
}

//发送到RabbitMQ
func Send2RabbitMQ(msgType int, objectId, message string) {
	if rabbitMQ == nil {
		initRabbitMQ()
	}

	publishMessage := publishMessage{
		MsgType:  msgType,
		ObjectId: objectId,
		Message:  message,
	}

	messageByte, _ := json.Marshal(publishMessage)

	rabbitMQ.PublishPub(string(messageByte))
}

//发送到本机分组
func (b *binder) SendMessage2LocalGroup(groupName, message string) {
	if len(groupName) > 0 {
		clientList := b.GetGroupList(groupName)
		if len(clientList) > 0 {
			for _, clientId := range clientList {
				//发送信息
				//todo key的销毁
				SendMessage2Client(clientId, message)
			}
		}
	}
}

func (b *binder) SetGroupList(groupName, clientId string) {
	//如果是集群则用redis共享数据
	if isCluster() {
		_, err := redis.SetAdd(define.REDIS_KEY_GROUP+groupName, clientId)
		if err != nil {
			panic(err)
		}
	} else {
		//如果是单机，就没比用redis了
		b.groupClientIds[groupName] = append(b.groupClientIds[groupName], clientId)
	}
}

func (b *binder) GetGroupList(groupName string) ([]string) {
	if isCluster() {
		groupList, err := redis.SetMembers(define.REDIS_KEY_GROUP + groupName)
		if err != nil {
			panic(err)
		}
		return groupList
	}

	return b.groupClientIds["groupName"]
}
