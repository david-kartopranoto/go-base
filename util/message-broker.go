package util

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

//MessageBrokerService implements MessageBrokerService interface
type MessageBrokerService struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	queues map[string]amqp.Queue
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

	err = ch.Qos(1, 0, false)
	if err != nil {
		return nil, err
	}

	queues := make(map[string]amqp.Queue)
	for name, qConf := range config.MessageBroker.Queue {
		q, err := ch.QueueDeclare(
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
		queues[name] = q
	}

	s := &MessageBrokerService{
		conn:   conn,
		ch:     ch,
		queues: queues,
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

func (s *MessageBrokerService) Consume(queue string, stopChan chan bool, process func([]byte) error) error {
	q, ok := s.queues[queue]
	if !ok {
		return fmt.Errorf("Queue is not found")
	}

	qChannel, err := s.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range qChannel {
			log.Printf("Received a message: %s", d.Body)

			err := process(d.Body)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}
	}()

	return nil
}
