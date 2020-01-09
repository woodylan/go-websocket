package define

import "sync"

//本机内网IP
var LocalHost string

//RPC端口号
var RPCPort string

//客户端的分组列表
var ClientGroupsMapMu sync.RWMutex
var ClientGroupsMap map[string][]string
