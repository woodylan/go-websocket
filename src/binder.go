package src

import (
	"fmt"
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
func (b *binder) BindToMap(clientId string, conn *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clintId2ConnMap[clientId] = &Conn{Conn: conn}
}

//删除客户端
func (b *binder) DelMap(clientId string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clintId2ConnMap, clientId)
	delete(b.clientGroupsMap, clientId)
}

//客户端数量
func (b *binder) ClientNumber() int {
	return len(b.clintId2ConnMap)
}

//是否为单机
func (b *binder) isStandalone() bool {
	return true
}

//发送信息到指定客户端
func (b *binder) SendMessage2Client(clientId, message string) {
	if b.isStandalone() {
		//如果是单机服务，则只发送到本机
		toClientChan <- [2]string{clientId, message}
		fmt.Println(clientId, message)
	} else {

	}
}

//发送信息到指定分组
func (b *binder) SendMessage2Group(groupName, message string) {
	if b.isStandalone() {
		//如果是单机服务，则只发送到本机
		if len(groupName) > 0 {
			if clientList, ok := b.groupClientIds[groupName]; ok {
				for _, client := range clientList {
					//发送信息
					toClientChan <- [2]string{client, message}
				}
			}
		}
	} else {
		
	}

}
