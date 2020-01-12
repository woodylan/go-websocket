package server

import (
	"encoding/json"
	"fmt"
	"go-websocket/pkg/rabbitmq"
	"go-websocket/tools/readconfig"
	"go-websocket/tools/util"
	"log"
)

//RabbitMQ 实例
var rabbitMQ *rabbitmq.RabbitMQ

//初始化rabbitMQ
func Init()  {
	//如果是集群，则读取初始化RabbitMQ实例
	if util.IsCluster() {
		initRabbitMQ()
		initRabbitMQReceive()
	}
}

//创建rabbitMQ实例
func initRabbitMQ() {
	rabbitMQ = rabbitmq.NewRabbitMQPubSub(
		readconfig.ConfigData.String("rabbitMQ::amqpurl"),
		readconfig.ConfigData.String("rabbitMQ::exchange"))
}

func initRabbitMQReceive() {
	msgs, err := rabbitMQ.ReceiveSub()
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for receiveData := range msgs {
			log.Printf("Received a message: %s", receiveData.Body)

			var publishMessage publishMessage
			err := json.Unmarshal([]byte(receiveData.Body), &publishMessage)
			if err == nil {
				//发送到指定分组
				SendMessage2LocalGroup(&publishMessage.ObjectId, &publishMessage.Message)
			} else {
				fmt.Println(err)
			}
		}
	}()
}
