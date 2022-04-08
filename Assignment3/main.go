package main

import (
	"assignment_3/package/amqpPkg"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main()  {
	connData := amqpPkg.Connection{URL: "amqp://guest:guest@localhost:5672/"}
	conn,channel,err := connData.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		channel.Close()
		conn.Close()
	}()

	err = channel.ExchangeDeclare(
		"x-exchange",
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		timeData := time.Now().Format(time.RFC822)
		jsonData,_ := json.Marshal(map[string]interface{}{
			"from":"bot",
			"data":timeData,
		})
		err = channel.Publish(
			"x-exchange",
			"",
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType: "application/json",
				Body: jsonData,
			})
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(timeData)
		time.Sleep(1*time.Minute)
	}
}