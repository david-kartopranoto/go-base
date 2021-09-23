package util

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

//MessageBrokerService implements MessageBrokerService interface
type MessageBrokerService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

//NewRabbitService create a new rabbit mq service
func NewRabbitMQService(config Config) (*MessageBrokerService, error) {
	host := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.MessageBroker.User,
		config.MessageBroker.Password,
		config.MessageBroker.Host,
		config.MessageBroker.Port)

	conn, err := amqp.Dial(host)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	for name, qConf := range config.MessageBroker.Queue {
		_, err = ch.QueueDeclare(
			name,
			qConf.Durable,
			qConf.DeleteUnused,
			qConf.Exclusive,
			qConf.NoWait,
			nil,
		)
		if err != nil {
			return nil, err
		}
	}

	s := &MessageBrokerService{
		conn: conn,
		ch:   ch,
	}

	return s, nil
}

func (s *MessageBrokerService) Publish(queue string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        b,
	}

	return s.ch.Publish(
		"",
		queue,
		false,
		false,
		message,
	)
}
