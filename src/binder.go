package src

import (
	"github.com/gorilla/websocket"
	"sync"
)

type binder struct {
	mu sync.RWMutex

	clintId2ConnMap map[string]*Conn
}

//给客户端绑定ID
func (b *binder) BindToMap(clintId string, conn *websocket.Conn) {
	b.clintId2ConnMap[clintId] = &Conn{Conn: conn}
}
