package client

import (
	"github.com/gorilla/websocket"
	"go-websocket/tools/util"
	"sync"
)

//连接列表
var Client2ConnMu sync.RWMutex
var ClintId2ConnMap map[string]*websocket.Conn

//客户端所属分组列表
var ClientGroupsMu sync.RWMutex
var ClientGroupsMap map[string][]string

//分组里的客户端列表
var GroupClientIdsMu sync.RWMutex
var GroupClientIds map[string][]string

//给客户端绑定ID
func AddClient(clientId string, conn *websocket.Conn) {
	Client2ConnMu.Lock()
	defer Client2ConnMu.Unlock()
	ClintId2ConnMap[clientId] = conn
}

//删除客户端
func DelClient(clientId string) {
	Client2ConnMu.Lock()
	delete(ClintId2ConnMap, clientId)
	Client2ConnMu.Unlock()
	if util.IsCluster() {
		//todo 删除redis
		//todo 删除集群里的分组信息
	} else {
		//删除单机里的分组
		ClientGroupsMu.Lock()
		defer ClientGroupsMu.Unlock()
		delete(ClientGroupsMap, clientId)
	}
}

func GetGroupClientIds(groupName string) ([]string) {
	GroupClientIdsMu.Lock()
	defer GroupClientIdsMu.Unlock()
	return GroupClientIds[groupName]
}

//客户端数量
func ClientNumber() int {
	return len(ClintId2ConnMap)
}

//客户端是否存在
func IsAlive(clientId string) (conn *websocket.Conn, ok bool) {
	conn, ok = ClintId2ConnMap[clientId];
	return
}
