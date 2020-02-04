package servers

import (
	"encoding/json"
	"fmt"
	"go-websocket/define"
	"go-websocket/pkg/rabbitmq"
	"go-websocket/tools/readconfig"
	"go-websocket/tools/util"
	"log"
)

//RabbitMQ 实例
var rabbitMQ *rabbitmq.RabbitMQ

//初始化rabbitMQ
func InitRabbitMQ() {
	//如果是集群，则读取初始化RabbitMQ实例
	if util.IsCluster() {
		rabbitMQ = rabbitmq.NewRabbitMQPubSub(
			readconfig.ConfigData.String("rabbitMQ::amqpurl"),
			readconfig.ConfigData.String("rabbitMQ::exchange"))
		initRabbitMQReceive()
	}
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
				if publishMessage.Type == define.RPC_MESSAGE_TYPE_GROUP {
					//发送到指定分组
					Manager.SendMessage2LocalGroup(&publishMessage.SystemId, &publishMessage.GroupName, publishMessage.Code, publishMessage.Msg, &publishMessage.Data)
				} else if publishMessage.Type == define.RPC_MESSAGE_TYPE_SYSTEM {
					Manager.SendMessage2LocalSystem(&publishMessage.SystemId, publishMessage.Code, publishMessage.Msg, &publishMessage.Data)
				}

			} else {
				fmt.Println(err)
			}
		}
	}()
}
