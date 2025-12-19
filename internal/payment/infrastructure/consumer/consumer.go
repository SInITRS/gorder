package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SInITRS/gorder/common/broker"
	"github.com/SInITRS/gorder/common/genproto/orderpb"
	"github.com/SInITRS/gorder/payment/app"
	"github.com/SInITRS/gorder/payment/app/command"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

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
			c.handleMessage(ch, msg, q)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
	logrus.Infof("Received a message from %s, msg: %s", q.Name, msg.Body)
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbitmq")
	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()
	var err error
	defer func() {
		if err != nil {
			_ = msg.Nack(false, false)
		} else {
			_ = msg.Ack(false)
		}
	}()
	o := &orderpb.Order{}
	if err = json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("failed to unmashal msg to order, err=%v", err)
		return
	}

	if _, err = c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
		logrus.Infof("failed to create payment, err=%v", err)
		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
			logrus.Warnf("retry_error, error handling retry, orderID = %s,messageID = %s, err = %v", o.ID, msg.MessageId, err)
		}
		return
	}
	span.AddEvent("payment.created")
	logrus.Infof("consume EnventOrderCreated success!, orderID = %s, messageID = %s", o.ID, msg.MessageId)
}
