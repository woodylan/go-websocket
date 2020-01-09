package connect

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go-websocket/clientvar"
	"go-websocket/servers"
	"go-websocket/tools/util"
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
//var toClientChan chan [2]string

type WebsocketHandler struct {
	Upgrader *websocket.Upgrader
}

type toClient struct {
	ClientId string `json:"clientId"`
}

func init() {
	servers.ToClientChan = make(chan [2]string, 10)
}

func (wh *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wh.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	clientId := util.GenClientId()

	//给客户端绑定ID
	clientvar.AddClient(clientId, conn)

	//返回给客户端
	if err = conn.WriteJSON(toClient{ClientId: clientId}); err != nil {
		_ = conn.Close()
	}

	log.Printf("客户端已连接:%s 总连接数：%d", clientId, clientvar.ClientNumber())

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	//发送心跳
	wh.SendJump(conn)

	//读取消息并发送 在这不提供
	//wh.readMessage(conn, clientId)

	//阻塞main线程
	select {}
}

func (wh *WebsocketHandler) WriteMessage() {
	for {
		select {
		case clientInfo := <-servers.ToClientChan:
			toConn, ok := clientvar.IsAlive(clientInfo[0]);
			if ok {
				err := toConn.WriteJSON(clientInfo[1]);
				if err != nil {
					go clientvar.DelClient(clientInfo[0])
					log.Println(err)
				} else {
					//todo 给redis续命
				}
			} else {
				go clientvar.DelClient(clientInfo[0])
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
				//todo 删除客户端
				fmt.Printf("删除客户端")
				return
			}
		}

	}()
}
