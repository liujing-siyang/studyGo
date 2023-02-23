package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://admin:123456@101.37.25.231:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	// 声明队列
	q, err := ch.QueueDeclare(
		"letter_queue", // name
		true,           // durable //将队列标记为持久的
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		amqp.Table{
			"x-message-ttl":             3000,          // 消息过期时间,毫秒
			"x-max-length":              3,             // 队列最大长度,即队列中可以存储处于ready状态消息的数量
			"x-max-length-bytes":        1024,          // 消息最大长度，队列中可以存储处于ready状态消息占用内存的大小(只计算消息体的字节数，不计算消息头、消息属性占用的字节数)
			"x-dead-letter-exchange":    "dead_direct", // 指定死信交换机
			"x-dead-letter-routing-key": "timeout",     // 指定死信routing-key
		}, // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 声明交换机
	err = ch.ExchangeDeclare(
		"letter_direct", // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	// 队列绑定（将队列、routing-key、交换机三者绑定到一起）
	err = ch.QueueBind(
		q.Name,          // queue name
		"info",          // routing key
		"letter_direct", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	// 声明死信队列
	dq, err := ch.QueueDeclare(
		"dead_queue", // name
		true,         // durable //将队列标记为持久的
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a dead queue")

	// 声明交换机
	err = ch.ExchangeDeclare(
		"dead_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an dead exchange")
	// 队列绑定（将队列、routing-key、交换机三者绑定到一起）
	err = ch.QueueBind(
		dq.Name,       // queue name
		"timeout",     // routing key
		"dead_direct", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(ctx,
		"letter_direct",       // exchange
		severityFrom(os.Args), // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}
