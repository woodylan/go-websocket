package src

import (
	"encoding/json"
	"fmt"
	"go-websocket/pkg/rabbitmq"
	"go-websocket/tools/readconfig"
	"log"
)

//RabbitMQ 实例
var rabbitMQ *rabbitmq.RabbitMQ

//创建rabbitMQ实例
func initRabbitMQ() {
	rabbitMQ = rabbitmq.NewRabbitMQPubSub(
		readconfig.ConfigData.String("rabbitMQ::amqpurl"),
		readconfig.ConfigData.String("rabbitMQ::exchange"))
}

func initRabbitMQReceive(b *binder) {
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
				SendMessage2LocalGroup(publishMessage.ObjectId, publishMessage.Message)
			} else {
				fmt.Println(err)
			}
		}
	}()
}
