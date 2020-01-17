package routers

import (
	"go-websocket/api/bindgroup"
	"go-websocket/api/send2client"
	"go-websocket/api/send2group"
	"go-websocket/servers"
	"net/http"
)

func Init() {
	//websocketHandler := &connect.Controller{}
	pushToClientHandler := &send2client.Controller{}
	pushToGroupHandler := &send2group.Controller{}
	bindToGroupHandler := &bindgroup.Controller{}

	//http.HandleFunc("/ws", websocketHandler.Run)
	http.HandleFunc("/api/send_to_client", pushToClientHandler.Run)
	http.HandleFunc("/api/send_to_group", pushToGroupHandler.Run)
	http.HandleFunc("/api/bind_to_group", bindToGroupHandler.Run)

	servers.StartWebSocket()

	go servers.WriteMessage()
}
