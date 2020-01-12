package routers

import (
	"go-websocket/api/bindgroup"
	"go-websocket/api/connect"
	"go-websocket/api/send2client"
	"go-websocket/api/send2group"
	"go-websocket/servers/server"
	"net/http"
)

func Init() {
	websocketHandler := &connect.Controller{}
	pushToClientHandler := &send2client.Controller{}
	pushToGroupHandler := &send2group.Controller{}
	bindToGroupHandler := &bindgroup.Controller{}

	http.HandleFunc("/ws", websocketHandler.Run)
	http.HandleFunc("/send_to_client", pushToClientHandler.Run)
	http.HandleFunc("/send_to_group", pushToGroupHandler.Run)
	http.HandleFunc("/bind_to_group", bindToGroupHandler.Run)

	go server.WriteMessage()
}
