package client

import (
	"github.com/gorilla/websocket"
	"sync"
)

//连接列表
var Client2ConnMu sync.RWMutex
var Clint2ConnMap map[string]*websocket.Conn

//客户端所属分组列表
var ClientGroupsMu sync.RWMutex
var ClientGroupsMap map[string][]string

//分组里的客户端列表
var GroupClientIdsMu sync.RWMutex
var GroupClientIds map[string][]string

//初始化变量
func Init() {
	ClientGroupsMap = make(map[string][]string, 0);
	Clint2ConnMap = make(map[string]*websocket.Conn);
	GroupClientIds = make(map[string][]string, 0);
}

//给客户端绑定ID
func AddClient(clientId string, conn *websocket.Conn) {
	Client2ConnMu.Lock()
	defer Client2ConnMu.Unlock()
	Clint2ConnMap[clientId] = conn
}

//删除客户端
func DelClient(clientId string) {
	Client2ConnMu.Lock()
	defer Client2ConnMu.Unlock()
	delete(Clint2ConnMap, clientId)

}

//删除客户端里的分组
func DelClientGroup(clientId string) {
	ClientGroupsMu.Lock()
	defer ClientGroupsMu.Unlock()
	delete(ClientGroupsMap, clientId)
}

//获取分组里的客户端列表
func GetGroupClientIds(groupName string) ([]string) {
	GroupClientIdsMu.Lock()
	defer GroupClientIdsMu.Unlock()
	return GroupClientIds[groupName]
}

//获取客户端分组列表
func GetClientGroups(clientId string) []string {
	return ClientGroupsMap[clientId]
}

//客户端数量
func ClientNumber() int {
	return len(Clint2ConnMap)
}

//客户端是否存在
func IsAlive(clientId string) (conn *websocket.Conn, ok bool) {
	conn, ok = Clint2ConnMap[clientId];
	return
}

//添加分组
func AddClientToGroup(groupName, clientId string) {
	ClientGroupsMu.Lock()
	defer ClientGroupsMu.Unlock()
	ClientGroupsMap[clientId] = append(ClientGroupsMap[clientId], groupName)
}
