package activemq

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-stomp/stomp/v3"

	"github.com/morzhanov/go-mq-examples/internal/mq"
)

type amq struct {
	conn  *stomp.Conn
	qName string
}

func (a *amq) Produce() {
	i := 0
	for {
		err := a.conn.Send(a.qName, "text/plain", []byte(fmt.Sprintf("__activemq_msg_value_%d__", i)))
		if err != nil {
			fmt.Println(fmt.Errorf("error in activemq producer: %w", err))
		}
		i++
		time.Sleep(time.Second * 5)
	}
}

func (a *amq) Listen() {
	sub, err := a.conn.Subscribe(a.qName, stomp.AckAuto)
	if err != nil {
		fmt.Println(fmt.Errorf("error in activemq listener: %w", err))
	}
	for {
		msg := <-sub.C
		log.Println(fmt.Sprintf(`activemq message received BODY: "%s"`, msg.Body))
	}
}

func NewActiveMQ(uri string, queue string) (res mq.MQ, close func(), err error) {
	netConn, err := net.DialTimeout("tcp", uri, 10*time.Second)
	if err != nil {
		return nil, nil, err
	}
	stompConn, err := stomp.Connect(netConn)
	if err != nil {
		return nil, nil, err
	}

	return &amq{stompConn, queue}, func() {
		errMsg := fmt.Errorf("error during activemq connection closing: %w", err)
		if err := stompConn.Disconnect(); err != nil {
			log.Println(errMsg)
		}
		if err := netConn.Close(); err != nil {
			log.Println(errMsg)
		}
	}, nil
}
