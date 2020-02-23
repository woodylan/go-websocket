package servers

import (
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	ClientId    string          // 标识ID
	SystemId    string          // 系统ID
	Socket      *websocket.Conn // 用户连接
	ConnectTime uint64          // 首次连接时间
}

type SendData struct {
	Code int
	Msg  string
	Data *interface{}
}

func NewClient(clientId string, systemId string, socket *websocket.Conn) (*Client) {
	return &Client{
		ClientId:    clientId,
		SystemId:    systemId,
		Socket:      socket,
		ConnectTime: uint64(time.Now().Unix()),
	}
}

func (c *Client) Read() {
	go func() {
	loop:
		for {
			messageType, _, err := c.Socket.ReadMessage()
			if err != nil {
				if messageType == -1 || messageType == websocket.CloseMessage {
					//下线
					Manager.DisConnect <- c
					break loop
				}
			}
		}
	}()
}
