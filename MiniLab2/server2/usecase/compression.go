package usecase

import (
	"MiniLab2/constant"
	"MiniLab2/package/amqpPkg"
	"compress/gzip"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

func CompressFile(
	file io.Reader,
	fileSize int64,
	fileName []string,
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
	time.Sleep(500 * time.Millisecond)
	newFileName := fmt.Sprintf("%s-%d.gz", fileName[0], time.Now().UnixNano())
	newFile, err := os.Create("./static/" + newFileName)
	if err != nil {
		_ = PushToExchange(fiber.Map{
			"type": "error",
		}, routingKey,broker)
		fmt.Println(err.Error())
		return
	}
	defer newFile.Close()

	gz, err := gzip.NewWriterLevel(newFile, gzip.BestCompression)
	gz.Name = strings.Join(fileName, ".")
	if err != nil {
		_ = PushToExchange(fiber.Map{
			"type": "error",
		}, routingKey,broker)
		fmt.Println(err.Error())
		return
	}
	defer gz.Close()
	buffer := make([]byte, 32*1024)
	amountRead := 0
	updatePercentage := float64(0)
	_ = PushToExchange(fiber.Map{
		"type":  "progress",
		"value": 0,
	}, routingKey,broker)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				_ = PushToExchange(fiber.Map{
					"type": "finished",
					"url":  fmt.Sprintf("http://localhost:8080/static/%s", newFileName),
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
				"url":  fmt.Sprintf("http://localhost:8080/static/%s", newFileName),
			}, routingKey,broker)
			return
		}
		_, err = gz.Write(buffer)
		if err != nil {
			_ = PushToExchange(fiber.Map{
				"type": "error",
			}, routingKey,broker)
			fmt.Println(err.Error())
			return
		}
		amountRead += n
		percentage := (float64(amountRead) / float64(fileSize)) * 100
		tempUpdateCalc := math.Floor(percentage / 10)
		if tempUpdateCalc > updatePercentage || percentage == 100 {
			updatePercentage = tempUpdateCalc
			_ = PushToExchange(fiber.Map{
				"type":  "progress",
				"value": percentage,
			}, routingKey,broker)
		}
	}
}

