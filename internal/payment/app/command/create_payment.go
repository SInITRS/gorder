package command

import (
	"context"

	"github.com/SInITRS/gorder/common/decorator"
	"github.com/SInITRS/gorder/common/genproto/orderpb"
	"github.com/SInITRS/gorder/payment/domain"
	"github.com/sirupsen/logrus"
)

type CreatePayment struct {
	Order *orderpb.Order
}

type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]

type createPaymentHandler struct {
	processer domain.Processor
	orderGRPC OrderService
}

// Handle implements decorator.CommandHandler.
func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
	link, err := c.processer.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}
	logrus.Infof("create payment link for order : %s success, payment link is %s", cmd.Order.ID, link)
	newOrder := &orderpb.Order{
		ID:          cmd.Order.ID,
		CustomerID:  cmd.Order.CustomerID,
		Status:      "waiting_for_payment",
		Items:       cmd.Order.Items,
		PaymentLink: link,
	}
	err = c.orderGRPC.UpdateOrder(ctx, newOrder)
	return link, err
}

func NewCreatePaymentHandler(
	processer domain.Processor,
	orderGRPC OrderService,
	logger *logrus.Entry,
	client decorator.MetricsClient,
) CreatePaymentHandler {
	return decorator.ApplyCommandDecorators(
		createPaymentHandler{processer: processer, orderGRPC: orderGRPC},
		logger,
		client,
	)
}
