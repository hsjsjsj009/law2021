package grpc

import (
	"MiniProject/constant"
	"MiniProject/downloader"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func InitClient() (*grpc.ClientConn,downloader.MiniProjectServiceClient) {
	conn,err := grpc.Dial(fmt.Sprintf("localhost%s",constant.GRPCPort),
		grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn,downloader.NewMiniProjectServiceClient(conn)
}
