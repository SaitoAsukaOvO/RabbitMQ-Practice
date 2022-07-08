package utils

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"strings"
)

const (
	HelloQueue = "hello"
	TaskQueue = "task_queue"
	LogQueue = "logs"
	LogQueueDirect = "logs_direct"
)

func GetChannels() *amqp.Channel {
	//create connection to mq (set host username and pwd)
	conn, err := amqp.Dial("amqp://dyh:dyh111@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	// a connection contains several channels, get channel to connect queue
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	//defer ch.Close()
	return ch
}

func BodyFrom(args []string) string{
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func SeverityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}
