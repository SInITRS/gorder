package query

import (
	"context"

	"github.com/SInITRS/gorder/common/decorator"
	domain "github.com/SInITRS/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
)

type GetCustomerOrder struct {
	CustomerID string
	OrderID    string
}

type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

type getCustomerOrderHandler struct {
	orderRepo domain.Repository
}

func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
	o, err := g.orderRepo.GetOrder(ctx, query.CustomerID, query.OrderID)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func NewGetCustomerOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	client decorator.MetricsClient,
) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
		getCustomerOrderHandler{orderRepo: orderRepo},
		logger,
		client,
	)
}
