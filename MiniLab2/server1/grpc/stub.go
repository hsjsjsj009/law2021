package grpc

import (
	"MiniLab2/compression"
	"MiniLab2/constant"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func InitClient() (*grpc.ClientConn,compression.CompressionClient) {
	conn,err := grpc.Dial(fmt.Sprintf("localhost%s",constant.GRPCPort),
		grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn,compression.NewCompressionClient(conn)
}
