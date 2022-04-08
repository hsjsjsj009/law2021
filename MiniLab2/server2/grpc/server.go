package grpc

import (
	"MiniLab2/compression"
	"MiniLab2/package/amqpPkg"
	"MiniLab2/server2/usecase"
	"bytes"
	"context"
	"github.com/streadway/amqp"
	"strings"
)

type server struct {
	compression.UnimplementedCompressionServer
	channel *amqp.Channel
	broker amqpPkg.IBroker
}

func NewServer(channel *amqp.Channel, broker amqpPkg.IBroker) *server {
	return &server{channel: channel, broker: broker}
}

func (s server) CompressFile(_ context.Context,req *compression.CompressionRequest) (*compression.CompressionResponse, error)  {
	file := bytes.NewBuffer(req.FileBytes)
	fileSize := len(req.FileBytes)
	fileName := strings.Split(req.FileName,".")
	go usecase.CompressFile(file, int64(fileSize),[]string{
		strings.Join(fileName[0:len(fileName)-1],""),
		fileName[len(fileName)-1],
	},req.RoutingKey,s.channel,s.broker)
	return &compression.CompressionResponse{Success: compression.SuccessStatus_SUCCESS_STATUS_SUCCESS},nil
}