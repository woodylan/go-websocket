# Golang实现的分布式WebSocket微服务

[![Go](https://img.shields.io/badge/Go-1.13-blue.svg)](https://golang.google.cn)
![GitHub release](https://img.shields.io/github/v/release/woodylan/go-websocket)
![Travis (.org)](https://api.travis-ci.com/woodylan/go-websocket.svg?branch=master)
[![star](https://img.shields.io/github/stars/woodylan/go-websocket?style=social)](https://github.com/woodylan/go-websocket/stargazers)

## 简介

本系统基于Golang、ETCD、RPC实现分布式WebSocket微服务，也可以单机部署，单机部署不需要ETCD、RPC。分布式部署可以支持nginx负责均衡、水平扩容部署，程序之间使用RPC通信。

基本流程为：用ws协议连接本服务，得到一个clientId，由客户端上报这个clinetId给服务端，服务端拿到这个clientId之后，可以给这个客户端发送信息，绑定这个客户端都分组，给分组发送消息。

目前实现的功能有，给指定客户端发送消息、绑定客户端到分组、给分组里的客户端批量发送消息、获取在线的客户端、上下线自动通知。适用于长连接的大部分场景，分组可以理解为聊天室，绑定客户端到分组相当于把客户端添加到聊天室，给分组发送信息相当于给聊天室的每个人发送消息。



## 文档

1. [技术方案架构](docs/introduction.md)
2. [接口文档](docs/api.md)



## SDK

1. PHP版：https://github.com/woodylan/go-websocket-php-sdk




## 使用

**下载本项目：**

这里已经打包好了，下载相应的环境，支持Linux、Windows、MacOS环境。

<https://github.com/woodylan/go-websocket/releases>

**你也可以选择自己编译：**

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

编译成功之后会得到一个二进制文件`go-websocket`，执行该二进制文件。

```shell
./go-websocket -c ./conf/app.ini
```

**连接测试：**

打开支持Websocket的客户端，输入 `ws://127.0.0.1:6000/ws` 进行连接，连接成功会返回`clientId`。



## docker体验

### 体验单机

1. 构建镜像

   ```shell
   docker build -t go-websocket .
   ```

2. 基于镜像运行容器

   ```shell
   docker run -tid -p 6000:6000 go-websocket
   ```

### 体验集群，同时运行ETCD集群
在当前目录下，直接运行` docker-compose up` 即可体验。



## 配置

**配置文件：**

配置文件位于项目根目录的`conf/app.ini`。

```ini
[common]
HttpPort = 6000
RPCPort = 7000
# 是否集群,单机则设为false
Cluster = true
# 对称加密key
CryptoKey = Adba723b7fe06819

[etcd]
Endpoints = 127.0.0.1:2379, 127.0.0.2:2379, 127.0.0.3:2379
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
    
    location /ws {
        proxy_pass http://ws_cluster; # 代理转发地址
        proxy_http_version 1.1;

        proxy_read_timeout 60s; # 超时设置

        # 启用支持websocket连接
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /api {
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header Host            $http_host;

        proxy_pass http://ws_cluster; # 代理转发地址
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
- [x] 关闭某个连接
- [x] 支持docker
- [ ] 查询某个客户端是否在线



## 沟通交流
QQ群：1028314856

<a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=adbd35a10f2c25bd765a7830d17579f6e645634d22151b6d480d1b05026125e9">
<img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="Golang websocket分布式交流群" title="Golang websocket分布式交流群"></a>