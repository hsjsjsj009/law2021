package amqpPkg

import (
	"encoding/json"
	"sync"

	"github.com/streadway/amqp"
)

// IBroker ...
type IBroker interface {
	PushQueue(data map[string]interface{}, types string) error
	PushQueueReconnect(url string, data map[string]interface{}, types, deadLetterKey string) (*amqp.Connection, *amqp.Channel, error)
	PushExchange(data interface{}, exchangeName,exchangeType,routingKey string) error
}

// broker ...
type broker struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Mux sync.Mutex
}

// NewBroker ...
func NewBroker(conn *amqp.Connection, channel *amqp.Channel) IBroker {
	return &broker{
		Connection: conn,
		Channel:    channel,
	}
}

func (m broker) PushExchange(data interface{}, exchangeName,exchangeType,routingKey string) error {
	m.Mux.Lock()
	defer m.Mux.Unlock()
	err := m.Channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		false,
		false,
		false,
		false,
		nil,
		)
	if err != nil {
		return err
	}
	dataByte,err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = m.Channel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: dataByte,
		},
		)
	return err
}

// PushQueue ...
func (m broker) PushQueue(data map[string]interface{}, types string) error {
	queue, err := m.Channel.QueueDeclare(types, true, false, false, false, nil)
	if err != nil {
		return err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = m.Channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	return err
}

// PushQueueReconnect ...
func (m broker) PushQueueReconnect(url string, data map[string]interface{}, types, deadLetterKey string) (*amqp.Connection, *amqp.Channel, error) {
	if m.Connection != nil {
		if m.Connection.IsClosed() {
			c := Connection{
				URL: url,
			}
			newConn, newChannel, err := c.Connect()
			if err != nil {
				return nil, nil, err
			}
			m.Connection = newConn
			m.Channel = newChannel
		}
	} else {
		c := Connection{
			URL: url,
		}
		newConn, newChannel, err := c.Connect()
		if err != nil {
			return nil, nil, err
		}
		m.Connection = newConn
		m.Channel = newChannel
	}

	args := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": deadLetterKey,
	}
	queue, err := m.Channel.QueueDeclare(types, true, false, false, false, args)
	if err != nil {
		return nil, nil, err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, nil, nil
	}

	err = m.Channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	return m.Connection, m.Channel, err
}
