package servers

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
)

func getXClient(addr string) (XClient client.XClient) {
	d := client.NewPeer2PeerDiscovery(addr, "")
	XClient = client.NewXClient("RPCServer", client.Failfast, client.RandomSelect, d, client.DefaultOption)
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

func SendRpcBindGroup(addr *string, groupName *string, clientId *string) {
	XClient := getXClient(*addr)
	defer XClient.Close()

	err := XClient.Call(context.Background(), "AddClient2Group", &AddClient2GroupArgs{GroupName: *groupName, ClientId: *clientId}, &Response{})
	if err != nil {
		_ = fmt.Errorf("failed to call: %v", err)
	}
}
