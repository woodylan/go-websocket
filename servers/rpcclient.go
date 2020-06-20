package servers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go-websocket/define"
	"go-websocket/servers/pb"
	"google.golang.org/grpc"
	"sync"
)

func grpcConn(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}
	return conn
}

func SendRpc2Client(addr string, messageId, sendUserId, clientId string, code int, message string, data *string) {
	conn := grpcConn(addr)
	defer conn.Close()

	log.WithFields(log.Fields{
		"host":     define.LocalHost,
		"port":     define.Port,
		"add":      addr,
		"clientId": clientId,
		"msg":      data,
	}).Info("发送到服务器")

	c := pb.NewCommonServiceClient(conn)
	_, err := c.Send2Client(context.Background(), &pb.Send2ClientReq{
		MessageId:  messageId,
		SendUserId: sendUserId,
		ClientId:   clientId,
		Code:       int32(code),
		Message:    message,
		Data:       *data,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

func CloseRpcClient(addr string, clientId, systemId string) {
	conn := grpcConn(addr)
	defer conn.Close()

	log.WithFields(log.Fields{
		"host":     define.LocalHost,
		"port":     define.Port,
		"add":      addr,
		"clientId": clientId,
	}).Info("发送关闭连接到服务器")

	c := pb.NewCommonServiceClient(conn)
	_, err := c.CloseClient(context.Background(), &pb.CloseClientReq{
		SystemId: systemId,
		ClientId: clientId,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

//绑定分组
func SendRpcBindGroup(addr string, systemId string, groupName string, clientId string, userId string, extend string) {
	conn := grpcConn(addr)
	defer conn.Close()

	c := pb.NewCommonServiceClient(conn)
	_, err := c.BindGroup(context.Background(), &pb.BindGroupReq{
		SystemId:  systemId,
		GroupName: groupName,
		ClientId:  clientId,
		UserId:    userId,
		Extend:    extend,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

//发送分组消息
func SendGroupBroadcast(systemId string, messageId, sendUserId, groupName string, code int, message string, data *string) {
	define.ServerListLock.Lock()
	defer define.ServerListLock.Unlock()
	for _, addr := range define.ServerList {
		conn := grpcConn(addr)
		defer conn.Close()

		c := pb.NewCommonServiceClient(conn)
		_, err := c.Send2Group(context.Background(), &pb.Send2GroupReq{
			SystemId:   systemId,
			MessageId:  messageId,
			SendUserId: sendUserId,
			GroupName:  groupName,
			Code:       int32(code),
			Message:    message,
			Data:       *data,
		})
		if err != nil {
			log.Errorf("failed to call: %v", err)
		}
	}
}

//发送系统信息
func SendSystemBroadcast(systemId string, messageId, sendUserId string, code int, message string, data *string) {
	define.ServerListLock.Lock()
	defer define.ServerListLock.Unlock()
	for _, addr := range define.ServerList {
		conn := grpcConn(addr)
		defer conn.Close()

		c := pb.NewCommonServiceClient(conn)
		_, err := c.Send2System(context.Background(), &pb.Send2SystemReq{
			SystemId:   systemId,
			MessageId:  messageId,
			SendUserId: sendUserId,
			Code:       int32(code),
			Message:    message,
			Data:       *data,
		})
		if err != nil {
			log.Errorf("failed to call: %v", err)
		}
	}
}

func GetOnlineListBroadcast(systemId *string, groupName *string) (clientIdList []string) {
	define.ServerListLock.Lock()
	defer define.ServerListLock.Unlock()

	serverCount := len(define.ServerList)

	onlineListChan := make(chan []string, serverCount)
	var wg sync.WaitGroup

	wg.Add(serverCount)
	for _, addr := range define.ServerList {
		go func(addr string) {
			conn := grpcConn(addr)
			defer conn.Close()
			c := pb.NewCommonServiceClient(conn)
			response, err := c.GetGroupClients(context.Background(), &pb.GetGroupClientsReq{
				SystemId:  *systemId,
				GroupName: *groupName,
			})
			if err != nil {
				log.Errorf("failed to call: %v", err)
			} else {
				onlineListChan <- response.List
			}
			wg.Done()

		}(addr)
	}

	wg.Wait()

	for i := 1; i <= serverCount; i++ {
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
