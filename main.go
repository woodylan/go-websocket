package main

import (
	"fmt"
	"go-websocket/define"
	"go-websocket/routers"
	"go-websocket/servers/client"
	"go-websocket/servers/server"
	"go-websocket/tools/readconfig"
	"go-websocket/tools/util"
	"net/http"
	"os"
)

func main() {
	port := getPort()

	//初始化配置文件
	readconfig.InitConfig()

	//初始化RPC服务
	initRPCServer(port)

	//记录本机内网IP地址
	define.LocalHost = util.GetIntranetIp()

	//初始化rabbitMQ
	server.Init()

	//初始化路由
	routers.Init()

	//初始化变量
	client.Init()

	fmt.Printf("服务器启动成功，端口号：%s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func initRPCServer(port string) {
	//如果是集群，则启用RPC进行通讯
	if util.IsCluster() {
		//初始化RPC服务
		rpcPort := util.GenRpcPort(port)
		fmt.Printf("启动RPC，端口号：%s\n", rpcPort)
		server.InitRpcServer(rpcPort)
	}
}

func getPort() string {
	port := "666"

	args := os.Args //获取用户输入的所有参数
	if args != nil && len(args) >= 2 && len(args[1]) != 0 {
		port = args[1];
	}

	return port
}
