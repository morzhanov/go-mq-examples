package kafka_test

import (
	kfka "github.com/morzhanov/go-mq-examples/internal/kafka"
	"github.com/segmentio/kafka-go"
)

type KafkaMock struct {
	createReaderMock func(groupId string) *kafka.Reader
	connMock         func() *kafka.Conn
	kafkaUriMock     func() string
	topicMock        func() string
}

func (m *KafkaMock) CreateReader(groupId string) *kafka.Reader {
	return m.createReaderMock(groupId)
}
func (m *KafkaMock) Conn() *kafka.Conn {
	return m.connMock()
}

func NewMqMock(
	createReader func(groupId string) *kafka.Reader,
	conn func() *kafka.Conn,
	kafkaUri func() string,
	topic func() string,
) kfka.Kafka {
	k := KafkaMock{}
	if createReader != nil {
		k.createReaderMock = createReader
	} else {
		k.createReaderMock = func(groupId string) *kafka.Reader { return nil }
	}
	if conn != nil {
		k.connMock = conn
	} else {
		k.connMock = func() *kafka.Conn { return nil }
	}
	if kafkaUri != nil {
		k.kafkaUriMock = kafkaUri
	} else {
		k.kafkaUriMock = func() string { return "uri" }
	}
	if topic != nil {
		k.topicMock = topic
	} else {
		k.topicMock = func() string { return "topic" }
	}
	return &k
}
