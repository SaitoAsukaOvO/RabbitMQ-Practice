package main

import (
	"RabbitMQ/utils"
	"log"

amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ch := utils.GetChannels()

	/**
	declare a queue to send msg
		1. name <string> queue name
		2. durable <bool> if the msg in the queue stored in disk/memory
		3. exclusive <bool> if the msg in the queue can be shared by multiple consumers
		4. autoDelete <bool> if the connection auto deleted after finishing
		5. others
	*/
	q, err := ch.QueueDeclare(
		utils.HelloQueue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	body := "Hello World!"
	/**
	use channel to send msg
		1. send to which exchange
		2. determine the route key (currently it is queue name)
	 */
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	utils.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
