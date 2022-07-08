package main

import (
	"RabbitMQ-Practice/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)


func main() {
	ch := utils.GetChannels()
	err := ch.ExchangeDeclare(utils.LogQueue, "fanout", true, false, false, false, nil)
	utils.FailOnError(err, "fail to declare exchange")
	body := utils.BodyFrom(os.Args)
	err = ch.Publish(utils.LogQueue, "", false, false,  amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	utils.FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
