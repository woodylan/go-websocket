package routers

import (
	"go-websocket/api/bind2group"
	"go-websocket/api/send2client"
	"go-websocket/api/send2group"
	"go-websocket/servers"
	"net/http"
)

func Init() {
	pushToClientHandler := &send2client.Controller{}
	pushToGroupHandler := &send2group.Controller{}
	bindToGroupHandler := &bind2group.Controller{}

	http.HandleFunc("/api/send_to_client", pushToClientHandler.Run)
	http.HandleFunc("/api/send_to_group", pushToGroupHandler.Run)
	http.HandleFunc("/api/bind_to_group", bindToGroupHandler.Run)

	servers.StartWebSocket()

	go servers.WriteMessage()
}
