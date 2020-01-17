package servers

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	ClientId      string          // 标识ID
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	AppId         uint32          // 登录的平台Id app/web/ios
	UserId        string          // 用户Id，用户登录以后才有
	FirstTime     uint64          // 首次连接事件
	HeartbeatTime uint64          // 用户上次心跳时间
	LoginTime     uint64          // 登录时间 登录以后才有
}

func NewClient(clientId string, addr string, socket *websocket.Conn, firstTime uint64) (*Client) {
	return &Client{
		ClientId:      clientId,
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
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