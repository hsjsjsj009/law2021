package usecase

import (
	"MiniProject/constant"
	"MiniProject/package/amqpPkg"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func DownloadFile(
	urlData *url.URL,
	routingKey string,
	channel *amqp.Channel,
	broker amqpPkg.IBroker) {
	err := channel.ExchangeDeclare(
		constant.ExchangeName,
		constant.ExchangeType,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}
	time.Sleep(1000 * time.Millisecond)
	parsedPath := strings.Split(urlData.Path,"/")
	filename := fmt.Sprintf("%s-%s",routingKey,parsedPath[len(parsedPath)-1])

	resp,err := http.Get(urlData.String())
	if err != nil {
		_ = PushToExchange(fiber.Map{
			"type": "error",
		}, routingKey,broker)
		fmt.Println(err.Error())
		return
	}

	newFile, err := os.Create(fmt.Sprintf("./static/%s",filename))
	if err != nil {
		_ = PushToExchange(fiber.Map{
			"type": "error",
		}, routingKey,broker)
		fmt.Println(err.Error())
		return
	}
	defer newFile.Close()

	fileSize := resp.ContentLength
	body := resp.Body

	buffer := make([]byte, 32*1024)
	amountRead := 0
	updatePercentage := float64(0)
	_ = PushToExchange(fiber.Map{
		"type":  "progress",
		"value": 0,
	}, routingKey,broker)
	for {
		n, err := body.Read(buffer)
		if err != nil {
			if err == io.EOF {
				_ = PushToExchange(fiber.Map{
					"type": "finished",
					"url":  fmt.Sprintf("/download/%s", filename),
				}, routingKey,broker)
				return
			}
			_ = PushToExchange(fiber.Map{
				"type": "error",
			}, routingKey,broker)
			return
		}
		if n == 0 {
			_ = PushToExchange(fiber.Map{
				"type": "finished",
				"url":  fmt.Sprintf("/download/%s", filename),
			}, routingKey,broker)
			return
		}
		_, err = newFile.Write(buffer[0:n])
		if err != nil {
			_ = PushToExchange(fiber.Map{
				"type": "error",
			}, routingKey,broker)
			fmt.Println(err.Error())
			return
		}
		amountRead += n
		percentage := (float64(amountRead) / float64(fileSize)) * 100
		tempUpdateCalc := math.Floor(percentage / 5)
		if tempUpdateCalc > updatePercentage || percentage == 100 {
			updatePercentage = tempUpdateCalc
			_ = PushToExchange(fiber.Map{
				"type":  "progress",
				"value": percentage,
			}, routingKey, broker)
		}
	}
}

