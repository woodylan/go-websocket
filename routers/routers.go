package routers

import (
	"github.com/gorilla/websocket"
	"go-websocket/api/bindgroup"
	"go-websocket/api/connect"
	"go-websocket/api/push2client"
	"go-websocket/api/push2group"
	"net/http"
)

var defaultUpgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Init() {
	websocketHandler := &connect.WebsocketHandler{
		defaultUpgrader,
	}

	pushToClientHandler := &push2client.Push2ClientHandler{}

	pushToGroupHandler := &push2group.Push2GroupHandler{}

	bindToGroupHandler := &bindgroup.BindGroupHandler{}

	http.Handle("/ws", websocketHandler)
	http.Handle("/push_to_client", pushToClientHandler)
	http.Handle("/push_to_group", pushToGroupHandler)
	http.Handle("/bind_to_group", bindToGroupHandler)

	go websocketHandler.WriteMessage()
}
