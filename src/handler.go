package src

import (
	"encoding/json"
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

//channel通道
var toClientChan chan [2]string

type WebsocketHandler struct {
	upgrader *websocket.Upgrader
	binder   *binder
}

type toClient struct {
	ClientId string `json:"clientId"`
}

type RetData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Date interface{} `json:"data"`
}

func init() {
	toClientChan = make(chan [2]string, 10)
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

	log.Printf("客户端已连接:%s 总连接数：%d", clientId, wh.binder.ClientNumber())

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	//发送心跳
	wh.SendJump(conn)

	//读取消息并发送 在这不提供
	//wh.readMessage(conn, clientId)

	//阻塞main线程
	select {}
}

//websocket客户端发送消息
//func (wh *WebsocketHandler) readMessage(conn *websocket.Conn, clientId string) {
//	for {
//		var inputData inputData
//
//		err := conn.ReadJSON(&inputData)
//		if err != nil {
//			log.Printf("read error: %v", err)
//			//删除这个客户端
//			wh.binder.DelMap(clientId)
//			return
//		}
//
//		toClientChan <- [2]string{inputData.ClientId, inputData.Message}
//	}
//}

func (wh *WebsocketHandler) WriteMessage() {
	for {
		select {
		case clientInfo := <-toClientChan:
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
		for {
			time.Sleep(heartbeatInterval)
			if err := conn.WriteJSON("heartbeat"); err != nil {
				return
			}
		}

	}()
}

func render(code int, msg string, data interface{}) (str string) {
	var retData RetData

	retData.Code = code
	retData.Msg = msg
	retData.Date = data

	retJson, _ := json.Marshal(retData)
	str = string(retJson)
	return
}
