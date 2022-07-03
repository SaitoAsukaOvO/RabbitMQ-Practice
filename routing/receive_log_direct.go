package main

import (
	"RabbitMQ/utils"
	"log"
	"os"
)


func main () {
	ch := utils.GetChannels()
	err := ch.ExchangeDeclare(utils.LogQueueDirect, "direct", true, false, false, false, nil)
	utils.FailOnError(err, "fail to declare exchange")

	//create a new binding for each severity we're interested in.
	q, err := ch.QueueDeclare("", true, false, false, false, nil)

	for _, s := range os.Args {
		err = ch.QueueBind(q.Name, s, utils.LogQueueDirect, false, nil)
	}
	utils.FailOnError(err, "fail to bind queue")
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
