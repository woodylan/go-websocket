package readconfig

import (
	"errors"
	"github.com/astaxie/beego/config"
	"os"
	"strings"
)

var ConfigData config.Configer

func InitConfig() (err error) {
	path, _ := os.Getwd()
	if strings.Contains(path, "readconfig") {
		path += "/../.."
	} else if strings.Contains(path, "servers") {
		path += "/.."
	}
	ConfigData, err = config.NewConfig("ini", path+"/configs/config.ini")
	if err != nil {
		return err
	}

	cluster, err := ConfigData.Bool("common::cluster")
	if (err != nil) {
		return err
	}

	//如果设置了集群，则amqpurl和exchange必须填写
	if cluster {
		if len(ConfigData.String("rabbitMQ::amqpurl")) == 0 || len(ConfigData.String("rabbitMQ::exchange")) == 0 {
			return errors.New("集群配置不完整")
		}
	}
	return
}
