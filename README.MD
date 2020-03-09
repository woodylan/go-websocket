# Golang实现的分布式WebSocket微服务

[![Go](https://img.shields.io/badge/Go-1.13-blue.svg)](https://golang.google.cn)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/woodylan/go-websocket)
![Travis (.org)](https://api.travis-ci.com/woodylan/go-websocket.svg?branch=master)
[![star](https://img.shields.io/github/stars/woodylan/go-websocket?style=social)](https://github.com/woodylan/go-websocket/stargazers)

## 简介

本系统基于Golang、Redis、RabbitMQ、RPC实现分布式WebSocket微服务，也可以单机部署，单机部署不需要Redis、RabbitMQ和RPC。分布式部署可以支持nginx负责均衡、水平扩容部署，程序之间使用RabbitMQ广播、RPC通信。

基本流程为：用ws协议连接本服务，得到一个clientId，由客户端上报这个clinetId给服务端，服务端拿到这个clientId之后，可以给这个客户端发送信息，绑定这个客户端都分组，给分组发送消息。

目前实现的功能有，给指定客户端发送消息、绑定客户端到分组、给分组里的客户端批量发送消息。适用于长连接的大部分场景，分组可以理解为聊天室，绑定客户端到分组相当于把客户端添加到聊天室，给分组发送信息相当于给聊天室的每个人发送消息。



## 文档

1. [技术方案架构](docs/introduction.md)
2. [接口文档](docs/api.md)




## 使用

**下载本项目：**

```shell
git clone https://github.com/woodylan/go-websocket.git
```

**编译：**

```shell
// 编译适用于本机的版本
go build

// 编译Linux版本
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

// 编译Windows 64位版本
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

// 编译MacOS版本
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build
```

**执行：**

编译成功之后会得到一个二进制文件`go-websocket`，执行该二进制文件，文件名后面跟着的是端口号，下面的命令`666`则表示端口号，你可以可以改成其他的。

```shell
./go-websocket 666
```

**连接测试：**

打开支持Websocket的客户端，输入 `ws://127.0.0.1:666/ws` 进行连接，连接成功会返回`clientId`。



## 单机部署

单机部署很简单，不需要配置Redis、RabbitMQ，只需要编译然后运行该二进制文件就可以了，步骤如上。



## 分布式部署

**安装Redis：** 参考网上教程



**配置文件：**

配置文件位于项目根目录的`configs/config.ini`，`cluster`为true表示分布式部署，分布式部署[rabbitMQ]、[redis]项必填。

```ini
[common]
# 是否分布式部署
cluster = true
# 对称加密key 16位
crypto_key = xxxxxxxxxxxxxxxx

[redis]
host = 127.0.0.1
port = 6379
password =
```

**运行项目：**

在不同的机器运行本项目，注意配置号端口号，项目如果在同一机器，则必须用不同的端口。你可以用`supervisor`做进程管理。

**配置Nginx负载均衡：**

```nginx
upstream ws_cluster {
    server 127.0.0.1:666;
    server 127.0.0.1:667;
}

server {
    listen  660;
    server_name ws.example.com;

    access_log /logs/access.log;
    error_log /logs/error.log;

    location / {
        proxy_pass http://ws_cluster; # 代理转发地址
        proxy_http_version 1.1;

        proxy_read_timeout 60s; # 超时设置

        # 启用支持websocket连接
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

至此，项目部署完成。




## 实现的功能

- [x] 分布式
- [x] 账户授权模式
- [x] 不同业务系统隔离
- [x] 发送给指定客户端
- [x] 发送给指定分组
- [x] 上下线通知
- [x] 群广播
- [x] 错误日志
- [x] 参数校验
- [ ] 查询某个客户端是否在线



## 沟通交流
QQ群：1028314856

<a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=adbd35a10f2c25bd765a7830d17579f6e645634d22151b6d480d1b05026125e9">
<img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="Golang websocket分布式交" title="gorose-orm"></a>