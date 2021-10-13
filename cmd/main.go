package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/morzhanov/go-mq-examples/internal/rabbitmq"

	"github.com/morzhanov/go-mq-examples/internal/activemq"
	"github.com/morzhanov/go-mq-examples/internal/config"
	"github.com/morzhanov/go-mq-examples/internal/kafka"
)

const initErr = "initialization error in step %s: %w"

func failOnError(step string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf(initErr, step, err))
	}
}

func main() {
	c, e := config.NewConfig()
	failOnError("config", e)
	a, closeAMQ, e := activemq.NewActiveMQ(c.ActiveMQURI, c.ActiveMQQueue)
	failOnError("activemq", e)
	defer closeAMQ()
	k, closeKafka, e := kafka.NewKafka(c.KafkaURI, c.KafkaTopic)
	failOnError("kafka", e)
	defer closeKafka()
	r, closeRMQ, e := rabbitmq.NewRabbitMQ(c.RabbitMQURI, c.RabbitMQQueue)
	failOnError("rabbitmq", e)
	defer closeRMQ()

	go a.Produce()
	go k.Produce()
	go r.Produce()

	go a.Listen()
	go k.Listen()
	go r.Listen()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	log.Println("MQs successfully started!")
	<-quit
	log.Println("received os.Interrupt, exiting...")
}
