package src

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type WebsocketHandler struct {
	upgrader     *websocket.Upgrader
	binder       *binder
	toClientChan chan [2]string
}

type toClient struct {
	ClientId string `json:"clientId"`
}

type inputData struct {
	ClientId string `json:"clientId"`
	Message  string `json:"message"`
}

func (wh *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	clintId := GenUUID()

	//给客户端绑定ID
	wh.binder.BindToMap(clintId, conn)

	//返回给客户端
	if err = conn.WriteJSON(toClient{ClientId: clintId}); err != nil {
		conn.Close()
	}

	wh.readMessage(conn)
}

func (wh *WebsocketHandler) readMessage(conn *websocket.Conn) {
	for {
		var inputData inputData
		err := conn.ReadJSON(&inputData)
		if err != nil {
			log.Printf("read error: %v", err)
		}

		fmt.Println(inputData.Message)

		wh.toClientChan <- []string{inputData.ClientId, inputData.Message}
	}
}

func (wh *WebsocketHandler) WriteMessage() {
	for {
		select {
		case clientInfo := <-wh.toClientChan:
			toConn, ok := wh.binder.clintId2ConnMap[clientInfo[0]];
			if ok {
				_ = toConn.Conn.WriteJSON(clientInfo[1]);
			}
		}
	}
}
