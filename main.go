package main

import (
	"fmt"
	"go-websocket/src"
	"go-websocket/src/readConfig"
	"os"
)

func main() {
	port := getPort()
	server := src.NewServer(":" + port)

	//读取配置文件
	readConfig.ReadConfig()

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
