# Golang实现的分布式WebSocket微服务

## 简介

本系统基于Golang、Redis、RabbitMQ、RPC实现分布式WebSocket微服务，也可以单机部署，单机部署不需要Redis、RabbitMQ和RPC。分布式部署可以支持nginx负责均衡、水平扩容部署，程序之间使用RabbitMQ广播、RPC通信。

目前实现的功能有，给指定客户端发送消息、绑定客户端到分组、给分组里的客户端批量发送消息。适用于长连接的大部分场景，分组可以理解为聊天室，绑定客户端到分组相当于把客户端添加到聊天室，给分组发送信息相当于给聊天室的每个人发送消息。



## 使用

### 下载本项目

```shell
git clone https://github.com/woodylan/go-websocket.git
```

### 编译

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

### 执行

编译成功之后会得到一个二进制文件`go-websocket`，执行该二进制文件，文件名后面跟着的是端口号，下面的命令`666`则表示端口号，你可以可以改成其他的。

```shell
./go-websocket 666
```

### 连接成功

打开支持Websocket的客户端，输入 `ws://127.0.0.1:666/ws` 进行连接，连接成功会返回`clientId`。



## 单机部署

单机部署很简单，不需要配置Redis、RabbitMQ，只需要编译然后运行该二进制文件就可以了，步骤如上。



## 分布式部署

## 安装Redis

参考网上教程

### 安装RabbitMQ

参考网上教程

### 配置文件

配置文件位于项目根目录的`configs/config.ini`，`cluster`为true表示分布式部署，分布式部署[rabbitMQ]、[redis]项必填。

```ini
[common]
# 是否分布式部署
cluster = true

[rabbitMQ]
# AMQP协议URL 格式：amqp://user:password@host:port/vhostname
amqpurl = 
# RabbitMQ交换机名称
exchange = 

[redis]
host = 127.0.0.1
port = 6379
password =
```

### 运行项目

在不同的机器运行本项目，注意配置号端口号，项目如果在同一机器，则必须用不同的端口。你可以用`supervisor`做进程管理。

### 配置Nginx负载均衡

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



## 接口

### 连接接口

**请求地址：**/ws

**协议：**websocket

**响应示例：**

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "clientId": "9fa54bdbbf2778cb"
  }
}
```

### 发送信息给指定客户端

**请求地址：**/api/send_to_client

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| clientId | string | 是       | 客户端ID |
| code | integer | 是       | 自定义的状态码 |
| msg | string | 是       | 自定义的状态消息 |
| data | sring、array、object | 是       | 消息内容 |

**响应示例：**

```json
{
  "code": 0,
  "msg": "success",
  "data": []
}
```

### 绑定客户端到分组

**请求地址：**/api/bind_to_group

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| clientId | string | 是       | 客户端ID |
| groupName | string | 是       | 分组名 |

**响应示例：**

```json
{
  "code": 0,
  "msg": "success",
  "data": []
}
```

### 发送信息给指定分组

**请求地址：**/api/send_to_group

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| groupName | string | 是       | 分组名 |
| code | integer | 是       | 自定义的状态码 |
| msg | string | 是       | 自定义的状态消息 |
| data | sring、array、object | 是       | 消息内容 |

**响应示例：**

```json
{
  "code": 0,
  "msg": "success",
  "data": []
}
```



## 实现的功能

- [x] 分布式
- [ ] 账户授权模式
- [ ] 日志
- [ ] 用户绑定客户端功能
- [ ] 消息持久化
- [ ] 消息记录
- [ ] 限流
- [ ] 黑名单