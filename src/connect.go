package src

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go-websocket/define"
	RabbitMQ "go-websocket/src/rabbitMq"
	"go-websocket/src/readConfig"
	"log"
	"net/http"
)

//RabbitMQ 实例
var rabbitMQ *RabbitMQ.RabbitMQ

const (
	defaultWSPath    = "/ws"
	pushToClientPath = "/push_to_client"
	pushToGroupPath  = "/push_to_group"
	bindToGroupPath  = "/bind_to_group"
)

type Server struct {
	Addr             string //监听地址
	WSPath           string //websocket路径，如'/ws'
	PushToClientPath string //推送消息到指定客户端地址,如'/push_to_client'
	PushToGroupPath  string //推送消息到指定分组地址,如'/push_to_group'
	BindToGroupPath  string //绑定到分组的地址，如'/bind_to_group'
	Upgrader         *websocket.Upgrader
	wh               *WebsocketHandler
}

type Conn struct {
	Conn *websocket.Conn
}

// http升级websocket协议的配置
var defaultUpgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewServer(addr string) *Server {
	return &Server{
		Addr:             addr,
		WSPath:           defaultWSPath,
		PushToClientPath: pushToClientPath,
		PushToGroupPath:  pushToGroupPath,
		BindToGroupPath:  bindToGroupPath,
	}
}

func (s *Server) ListenAndServer() error {
	b := &binder{
		clintId2ConnMap: make(map[string]*Conn),
		clientGroupsMap: make(map[string][]string, 0),
		groupClientIds:  make(map[string][]string, 0),
	}

	//如果是集群，则读取初始化RabbitMQ实例
	if isCluster() {
		initRabbitMQ()
		initRabbitMQReceive(b)
	}

	websocketHandler := &WebsocketHandler{
		defaultUpgrader,
		b,
	}

	pushToClientHandler := &PushToClientHandler{
		binder: b,
	}

	pushToGroupHandler := &PushToGroupHandler{
		binder: b,
	}

	bindToGroupHandler := &BindToGroupHandler{
		binder: b,
	}

	http.Handle(s.WSPath, websocketHandler)
	http.Handle(s.PushToClientPath, pushToClientHandler)
	http.Handle(s.PushToGroupPath, pushToGroupHandler)
	http.Handle(s.BindToGroupPath, bindToGroupHandler)

	go websocketHandler.WriteMessage()

	return http.ListenAndServe(s.Addr, nil)
}

//创建rabbitMQ实例
func initRabbitMQ() {
	rabbitMQ = RabbitMQ.NewRabbitMQPubSub(
		readConfig.ConfigData.String("rabbitMQ::amqpurl"),
		readConfig.ConfigData.String("rabbitMQ::exchange"))

}

func initRabbitMQReceive(b *binder) {
	msgs, err := rabbitMQ.ReceiveSub()
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for receiveData := range msgs {
			log.Printf("Received a message: %s", receiveData.Body)

			var publishMessage publishMessage
			err := json.Unmarshal([]byte(receiveData.Body), &publishMessage)
			if err == nil {
				if publishMessage.MsgType == define.MESSAGE_TYPE_CLIENT {
					//发送到指定客户端
					SendMessage2Client(publishMessage.ObjectId, publishMessage.Message)
				} else if publishMessage.MsgType == define.MESSAGE_TYPE_GROUP {
					//发送到指定分组
					b.SendMessage2LocalGroup(publishMessage.ObjectId, publishMessage.Message)
				}
			} else {
				fmt.Println(err)
			}
		}
	}()
}
