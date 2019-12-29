package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string //队列名称
	Exchange  string //交换机名称
	Key       string
	AMQPUrl   string //amqp协议的连接信息
}

//创建结构体
func NewRabbitMQ(queueName string, AMQPURL string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		AMQPUrl:   AMQPURL,
	}
}

//发布订阅模式创建实例
func NewRabbitMQPubSub(AMQPURL string, exchangeName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitMQ := NewRabbitMQ("", AMQPURL, exchangeName, "")

	var err error

	//获取connection
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.AMQPUrl)
	rabbitMQ.failOnErr(err, "failed to connect rabbitMQ!")

	//获取channel
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnErr(err, "failed to open a channel")

	return rabbitMQ
}

//生产消息
func (r *RabbitMQ) PublishPub(message string) {
	//1.申请交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout", //广播类型
		true,     //是否持久化
		false,    //是否自动删除
		false,
		false,
		nil,
	)

	r.failOnErr(err, "failed to declare an exchange")

	//2.发送消息
	_ = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

//消费消息
func (r *RabbitMQ) ReceiveSub() (<-chan amqp.Delivery, error) {
	//1.申请交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout", //广播类型
		true,     //是否持久化
		false,    //是否自动删除
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")

	//2.创建队列
	q, err := r.channel.QueueDeclare(
		"",    //随机生产队列名称
		false, //是否持久化
		false, //是否自动删除
		true,  //是否具有排他性
		false, //是否阻塞
		nil,   //额外属性
	)
	r.failOnErr(err, "failed to declare an queue")

	//3.绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		"", //发布订阅模式下，这里留空
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err, "failed to bind queue")

	//消费消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",    //区分多个多个消费者
		true,  //是否自动应答
		false, //是否具有排他性
		false,
		false, //是否阻塞
		nil,
	)

	return msgs, err
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func (r *RabbitMQ) Destroy() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}
