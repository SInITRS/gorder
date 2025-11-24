package consumer

import (
	"context"
	"encoding/json"

	"github.com/SInITRS/gorder/common/broker"
	"github.com/SInITRS/gorder/common/genproto/orderpb"
	"github.com/SInITRS/gorder/payment/app"
	"github.com/SInITRS/gorder/payment/app/command"
	"github.com/sirupsen/logrus"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{
		app: app,
	}
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
	logrus.Infof("Received a message from %s, msg: %s", q.Name, msg.Body)
	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("failed to unmashal msg to order, err=%v", err)
		_ = msg.Nack(false, false)
		return
	}

	if _, err := c.app.Commands.CreatePayment.Handle(context.TODO(), command.CreatePayment{Order: o}); err != nil {
		logrus.Infof("failed to create order, err=%v", err)
		_ = msg.Nack(false, false)
		return
	}
	_ = msg.Ack(false)
	logrus.Info("consume EnventOrderCreated success")
}
