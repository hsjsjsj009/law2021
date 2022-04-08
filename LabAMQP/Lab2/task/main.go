package main

import (
	"LabAMQP/package/amqpPkg"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func main() {
	connData := amqpPkg.Connection{URL: "amqp://guest:guest@localhost:5672/"}
	conn,channel,err := connData.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		channel.Close()
		conn.Close()
	}()

	q,err := channel.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	body := bodyFrom(os.Args)
	err = channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body: []byte(body),
		})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf(" [x] Sent %s", body)
}
