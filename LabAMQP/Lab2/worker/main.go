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
	err = channel.Qos(
		1,0,false)
	if err != nil {
		log.Fatal(err.Error())
	}
	msgs,err := channel.Consume(
		q.Name,
		"",
		false,
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
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
