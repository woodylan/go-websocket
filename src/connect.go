package src

import (
	"github.com/gorilla/websocket"
	"go-websocket/tools/util"
	"net/http"
)

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
	binder := NewBinder()

	//如果是集群，则读取初始化RabbitMQ实例
	if util.IsCluster() {
		initRabbitMQ()
		initRabbitMQReceive(binder)
	}

	websocketHandler := &WebsocketHandler{
		defaultUpgrader,
		binder,
	}

	pushToClientHandler := &PushToClientHandler{
		binder: binder,
	}

	pushToGroupHandler := &PushToGroupHandler{
		binder: binder,
	}

	bindToGroupHandler := &BindToGroupHandler{
		binder: binder,
	}

	http.Handle(s.WSPath, websocketHandler)
	http.Handle(s.PushToClientPath, pushToClientHandler)
	http.Handle(s.PushToGroupPath, pushToGroupHandler)
	http.Handle(s.BindToGroupPath, bindToGroupHandler)

	go websocketHandler.WriteMessage()

	return http.ListenAndServe(s.Addr, nil)
}
