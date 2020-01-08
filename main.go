package main

import (
	"fmt"
	"go-websocket/define"
	"go-websocket/src"
	"go-websocket/tools/readconfig"
	"go-websocket/tools/util"
	"os"
)

func main() {
	port := getPort()
	server := src.NewServer(":" + port)

	//读取配置文件
	readconfig.ReadConfig()

	//如果是集群，则启用RPC进行通讯
	if util.IsCluster() {
		//初始化RPC服务
		rpcPort := util.GenRpcPort(port)
		fmt.Printf("启动RPC，端口号：%s\n", rpcPort)
		src.InitRpcServer(rpcPort)
	}

	//记录本机内网IP地址
	define.LocalHost = util.GetIntranetIp()

	fmt.Printf("服务器启动成功，端口号：%s\n", port)
	if err := server.ListenAndServer(); err != nil {
		panic(err)
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
