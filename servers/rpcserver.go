package servers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go-websocket/pkg/setting"
	"go-websocket/servers/pb"
	"go-websocket/tools/util"
	"google.golang.org/grpc"
	"net"
)

type CommonServiceServer struct{}

func (this *CommonServiceServer) Send2Client(ctx context.Context, req *pb.Send2ClientReq) (*pb.Send2ClientReply, error) {
	log.WithFields(log.Fields{
		"host":     setting.GlobalSetting.LocalHost,
		"port":     setting.CommonSetting.HttpPort,
		"clientId": req.ClientId,
	}).Info("接收到RPC指定客户端消息")
	SendMessage2LocalClient(req.MessageId, req.ClientId, req.SendUserId, int(req.Code), req.Message, &req.Data)
	return &pb.Send2ClientReply{}, nil
}

func (this *CommonServiceServer) CloseClient(ctx context.Context, req *pb.CloseClientReq) (*pb.CloseClientReply, error) {
	log.WithFields(log.Fields{
		"host":     setting.GlobalSetting.LocalHost,
		"port":     setting.CommonSetting.HttpPort,
		"clientId": req.ClientId,
	}).Info("接收到RPC关闭连接")
	CloseLocalClient(req.ClientId, req.SystemId)
	return &pb.CloseClientReply{}, nil
}

//添加分组到group
func (this *CommonServiceServer) BindGroup(ctx context.Context, req *pb.BindGroupReq) (*pb.BindGroupReply, error) {
	if client, err := Manager.GetByClientId(req.ClientId); err == nil {
		//添加到本地
		Manager.AddClient2LocalGroup(req.GroupName, client, req.UserId, req.Extend)
	} else {
		log.Error("添加分组失败" + err.Error())
	}
	return &pb.BindGroupReply{}, nil
}

func (this *CommonServiceServer) Send2Group(ctx context.Context, req *pb.Send2GroupReq) (*pb.Send2GroupReply, error) {
	log.WithFields(log.Fields{
		"host": setting.GlobalSetting.LocalHost,
		"port": setting.CommonSetting.HttpPort,
	}).Info("接收到RPC发送分组消息")
	Manager.SendMessage2LocalGroup(req.SystemId, req.MessageId, req.SendUserId, req.GroupName, int(req.Code), req.Message, &req.Data)
	return &pb.Send2GroupReply{}, nil
}

func (this *CommonServiceServer) Send2System(ctx context.Context, req *pb.Send2SystemReq) (*pb.Send2SystemReply, error) {
	log.WithFields(log.Fields{
		"host": setting.GlobalSetting.LocalHost,
		"port": setting.CommonSetting.HttpPort,
	}).Info("接收到RPC发送系统消息")
	Manager.SendMessage2LocalSystem(req.SystemId, req.MessageId, req.SendUserId, int(req.Code), req.Message, &req.Data)
	return &pb.Send2SystemReply{}, nil
}

//获取分组在线用户列表
func (this *CommonServiceServer) GetGroupClients(ctx context.Context, req *pb.GetGroupClientsReq) (*pb.GetGroupClientsReply, error) {
	response := pb.GetGroupClientsReply{}
	response.List = Manager.GetGroupClientList(util.GenGroupKey(req.SystemId, req.GroupName))
	return &response, nil
}

func InitGRpcServer() {
	go createGRPCServer(":" + setting.CommonSetting.RPCPort)
}

func createGRPCServer(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterCommonServiceServer(s, &CommonServiceServer{})

	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}
