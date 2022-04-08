package usecase

import (
	"MiniLab2/constant"
	"MiniLab2/package/amqpPkg"
)

func PushToExchange(data interface{}, uid string, broker amqpPkg.IBroker) error {
	return broker.PushExchange(data, constant.ExchangeName, constant.ExchangeType, uid)
}

