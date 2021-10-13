package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/morzhanov/go-mq-examples/internal/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rmq struct {
	ch    *amqp.Channel
	qName string
}

func (r *rmq) Produce() {
	i := 0
	for {
		pub := amqp.Publishing{
			ContentType: "text/plain",
			MessageId:   fmt.Sprintf("__rabbitmq_msg_id_%d__", i),
			Body:        []byte(fmt.Sprintf("__rabbitmq_msg_value_%d__", i)),
		}
		i++
		if err :=
			r.ch.Publish("", r.qName, false, false, pub); err != nil {
			fmt.Println(fmt.Errorf("error in rabbitmq producer: %w", err))
		}
		time.Sleep(time.Second * 5)
	}
}

func (r *rmq) Listen() {
	msgs, err := r.ch.Consume(
		r.qName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(fmt.Errorf("error in rabbitmq listener: %w", err))
	}
	for d := range msgs {
		log.Println(fmt.Sprintf(`rabbitmq message received MID: "%s", BODY: "%s"`, d.MessageId, d.Body))
	}
}

func (r *rmq) createQueue() (amqp.Queue, error) {
	return r.ch.QueueDeclare(
		r.qName,
		false,
		false,
		false,
		false,
		nil,
	)
}

func NewRabbitMQ(uri string, queue string) (res mq.MQ, close func(), err error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	r := &rmq{ch, queue}
	if _, err := r.createQueue(); err != nil {
		return nil, nil, err
	}
	return r, func() {
		closeMsg := fmt.Errorf("error during rabbitmq connection closing: %w", err)
		if err := ch.Close(); err != nil {
			log.Println(closeMsg)
		}
		if err := conn.Close(); err != nil {
			log.Println(closeMsg)
		}
	}, nil
}
