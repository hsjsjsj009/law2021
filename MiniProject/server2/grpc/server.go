package grpc

import (
	"MiniProject/downloader"
	"MiniProject/package/amqpPkg"
	"MiniProject/server2/usecase"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"net/url"
)

type server struct {
	downloader.UnimplementedMiniProjectServiceServer
	channel *amqp.Channel
	broker amqpPkg.IBroker
}

func NewServer(channel *amqp.Channel, broker amqpPkg.IBroker) *server {
	return &server{channel: channel, broker: broker}
}

func (s server) Download(_ context.Context,req *downloader.DownloadRequest) (*downloader.DownloadResponse, error)  {
	uid := fmt.Sprintf("%s",uuid.New())
	urlData,_ := url.Parse(req.Url)
	go usecase.DownloadFile(urlData,uid,s.channel,s.broker)
	return &downloader.DownloadResponse{
		Url:    req.Url,
		UniqId: uid,
	},nil
}