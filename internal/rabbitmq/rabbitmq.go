package rabbitmq

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type rmq struct {
	conn     *kafka.Conn
	kafkaUri string
	topic    string
}

type RabbitMQ interface {
	CreateReader(groupId string) *kafka.Reader
	Conn() *kafka.Conn
}

func (k *rmq) createTopic() error {
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             k.topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	return k.conn.CreateTopics(topicConfigs...)
}

func (k *rmq) CreateReader(groupId string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.kafkaUri},
		Topic:    k.topic,
		GroupID:  groupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

func (k *rmq) Conn() *kafka.Conn {
	return k.conn
}

func NewRabbitMQ(uri string, topic string) (res RabbitMQ, err error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", uri, topic, 0)
	if err != nil {
		return nil, err
	}
	msgQ := rmq{
		conn,
		uri,
		topic,
	}
	if err := msgQ.createTopic(); err != nil {
		return nil, err
	}
	return &msgQ, nil
}
