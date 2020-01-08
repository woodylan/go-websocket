package src

import (
	"context"
	"github.com/smallnest/rpcx/server"
	"go-websocket/define"
)

type RPCServer struct {
}

type Push2ClientArgs struct {
	ClientId string
	Message  string
}

type Response struct {
	Success bool
}

func (s *RPCServer) Push2Client(ctx context.Context, args *Push2ClientArgs, response *Response) error {
	SendMessage2LocalClient(args.ClientId, args.Message)
	return nil
}

func InitRpcServer(port string) {
	define.RPCPort = port
	go createServer("tcp", ":"+port);
	return
}

func createServer(network string, address string) {
	s := server.NewServer()
	err := s.Register(new(RPCServer), "")
	if err != nil {
		panic(err)
	}
	err = s.Serve(network, address)
	if err != nil {
		panic(err)
	}

	return
}
