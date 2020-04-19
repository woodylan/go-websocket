package servers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
	"go-websocket/define"
	"go-websocket/pkg/redis"
	"sync"
	"time"
)

//客户端列表
func getServerList() []*client.KVPair {
	serverList, err := redis.SMEMBERS(define.REDIS_KEY_SERVER_LIST)
	if err != nil {
		log.Errorf("failed to get server list: %v", err)
		return []*client.KVPair{
			{Key: define.LocalHost + ":" + define.RPCPort},
		}
	}

	var clientList []*client.KVPair
	for _, host := range serverList {
		clientList = append(clientList, &client.KVPair{Key: host})
	}

	return clientList
}

//获取单台客户端
func getXClient(addr string) (XClient client.XClient) {
	d := client.NewPeer2PeerDiscovery(addr, "")
	XClient = client.NewXClient("RPCServer", client.Failfast, client.RandomSelect, d, client.Option{
		Retries:        3,
		RPCPath:        share.DefaultRPCPath,
		ConnectTimeout: 100 * time.Millisecond,
		SerializeType:  protocol.MsgPack,
		CompressType:   protocol.None,
		BackupLatency:  10 * time.Millisecond,
	})
	return
}

//获取多台客户端
func getXClients() (XClient client.XClient) {
	d := client.NewMultipleServersDiscovery(getServerList())
	XClient = client.NewXClient("RPCServer", client.Failtry, client.RandomSelect, d, client.Option{
		Retries:        3,
		RPCPath:        share.DefaultRPCPath,
		ConnectTimeout: 100 * time.Millisecond,
		SerializeType:  protocol.MsgPack,
		CompressType:   protocol.None,
		BackupLatency:  10 * time.Millisecond,
	})
	return
}

func SendRpc2Client(addr string, messageId, sendUserId, clientId string, code int, message string, data *interface{}) {
	XClient := getXClient(addr)
	defer XClient.Close()

	log.WithFields(log.Fields{
		"host":     define.LocalHost,
		"port":     define.Port,
		"add":      addr,
		"clientId": clientId,
		"msg":      (*data).(string),
	}).Info("发送到服务器")
	err := XClient.Call(context.Background(), "Push2Client", &Push2ClientArgs{MessageId: messageId, SendUserId: sendUserId, ClientId: clientId, Code: code, Message: message, Data: data}, &Response{})

	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

func CloseRpcClient(addr string, clientId, systemId string) {
	XClient := getXClient(addr)
	defer XClient.Close()

	log.WithFields(log.Fields{
		"host":     define.LocalHost,
		"port":     define.Port,
		"add":      addr,
		"clientId": clientId,
	}).Info("发送关闭连接到服务器")
	err := XClient.Call(context.Background(), "CloseClient", &CloseClientArgs{SystemId: systemId, ClientId: clientId}, &Response{})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

//绑定分组
func SendRpcBindGroup(addr *string, systemId string, groupName string, clientId string, userId string, extend string) {
	XClient := getXClient(*addr)
	defer XClient.Close()

	err := XClient.Call(context.Background(), "AddClient2Group", &AddClient2GroupArgs{SystemId: systemId, GroupName: groupName, ClientId: clientId, UserId: userId, Extend: extend}, &Response{})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

//发送分组消息
func SendGroupBroadcast(systemId string, messageId, sendUserId, groupName string, code int, message string, data *interface{}) {
	XClient := getXClients()
	defer XClient.Close()

	err := XClient.Broadcast(context.Background(), "Push2Group", &Push2GroupArgs{MessageId: messageId, SystemId: systemId, SendUserId: sendUserId, GroupName: groupName, Code: code, Message: message, Data: data}, &Response{})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

//发送系统信息
func SendSystemBroadcast(systemId, messageId, sendUserId string, code int, message string, data *interface{}) {
	XClient := getXClients()
	defer XClient.Close()

	err := XClient.Broadcast(context.Background(), "Push2System", &Push2SystemArgs{MessageId: messageId, SystemId: systemId, SendUserId: sendUserId, Code: code, Message: message, Data: data}, &Response{})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

func GetOnlineListBroadcast(systemId *string, groupName *string) (clientIdList []string) {
	serverList := getServerList()
	serverCount := len(serverList)

	onlineListChan := make(chan []string, serverCount)
	var wg sync.WaitGroup

	wg.Add(serverCount)
	for _, server := range serverList {
		go func(add string) {
			XClient := getXClient(add)

			response := &GroupListResponse{}
			err := XClient.Call(context.Background(), "GetOnlineList", &GetGroupListArgs{SystemId: *systemId, GroupName: *groupName}, response)
			_ = XClient.Close()
			if err != nil {
				log.Errorf("failed to call: %v", err)
			} else {
				onlineListChan <- response.List
			}
			wg.Done()

		}(server.Key)
	}

	wg.Wait()

	for i := 1; i <= len(serverList); i++ {
		list, ok := <-onlineListChan
		if ok {
			clientIdList = append(clientIdList, list...)
		} else {
			return
		}
	}
	close(onlineListChan)

	return
}
