package define

import "sync"

//本机内网IP
var LocalHost string

//监听端口号
var Port string

//RPC端口号
var RPCPort string

//服务器列表
var ServerList map[string]string
var ServerListLock sync.RWMutex

var EtcdEndpoints []string
