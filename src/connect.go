package src

import (
	"github.com/gorilla/websocket"
	"net/http"
)

const (
	defaultWSPath   = "/ws"
	defaultPushPath = "/push"
)

type Server struct {
	Addr     string //监听地址
	WSPath   string //websocket路径，如'/ws'
	PushPath string //推送消息地址,如'/push'
	Upgrader *websocket.Upgrader
	wh       *WebsocketHandler
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
		Addr:     addr,
		WSPath:   defaultWSPath,
		PushPath: defaultPushPath,
	}
}

func (s *Server) ListenAndServer() error {
	b := &binder{
		clintId2ConnMap: make(map[string]*Conn),
	}

	websocketHandler := &WebsocketHandler{
		defaultUpgrader,
		b,
	}

	pushHandler := &PushHandler{
		binder: b,
	}

	http.Handle(s.WSPath, websocketHandler)
	http.Handle(s.PushPath,pushHandler)

	go websocketHandler.WriteMessage()

	return http.ListenAndServe(s.Addr, nil)
}
