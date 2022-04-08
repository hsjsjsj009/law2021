package main

import (
	"LabAMQP/package/amqpPkg"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

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
	body := "Hello World!"
	number := 0
	for {
		number += 1
		err = channel.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body: []byte(fmt.Sprintf("%s%d",body,number)),
			})
		if err != nil {
			log.Fatal(err.Error())
		}
		time.Sleep(time.Millisecond*500)
	}

}
