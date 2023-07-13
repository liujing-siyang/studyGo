package main

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://admin:admin123456@101.37.25.231:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_direct", "info")
	err = ch.QueueBind(
		q.Name,        // queue name
		"info",        // routing key
		"logs_direct", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
	// case1
	// 这种消费方式类似kafka的消费组，四个消费者消费同一个队列（相当于kafka的topic）,每个消费者消费一个分区
	// for i := 0; i < 3; i++ {
	// 	go consumermsg(ch, q.Name, i+1)
	// }
	// go consumermsg2(ch, q.Name, 4)

	// case2
	// 虽然test50晚了5s确认，但没有阻塞其它消息的消费
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	for d := range msgs {
		// msg := strings.TrimLeft(string(d.Body),"test")
		// id ,_ := strconv.ParseInt(msg,10,32)
		if string(d.Body) == "test50" {
			go func(msg []byte) {
				log.Printf("[x] %s", msg)
				time.Sleep(5000 * time.Millisecond)
				d.Ack(false)
			}(d.Body)
			continue
		}
		log.Printf("[x] %s", d.Body)
		d.Ack(false)
	}

	var forever chan struct{}

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func consumermsg(ch *amqp.Channel, queue string, id int) {
	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")
	for d := range msgs {
		log.Printf("[id] %d [x] %s", id, d.Body)
		d.Ack(false)
	}
}

func consumermsg2(ch *amqp.Channel, queue string, id int) {
	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")
	for d := range msgs {
		log.Printf("[id] %d [x] %s", id, d.Body)
		time.Sleep(2000 * time.Millisecond)
		d.Ack(false)
	}
}
