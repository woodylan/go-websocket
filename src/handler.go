package src

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// 最大的消息大小
	maxMessageSize = 8192

	// 心跳间隔
	heartbeatInterval = 10 * time.Second
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

	clientId := GenUUID()

	//给客户端绑定ID
	wh.binder.BindToMap(clientId, conn)

	//返回给客户端
	if err = conn.WriteJSON(toClient{ClientId: clientId}); err != nil {
		_ = conn.Close()
	}

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	//发送心跳
	wh.SendJump(conn)

	//读取消息并发送
	wh.readMessage(conn, clientId)
}

func (wh *WebsocketHandler) readMessage(conn *websocket.Conn, clientId string) {
	for {
		var inputData inputData

		err := conn.ReadJSON(&inputData)
		if err != nil {
			log.Printf("read error: %v", err)
			//删除这个客户端
			wh.binder.DelMap(clientId)
			return
		}

		wh.toClientChan <- [2]string{inputData.ClientId, inputData.Message}
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

//发送心跳数据
func (wh *WebsocketHandler) SendJump(conn *websocket.Conn) {
	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteJSON("heartbeat"); err != nil {
				return
			}
			time.Sleep(heartbeatInterval)
		}

	}()
}
