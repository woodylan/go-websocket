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

func SendRpc2Client(addr string, clientId string, message string) {
	XClient := getXClient(addr)
	defer XClient.Close()
	//d := client.NewPeer2PeerDiscovery(addr, "")
	//xclient := client.NewXClient("RPCServer", client.Failfast, client.RandomSelect, d, client.DefaultOption)
	//defer xclient.Close()

	err := XClient.Call(context.Background(), "Push2Client", &Push2ClientArgs{ClientId: clientId, Message: message}, &Response{})
	if err != nil {
		_ = fmt.Errorf("failed to call: %v", err)
	}
}

func SendRpcBindGroup(addr string, groupName string, clientId string) {
	XClient := getXClient(addr)
	defer XClient.Close()
	//d := client.NewPeer2PeerDiscovery(addr, "")
	//xclient := client.NewXClient("RPCServer", client.Failfast, client.RandomSelect, d, client.DefaultOption)
	//defer xclient.Close()

	err := XClient.Call(context.Background(), "AddClient2Group", &AddClient2GroupArgs{GroupName: groupName, ClientId: clientId}, &Response{})
	if err != nil {
		_ = fmt.Errorf("failed to call: %v", err)
	}
}
