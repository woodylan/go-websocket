package routers

import (
	"go-websocket/api/bindgroup"
	"go-websocket/api/connect"
	"go-websocket/api/push2client"
	"go-websocket/api/push2group"
	"go-websocket/servers"
	"net/http"
)

func Init() {
	websocketHandler := &connect.Controller{}
	pushToClientHandler := &push2client.Controller{}
	pushToGroupHandler := &push2group.Controller{}
	bindToGroupHandler := &bindgroup.Controller{}

	http.HandleFunc("/ws", websocketHandler.Run)
	http.HandleFunc("/push_to_client", pushToClientHandler.Run)
	http.HandleFunc("/push_to_group", pushToGroupHandler.Run)
	http.HandleFunc("/bind_to_group", bindToGroupHandler.Run)

	go servers.WriteMessage()
}
