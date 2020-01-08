package src

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
)

func SendRpc2Client(addr string, clientId string, message string) {
	d := client.NewPeer2PeerDiscovery(addr, "")
	xclient := client.NewXClient("RPCServer", client.Failfast, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	err := xclient.Call(context.Background(), "Push2Client", &Push2ClientArgs{ClientId: clientId, Message: message}, &Response{})
	if err != nil {
		_ = fmt.Errorf("failed to call: %v", err)
	}
}
