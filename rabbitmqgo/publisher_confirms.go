package main

import (
	"context"
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
	conn, err := amqp.Dial("amqp://admin:123456@101.37.25.231:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	confirms := make(chan amqp.Confirmation)
	ch.NotifyPublish(confirms)
	// broker回应消息是否到达exchange
	go func() {
		for confirm := range confirms {
			if confirm.Ack {
				// code when messages is confirmed
				ch.Ack(confirm.DeliveryTag, false)
				log.Printf("Confirmed")
			} else {
				// code when messages is nack-ed
				ch.Nack(confirm.DeliveryTag, false, false)
				log.Printf("Nacked")

			}
		}
	}()

	err = ch.Confirm(false)
	failOnError(err, "Failed to confirm")
	//broker回应消息是否正确入列
	returnConfirms := make(chan amqp.Return)
	ch.NotifyReturn(returnConfirms)
	go func() {
		for ret := range returnConfirms {
			if len(string(ret.Body)) != 0 {
				log.Println("message can be route to queue:", string(ret.Body))
			}
		}
	}()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	consume(ch, q.Name)
	publish(ch, q.Name, "hello")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	var forever chan struct{}
	<-forever
}

func consume(ch *amqp.Channel, qName string) {
	msgs, err := ch.Consume(
		qName, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
}

func publish(ch *amqp.Channel, qName, text string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := ch.PublishWithContext(
		ctx,
		"",
		qName,
		true,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(text),
		})
	failOnError(err, "Failed to publish a message")
}
