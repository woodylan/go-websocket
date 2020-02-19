package util

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"go-websocket/define"
	"go-websocket/tools/crypto"
	"go-websocket/tools/readconfig"
	"net"
	"strconv"
	"strings"
)

//GenUUID 生成uuid
func GenUUID() string {
	uuidFunc := uuid.NewV4()
	uuidStr := uuidFunc.String()
	uuidStr = strings.Replace(uuidStr, "-", "", -1)
	uuidByt := []rune(uuidStr)
	return string(uuidByt[8:24])
}

//对称加密IP和端口，当做clientId
func GenClientId() string {
	raw := []byte(define.LocalHost + ":" + define.RPCPort)
	key := readconfig.ConfigData.String("common::crypto_key")
	str, err := crypto.Encrypt(raw, []byte(key))
	if err != nil {
		panic(err)
	}

	return str
}

//解析redis的地址格式
func ParseRedisAddrValue(redisValue string) (host string, port string, err error) {
	if redisValue == "" {
		err = errors.New("解析地址错误")
		return
	}
	addr := strings.Split(redisValue, ":")
	if len(addr) != 2 {
		err = errors.New("解析地址错误")
		return
	}
	host, port = addr[0], addr[1]

	return
}

//判断地址是否为本机
func IsAddrLocal(host string, port string) bool {
	return host == define.LocalHost && port == define.RPCPort
}

//是否集群
func IsCluster() bool {
	cluster, _ := readconfig.ConfigData.Bool("common::cluster")
	return cluster
}

//生成RPC通信端口号，目前是ws端口号+1000
func GenRpcPort(port string) string {
	iPort, _ := strconv.Atoi(port)
	return strconv.Itoa(iPort + 1000)
}

//获取group分组key
func GetGroupKey(groupName string) string {
	//在redis每个服务都有一个单独的key，用服务器地址(ip:port)区分
	return define.REDIS_KEY_GROUP + define.LocalHost + ":" + define.RPCPort + ":" + groupName
}

//获取client key地址信息
func GetAddrInfoAndIsLocal(clientId string) (addr string, host string, port string, isLocal bool, err error) {
	//解密ClientId
	key := readconfig.ConfigData.String("common::crypto_key")

	addr, err = crypto.Decrypt(clientId, []byte(key))
	if err != nil {
		return
	}

	host, port, err = ParseRedisAddrValue(addr)
	if err != nil {
		return
	}

	isLocal = IsAddrLocal(host, port)
	return
}

//获取本机内网IP
func GetIntranetIp() string {
	addrs, _ := net.InterfaceAddrs()

	for _, addr := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}

	return ""
}

func GenGroupKey(systemId, groupName string) string {
	return systemId + ":" + groupName
}
