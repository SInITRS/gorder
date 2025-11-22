package consumer

import (
	"github.com/SInITRS/gorder/common/broker"
	"github.com/sirupsen/logrus"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct{}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("failed to consume: quene = %s, err = %v", q.Name, err)
	}
	forever := make(chan struct{})
	go func() {
		for msg := range msgs {
			c.handleMessage(msg, q)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
	// TODO: 处理消息
	msg.Ack(false)
	logrus.Infof("Received a message from %s, msg: %s", q.Name, msg.Body)
}
