package servers

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/server"
	"go-websocket/define"
)

type RPCServer struct {
}

type Push2ClientArgs struct {
	ClientId string
	Code     int
	Message  string
	Data     interface{}
}

type Push2GroupArgs struct {
	SystemId  string
	GroupName string
	Code      int
	Message   string
	Data      interface{}
}

type Push2SystemArgs struct {
	SystemId string
	Code     int
	Message  string
	Data     interface{}
}

type AddClient2GroupArgs struct {
	SystemId  string
	GroupName string
	ClientId  string
}

type Response struct {
	Success bool
}

func (s *RPCServer) Push2Client(ctx context.Context, args *Push2ClientArgs, response *Response) error {
	fmt.Println("接收到RPC消息:发送指定客户端消息")
	SendMessage2LocalClient(&args.ClientId, args.Code, args.Message, &args.Data)
	return nil
}

func (s *RPCServer) Push2Group(ctx context.Context, args *Push2GroupArgs, response *Response) error {
	fmt.Println("接收到RPC消息:发送分组消息")
	Manager.SendMessage2LocalGroup(&args.SystemId, &args.GroupName, args.Code, args.Message, &args.Data)
	return nil
}

func (s *RPCServer) Push2System(ctx context.Context, args *Push2SystemArgs, response *Response) error {
	fmt.Println("接收到RPC消息:发送系统消息")
	Manager.SendMessage2LocalSystem(&args.SystemId, args.Code, args.Message, &args.Data)
	return nil
}

//添加分组到group
func (s *RPCServer) AddClient2Group(ctx context.Context, args *AddClient2GroupArgs, response *Response) error {
	AddClient2Group(&args.SystemId, &args.GroupName, args.ClientId)
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
