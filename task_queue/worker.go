package main

import (
	"RabbitMQ-Practice/utils"
	"bytes"
	"log"
	"time"
)



func main() {
	ch := utils.GetChannels()
	q, err := ch.QueueDeclare(
		utils.TaskQueue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	/*
	This consumer receives at most messages for processing at a time,
	and only after the message processing is completed and a manual reply is received,
	will it be distributed.
	*/
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FailOnError(err, "Failed to set QoS")

	// a worker to receive msg
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")


	forever := make(chan bool)

	/*
	Message acknowledgment:
	make sure a message is never lost
	we will use manual message acknowledgements by passing a false for the "auto-ack" argument
	and then send a proper acknowledgment from the worker with d.Ack(false)
	(this acknowledges a single delivery), once we're done with a task.
	When multiple is true, this delivery and all prior unacknowledged deliveries
	on the same channel will be acknowledged.
	 */
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
