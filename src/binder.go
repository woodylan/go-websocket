package src

import (
	"github.com/gorilla/websocket"
	"sync"
)

type binder struct {
	mu sync.RWMutex

	clintId2ConnMap map[string]*Conn
	clientGroupsMap map[string][]string
	groupClientIds  map[string][]string
}

//给客户端绑定ID
func (b *binder) BindToMap(clintId string, conn *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clintId2ConnMap[clintId] = &Conn{Conn: conn}
}

//删除客户端
func (b *binder) DelMap(clintId string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clintId2ConnMap, clintId)
}

//客户端数量
func (b *binder) ClientNumber() int {
	return len(b.clintId2ConnMap)
}
