package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/morzhanov/go-mq-examples/internal/mq"

	"github.com/segmentio/kafka-go"
)

type kmq struct {
	conn     *kafka.Conn
	kafkaUri string
	topic    string
}

func (k *kmq) createTopic() error {
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             k.topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	return k.conn.CreateTopics(topicConfigs...)
}

func (k *kmq) createReader(groupId string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.kafkaUri},
		Topic:    k.topic,
		GroupID:  groupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

func (k *kmq) Produce() {
	i := 0
	for {
		m := kafka.Message{
			Key:   []byte(fmt.Sprintf("__kafka_msg_key_%d__", i)),
			Value: []byte(fmt.Sprintf("__kafka_msg_value_%d__", i)),
		}
		i++
		if _, err := k.conn.WriteMessages(m); err != nil {
			fmt.Println(fmt.Errorf("error in kafka producer: %w", err))
		}
		time.Sleep(time.Second * 5)
	}
}

func (k *kmq) Listen() {
	r := k.createReader(k.topic)
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(fmt.Errorf("error in kafka listener: %w", err))
			continue
		}
		log.Println(fmt.Sprintf(`kafka message received K: "%s", V: "%s"`, m.Key, m.Value))
	}
}

func NewKafka(uri string, topic string) (res mq.MQ, close func(), err error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", uri, topic, 0)
	if err != nil {
		return nil, nil, err
	}
	k := kmq{conn, uri, topic}
	if err := k.createTopic(); err != nil {
		return nil, nil, err
	}
	return &k, func() {
		if err := conn.Close(); err != nil {
			log.Println(fmt.Errorf("error during kafka connection closing: %w", err))
		}
	}, nil
}
