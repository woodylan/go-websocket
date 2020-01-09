package clientvar

import (
	"github.com/gorilla/websocket"
	"go-websocket/tools/util"
	"sync"
)

//客户端的分组列表
var ClientGroupsMapMu sync.RWMutex
var ClientGroupsMap map[string][]string

var MuClient2ConnMap sync.RWMutex
var ClintId2ConnMap map[string]*websocket.Conn

var MuGroupClientIds sync.RWMutex
var GroupClientIds map[string][]string

//给客户端绑定ID
func AddClient(clientId string, conn *websocket.Conn) {
	MuClient2ConnMap.Lock()
	defer MuClient2ConnMap.Unlock()
	ClintId2ConnMap[clientId] = conn
}

//删除客户端
func DelClient(clientId string) {
	MuClient2ConnMap.Lock()
	delete(ClintId2ConnMap, clientId)
	MuClient2ConnMap.Unlock()
	if util.IsCluster() {
		//todo 删除redis
		//todo 删除集群里的分组信息
	} else {
		//删除单机里的分组
		ClientGroupsMapMu.Lock()
		defer ClientGroupsMapMu.Unlock()
		delete(ClientGroupsMap, clientId)
	}
}

func GetGroupClientIds(groupName string) ([]string) {
	MuGroupClientIds.Lock()
	defer MuGroupClientIds.Unlock()
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
