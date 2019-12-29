package readconfig

import (
	"github.com/astaxie/beego/config"
)

var ConfigData config.Configer

func ReadConfig() {
	var err error

	ConfigData, err = config.NewConfig("ini", "configs/config.ini")
	if err != nil {
		panic("读取配置文件错误")
	}

	cluster, err := ConfigData.Bool("common::cluster")
	if (err != nil) {
		panic("读取配置文件错误")
	}

	//如果设置了集群，则amqpurl和exchange必须填写
	if cluster {
		if len(ConfigData.String("rabbitMQ::amqpurl")) == 0 || len(ConfigData.String("rabbitMQ::exchange")) == 0 {
			panic("集群配置不完整")
		}
	}
}
