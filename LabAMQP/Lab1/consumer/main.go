package main

import (
	"LabAMQP/package/amqpPkg"
	"log"
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
	msgs,err := channel.Consume(
		"hello",
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
