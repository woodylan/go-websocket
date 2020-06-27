package main

import (
	"fmt"
	"go-websocket/configs"
	"go-websocket/define"
	"go-websocket/pkg/etcd"
	"go-websocket/routers"
	"go-websocket/servers"
	_ "go-websocket/tools/log"
	"go-websocket/tools/util"
	"net"
	"net/http"
)

func main() {
	//初始化配置文件
	if err := configs.InitConfig(); err != nil {
		panic(err)
	}

	//初始化RPC服务
	initRPCServer()

	//将服务器地址、端口注册到etcd中
	registerServer()

	//初始化路由
	routers.Init()

	//启动一个定时器用来发送心跳
	servers.PingTimer()

	fmt.Printf("服务器启动成功，端口号：%s\n", configs.Conf.CommonConf.Port)

	if err := http.ListenAndServe(":"+configs.Conf.CommonConf.Port, nil); err != nil {
		panic(err)
	}
}

func initRPCServer() {
	//如果是集群，则启用RPC进行通讯
	if util.IsCluster() {
		//初始化RPC服务
		servers.InitGRpcServer()
		fmt.Printf("启动RPC，端口号：%s\n", configs.Conf.CommonConf.RPCPort)
	}
}

//ETCD注册发现服务
func registerServer() {
	if util.IsCluster() {
		//注册租约
		ser, err := etcd.NewServiceReg(configs.Conf.EtcdEndpoints, 5)
		if err != nil {
			panic(err)
		}

		hostPort := net.JoinHostPort(configs.Conf.CommonConf.LocalHost, configs.Conf.CommonConf.RPCPort)
		//添加key
		err = ser.PutService(define.ETCD_SERVER_LIST+hostPort, hostPort)
		if err != nil {
			panic(err)
		}

		cli, err := etcd.NewClientDis(configs.Conf.EtcdEndpoints)
		if err != nil {
			panic(err)
		}
		_, err = cli.GetService(define.ETCD_SERVER_LIST)
		if err != nil {
			panic(err)
		}
	}
}
