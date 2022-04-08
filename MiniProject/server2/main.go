package main

import (
	"MiniProject/constant"
	"MiniProject/downloader"
	"MiniProject/package/amqpPkg"
	grpcServer "MiniProject/server2/grpc"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"log"
	"net"
)

//Track progress by write the byte manually through the compressor

var (
	broker  amqpPkg.IBroker
	channel *amqp.Channel
	conn    *amqp.Connection
)

func main() {
	var err error

	connData := amqpPkg.Connection{URL: "amqp://guest:guest@localhost:5672/"}
	conn, channel,err = connData.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	broker = amqpPkg.NewBroker(conn, channel)

	holdChannel := make(chan bool)

	go func() {
		lis, err := net.Listen("tcp", constant.GRPCPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*1024))
		downloader.RegisterMiniProjectServiceServer(s, grpcServer.NewServer(channel,broker))
		log.Printf("GRPC Server Listening on port %s",constant.GRPCPort)
		if err := s.Serve(lis); err != nil {
			log.Fatal(err.Error())
		}
	}()

	<-holdChannel
}