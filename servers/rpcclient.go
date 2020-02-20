package servers

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
)

//客户端列表
func getServerList() []*client.KVPair {
	return []*client.KVPair{
		{Key: "127.0.0.1:8777"},
		{Key: "127.0.0.1:8778"}}
}

//获取单台客户端
func getXClient(addr string) (XClient client.XClient) {
	d := client.NewPeer2PeerDiscovery(addr, "")
	XClient = client.NewXClient("RPCServer", client.Failfast, client.RandomSelect, d, client.DefaultOption)
	return
}

//获取多台客户端
func getXClients() (XClient client.XClient) {
	d := client.NewMultipleServersDiscovery(getServerList())
	XClient = client.NewXClient("RPCServer", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	return
}

func SendRpc2Client(addr string, clientId *string, message string, data *interface{}) {
	XClient := getXClient(addr)
	defer XClient.Close()

	go fmt.Println("发送到服务器：" + addr + " 客户端：" + *clientId + " 消息：" + (*data).(string))
	err := XClient.Call(context.Background(), "Push2Client", &Push2ClientArgs{ClientId: *clientId, Message: message, Data: data}, &Response{})
	if err != nil {
		_ = fmt.Errorf("failed to call: %v", err)
	}
}

//绑定分组
func SendRpcBindGroup(addr *string, systemId *string, groupName *string, clientId *string) {
	XClient := getXClient(*addr)
	defer XClient.Close()

	err := XClient.Call(context.Background(), "AddClient2Group", &AddClient2GroupArgs{SystemId: *systemId, GroupName: *groupName, ClientId: *clientId}, &Response{})
	if err != nil {
		_ = fmt.Errorf("failed to call: %v", err)
	}
}

//发送分组消息
func SendGroupBroadcast(systemId *string, groupName *string, code int, message string, data *interface{}) {
	XClient := getXClients()
	defer XClient.Close()

	err := XClient.Broadcast(context.Background(), "Push2Group", &Push2GroupArgs{SystemId: *systemId, GroupName: *groupName, Code: code, Message: message, Data: data}, &Response{})
	if err != nil {
		_ = fmt.Errorf("failed to call: %v", err)
	}
}

//发送系统信息
func SendSystemBroadcast(systemId *string, code int, message string, data *interface{}) {
	XClient := getXClients()
	defer XClient.Close()

	err := XClient.Broadcast(context.Background(), "Push2System", &Push2SystemArgs{SystemId: *systemId, Code: code, Message: message, Data: data}, &Response{})
	if err != nil {
		_ = fmt.Errorf("failed to call: %v", err)
	}
}
