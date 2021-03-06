package main

import (
	"RabbitMQ-Practice/utils"
	"fmt"
	"log"
)


func main () {
	ch := utils.GetChannels()
	err := ch.ExchangeDeclare(utils.LogQueue, "fanout", true, false, false, false, nil)
	utils.FailOnError(err, "fail to declare exchange")
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	err = ch.QueueBind(q.Name, "", utils.LogQueue, false, nil)
	fmt.Println(q.Name)
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
