package connect

import (
	"github.com/gorilla/websocket"
	"go-websocket/servers/client"
	"go-websocket/servers/server"
	"go-websocket/tools/util"
	"log"
	"net/http"
)

const (
	// 最大的消息大小
	maxMessageSize = 8192
)

// 关闭连接的信号
var closeSign chan bool

type Controller struct {
	Upgrader *websocket.Upgrader
}

type toClient struct {
	ClientId string `json:"clientId"`
}

var defaultUpgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Controller) Run(w http.ResponseWriter, r *http.Request) {
	c.Upgrader = defaultUpgrader
	conn, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	closeSign = make(chan bool, 100)

	clientId := util.GenClientId()

	//给客户端绑定ID
	client.AddClient(clientId, conn)

	//返回给客户端
	if err = conn.WriteJSON(toClient{ClientId: clientId}); err != nil {
		_ = conn.Close()
	}

	log.Printf("客户端已连接: %s 总连接数：%d", clientId, client.ClientNumber())

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	//读取客户端消息
	readMessage(conn, clientId)
}

//websocket客户端发送消息
func readMessage(conn *websocket.Conn, clientId string) {
	go func() {
		for {
			messageType, _, err := conn.ReadMessage()
			if err != nil {
				if messageType == -1 || messageType == websocket.CloseMessage {
					//关闭连接
					_ = conn.Close()
					server.DelClient(clientId)
					closeSign <- true
					log.Printf("客户端已下线: %s 总连接数：%d", clientId, client.ClientNumber())
					return
				}
			}
		}
	}()

	//接收关闭信号并关闭
	for {
		select {
		case <-closeSign:
			server.DelClient(clientId)
			_ = conn.Close()
			return
		}
	}
}
