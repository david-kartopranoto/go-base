package util

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type MetricProvider interface {
	SaveHistogram(handler, method, statusCode string, duration float64)
}

//MessageBrokerService implements MessageBrokerService interface
type MessageBrokerService struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	queues map[string]amqp.Queue
	metric MetricProvider
}

//NewRabbitService create a new rabbit mq service
func NewRabbitMQService(config Config, metric MetricProvider) (*MessageBrokerService, error) {
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
		metric: metric,
	}

	return s, nil
}

func (s *MessageBrokerService) Publish(queue string, body interface{}) error {
	var err error
	start := time.Now()

	defer func(t time.Time, e error) {
		duration := time.Since(t).Seconds()
		s.metric.SaveHistogram("Publish", queue, fmt.Sprintln(e), duration)
	}(start, err)

	var b []byte
	b, err = json.Marshal(body)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        b,
	}

	err = s.ch.Publish(
		"",
		queue,
		false,
		false,
		message,
	)

	return err
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
			var errP error
			start := time.Now()

			log.Printf("Received a message: %s", d.Body)
			errP = process(d.Body)
			if errP != nil {
				log.Printf("Error decoding JSON: %s", errP)
			}
			duration := time.Since(start).Seconds()
			s.metric.SaveHistogram("Consume-Process", queue, fmt.Sprintf("%v", errP), duration)

			start = time.Now()
			if errP = d.Ack(false); errP != nil {
				log.Printf("Error acknowledging message : %s", errP)
			} else {
				log.Printf("Acknowledged message")
			}
			duration = time.Since(start).Seconds()
			s.metric.SaveHistogram("Consume-Ack", queue, fmt.Sprintf("%v", errP), duration)

		}
	}()

	return nil
}
