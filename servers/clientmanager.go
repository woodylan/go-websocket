package servers

import (
	"errors"
	"log"
	"sync"
	"time"
)

// 连接管理
type ClientManager struct {
	ClientIdMap     map[string]*Client // 全部的连接
	ClientIdMapLock sync.RWMutex       // 读写锁

	Connect    chan *Client // 连接处理
	DisConnect chan *Client // 断开连接处理

	GroupLock sync.RWMutex
	Groups    map[string][]*Client
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		ClientIdMap: make(map[string]*Client),
		Connect:     make(chan *Client, 1000),
		DisConnect:  make(chan *Client, 1000),
		Groups:      make(map[string][]*Client, 100),
	}

	return
}

// 管道处理程序
func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.Connect:
			// 建立连接事件
			manager.EventConnect(client)
		case conn := <-manager.DisConnect:
			// 断开连接事件
			manager.EventDisconnect(conn)
		}
	}
}

// 建立连接事件
func (manager *ClientManager) EventConnect(client *Client) {
	manager.AddClient(client)

	log.Printf("客户端已连接: %s 总连接数：%d", client.ClientId, Manager.Count())
}

// 断开连接时间
func (manager *ClientManager) EventDisconnect(client *Client) {
	//关闭连接
	_ = client.Socket.Close()
	manager.DelClient(client)
	log.Printf("客户端已断开: %s 总连接数：%d 连接时间:%d秒 ", client.ClientId, Manager.Count(), uint64(time.Now().Unix())-client.ConnectTime)
}

// 添加客户端
func (manager *ClientManager) AddClient(client *Client) {
	manager.ClientIdMapLock.Lock()
	defer manager.ClientIdMapLock.Unlock()

	manager.ClientIdMap[client.ClientId] = client
}

// 获取所有的客户端
func (manager *ClientManager) AllClient() map[string]*Client {
	manager.ClientIdMapLock.RLock()
	defer manager.ClientIdMapLock.RUnlock()

	return manager.ClientIdMap
}

// 客户端数量
func (manager *ClientManager) Count() int {
	manager.ClientIdMapLock.RLock()
	defer manager.ClientIdMapLock.RUnlock()
	return len(manager.ClientIdMap)
}

// 删除客户端
func (manager *ClientManager) DelClient(client *Client) {
	manager.ClientIdMapLock.Lock()
	defer manager.ClientIdMapLock.Unlock()

	if _, ok := manager.ClientIdMap[client.ClientId]; ok {
		delete(manager.ClientIdMap, client.ClientId)
	}
}

// 通过clientId获取
func (manager *ClientManager) GetByClientId(clientId string) (*Client, error) {
	manager.ClientIdMapLock.RLock()
	defer manager.ClientIdMapLock.RUnlock()

	if client, ok := manager.ClientIdMap[clientId]; !ok {
		return nil, errors.New("客户端不存在")
	} else {
		return client, nil
	}
}

// 发送到本机分组
func (manager *ClientManager) SendMessage2LocalGroup(groupName *string, code int, msg string, data *interface{}) {
	if len(*groupName) > 0 {
		clientList := manager.GetGroupClientList(*groupName)
		if len(clientList) > 0 {
			for _, client := range clientList {
				SendMessage2LocalClient(&client.ClientId, code, msg, data)
			}
		}
	}
}

// 添加到本地分组
func (manager *ClientManager) AddClient2LocalGroup(groupName *string, client *Client) {
	manager.GroupLock.RLock()
	defer manager.GroupLock.RUnlock()
	manager.Groups[*groupName] = append(manager.Groups[*groupName], client)
}

// 获取本地分组的成员
func (manager *ClientManager) GetGroupClientList(groupName string) []*Client {
	manager.GroupLock.RLock()
	defer manager.GroupLock.RUnlock()
	return manager.Groups[groupName]
}
