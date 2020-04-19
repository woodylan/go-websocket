package servers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"go-websocket/define"
	"go-websocket/tools/util"
)

type RPCServer struct {
}

type Push2ClientArgs struct {
	MessageId  string
	ClientId   string
	SendUserId string
	Code       int
	Message    string
	Data       interface{}
}

type CloseClientArgs struct {
	SystemId string
	ClientId string
}

type Push2GroupArgs struct {
	MessageId  string
	SystemId   string
	SendUserId string
	GroupName  string
	Code       int
	Message    string
	Data       interface{}
}

type Push2SystemArgs struct {
	MessageId  string
	SystemId   string
	SendUserId string
	Code       int
	Message    string
	Data       interface{}
}

type AddClient2GroupArgs struct {
	SystemId  string
	GroupName string
	ClientId  string
	UserId    string
	Extend    string
}

type GetGroupListArgs struct {
	SystemId  string
	GroupName string
}

type Response struct {
	Success bool
}

type GroupListResponse struct {
	List []string
}

func (s *RPCServer) Push2Client(ctx context.Context, args *Push2ClientArgs, response *Response) error {
	log.WithFields(log.Fields{
		"host":     define.LocalHost,
		"port":     define.Port,
		"clientId": args.ClientId,
	}).Info("接收到RPC指定客户端消息")
	SendMessage2LocalClient(args.MessageId, args.ClientId, args.SendUserId, args.Code, args.Message, &args.Data)
	return nil
}

func (s *RPCServer) CloseClient(ctx context.Context, args *CloseClientArgs, response *Response) error {
	log.WithFields(log.Fields{
		"host":     define.LocalHost,
		"port":     define.Port,
		"clientId": args.ClientId,
	}).Info("接收到RPC关闭连接")
	CloseLocalClient(args.ClientId, args.SystemId)
	return nil
}

func (s *RPCServer) Push2Group(ctx context.Context, args *Push2GroupArgs, response *Response) error {
	log.WithFields(log.Fields{
		"host": define.LocalHost,
		"port": define.Port,
	}).Info("接收到RPC发送分组消息")
	Manager.SendMessage2LocalGroup(args.SystemId, args.MessageId, args.SendUserId, args.GroupName, args.Code, args.Message, &args.Data)
	return nil
}

func (s *RPCServer) Push2System(ctx context.Context, args *Push2SystemArgs, response *Response) error {
	log.WithFields(log.Fields{
		"host": define.LocalHost,
		"port": define.Port,
	}).Info("接收到RPC发送系统消息")
	Manager.SendMessage2LocalSystem(args.SystemId, args.MessageId, args.SendUserId, args.Code, args.Message, &args.Data)
	return nil
}

//添加分组到group
func (s *RPCServer) AddClient2Group(ctx context.Context, args *AddClient2GroupArgs, response *Response) error {
	if client, err := Manager.GetByClientId(args.ClientId); err == nil {
		//添加到本地
		Manager.AddClient2LocalGroup(args.GroupName, client, args.UserId, args.Extend)
	} else {
		log.Error("添加分组失败" + err.Error())
	}
	return nil
}

//获取分组在线用户列表
func (s *RPCServer) GetOnlineList(ctx context.Context, args *GetGroupListArgs, response *GroupListResponse) error {
	response.List = Manager.GetGroupClientList(util.GenGroupKey(args.SystemId, args.GroupName))
	return nil
}

func InitRpcServer(port string) {
	define.RPCPort = port
	go createServer("tcp", ":"+port)
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
}
