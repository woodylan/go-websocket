package connect

import (
	"github.com/gorilla/websocket"
	"go-websocket/api"
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

type Controller struct {
	Upgrader *websocket.Upgrader
}

type renderData struct {
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

	clientId := util.GenClientId()

	//给客户端绑定ID
	client.AddClient(&clientId, conn)

	//返回给客户端
	if err = api.ConnRender(conn, renderData{ClientId: clientId}); err != nil {
		_ = conn.Close()
		return
	}

	log.Printf("客户端已连接: %s 总连接数：%d", clientId, client.ClientNumber())

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	//读取客户端消息
	readMessage(conn, &clientId)

	select {}
}

//websocket客户端发送消息
func readMessage(conn *websocket.Conn, clientId *string) {
	go func() {
	loop:
		for {
			messageType, _, err := conn.ReadMessage()
			if err != nil {
				if messageType == -1 || messageType == websocket.CloseMessage {
					//关闭连接
					_ = conn.Close()
					server.DelClient(clientId)
					log.Printf("客户端已下线: %s 总连接数：%d", *clientId, client.ClientNumber())
					break loop
				}
			}
		}
	}()
}
