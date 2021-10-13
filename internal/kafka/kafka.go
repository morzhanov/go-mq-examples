package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type kfka struct {
	conn     *kafka.Conn
	kafkaUri string
	topic    string
}

type Kafka interface {
	CreateReader(groupId string) *kafka.Reader
	Conn() *kafka.Conn
}

func (k *kfka) createTopic() error {
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             k.topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	return k.conn.CreateTopics(topicConfigs...)
}

func (k *kfka) CreateReader(groupId string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.kafkaUri},
		Topic:    k.topic,
		GroupID:  groupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

func (k *kfka) Conn() *kafka.Conn {
	return k.conn
}

func NewKafka(uri string, topic string) (res Kafka, err error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", uri, topic, 0)
	if err != nil {
		return nil, err
	}
	msgQ := kfka{
		conn,
		uri,
		topic,
	}
	if err := msgQ.createTopic(); err != nil {
		return nil, err
	}
	return &msgQ, nil
}
