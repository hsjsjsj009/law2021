package main

import (
	"MiniLab2/compression"
	"MiniLab2/constant"
	"MiniLab2/package/amqpPkg"
	grpcServer "MiniLab2/server2/grpc"
	"MiniLab2/server2/usecase"
	"MiniLab2/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
)

//Track progress by write the byte manually through the compressor

var (
	broker  amqpPkg.IBroker
	channel *amqp.Channel
	conn    *amqp.Connection
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 1024*1024*1024,
	})

	var err error

	app.Use(recover.New())
	app.Static("/static","./static")

	connData := amqpPkg.Connection{URL: "amqp://guest:guest@localhost:5672/"}
	conn, channel,err = connData.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	broker = amqpPkg.NewBroker(conn, channel)

	app.Post("/", func(ctx *fiber.Ctx) error {
		routingKey := ctx.Get("X-ROUTING-KEY","")
		if routingKey == "" {
			return ctx.Status(400).JSON(fiber.Map{
				"message":"X-ROUTING-KEY header required",
			})
		}

		form,err := ctx.MultipartForm()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		files,err := util.MultipartFileConverter(form.File)

		if len(files["file"]) > 1 {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"only process one file at a time",
			})
		}

		file := files["file"][0]

		fileName := strings.Split(file.FileName,".")

		go usecase.CompressFile(file.File,file.Size,[]string{
			strings.Join(fileName[0:len(fileName)-1],""),
			fileName[len(fileName)-1],
		}, routingKey,channel,broker)

		return ctx.JSON(fiber.Map{
			"status":"ok",
		})
	})

	holdChannel := make(chan bool)

	go func() {
		if err := app.Listen(":8080"); err !=nil {
			log.Fatal(err.Error())
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", constant.GRPCPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*1024))
		compression.RegisterCompressionServer(s, grpcServer.NewServer(channel,broker))
		log.Printf("GRPC Server Listening on port %s",constant.GRPCPort)
		if err := s.Serve(lis); err != nil {
			log.Fatal(err.Error())
		}
	}()

	<-holdChannel
}