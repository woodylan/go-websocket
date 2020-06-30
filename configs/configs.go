package configs

import (
	"github.com/astaxie/beego/config"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	Conf *configs
)

type configs struct {
	CommonConf *commonConf
	EtcdConf   *etcdConf

	ServerList     map[string]string
	ServerListLock sync.RWMutex
	EtcdEndpoints  []string
}

type commonConf struct {
	IscCluster bool   //是否集群部署
	CryptoKey  string //对称加密key
	LocalHost  string //本机内网IP
	Port       string //监听端口号
	RPCPort    string //RPC端口号
}

type etcdConf struct {
	Host string
}

func InitConfig() (err error) {
	Conf = defaultConf()

	lasTwoPath := map[string]bool{
		"readconfig":    true,
		"send2client":   true,
		"bind2group":    true,
		"send2group":    true,
		"getonlinelist": true,
		"register":      true,
		"closeclient":   true,
	}

	path, _ := os.Getwd()
	if strings.Contains(path, "servers") {
		path += "/.."
	} else {
		for key := range lasTwoPath {
			if strings.Contains(path, key) {
				path += "/../.."
				break
			}
		}
	}
	configData, err := config.NewConfig("ini", path+"/configs/config.ini")
	if err != nil {
		return err
	}

	Conf.CommonConf.IscCluster, err = configData.Bool("common::cluster")
	if err != nil {
		return err
	}

	Conf.CommonConf.CryptoKey = configData.String("common::crypto_key")

	Conf.EtcdConf.Host = configData.String("etcd::host")
	if len(Conf.EtcdConf.Host) > 0 {
		Conf.EtcdEndpoints = append(Conf.EtcdEndpoints, Conf.EtcdConf.Host)
	}

	return err
}

func defaultConf() *configs {
	port := getPort()
	return &configs{
		CommonConf: &commonConf{
			IscCluster: false,
			CryptoKey:  "Adba723b7fe06819",
			LocalHost:  getIntranetIp(),
			Port:       port,
			RPCPort:    genRpcPort(port),
		},
		EtcdConf: &etcdConf{
			Host: "127.0.0.1:2379",
		},
		ServerList:    make(map[string]string),
		EtcdEndpoints: make([]string, 0),
	}
}

func getIntranetIp() string {
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

func getPort() string {
	port := "666"

	args := os.Args //获取用户输入的所有参数
	if len(args) >= 2 && len(args[1]) != 0 {
		port = args[1]
	}

	return port
}

//生成RPC通信端口号，目前是ws端口号+1000
func genRpcPort(port string) string {
	iPort, _ := strconv.Atoi(port)
	return strconv.Itoa(iPort + 1000)
}
