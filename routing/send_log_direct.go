package main

import (
	"RabbitMQ/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

//a message goes to the queues whose binding key exactly matches the routing key of the message.

func main() {
	ch := utils.GetChannels()
	err := ch.ExchangeDeclare(utils.LogQueueDirect, "direct", true, false, false, false, nil)
	utils.FailOnError(err, "fail to declare exchange")
	body := utils.BodyForm(os.Args)
	err = ch.Publish(utils.LogQueueDirect,  utils.SeverityFrom(os.Args), false, false,  amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	utils.FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
