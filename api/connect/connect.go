package connect

import (
	"github.com/gorilla/websocket"
	"go-websocket/servers/client"
	"go-websocket/servers/server"
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

	closeMessage = make(chan bool)

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

	//发送心跳
	sendJump(clientId, conn)
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
}

//发送心跳数据
func sendJump(clientId string, conn *websocket.Conn) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second)); err != nil {
				_ = conn.Close()
				server.DelClient(clientId)
				log.Printf("发送心跳失败: %s 总连接数：%d", clientId, client.ClientNumber())
				return
			}
		case <-closeSign:
			server.DelClient(clientId)
			_ = conn.Close()
			return
		}
	}
}
