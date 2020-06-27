package util

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"go-websocket/configs"
	"go-websocket/tools/crypto"
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
	raw := []byte(configs.Conf.CommonConf.LocalHost + ":" + configs.Conf.CommonConf.RPCPort)
	str, err := crypto.Encrypt(raw, []byte(configs.Conf.CommonConf.CryptoKey))
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
	return host == configs.Conf.CommonConf.LocalHost && port == configs.Conf.CommonConf.RPCPort
}

//是否集群
func IsCluster() bool {
	return configs.Conf.CommonConf.IscCluster
}

//生成RPC通信端口号，目前是ws端口号+1000
func GenRpcPort(port string) string {
	iPort, _ := strconv.Atoi(port)
	return strconv.Itoa(iPort + 1000)
}

//获取client key地址信息
func GetAddrInfoAndIsLocal(clientId string) (addr string, host string, port string, isLocal bool, err error) {
	//解密ClientId
	addr, err = crypto.Decrypt(clientId, []byte(configs.Conf.CommonConf.CryptoKey))
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

func GenGroupKey(systemId, groupName string) string {
	return systemId + ":" + groupName
}
