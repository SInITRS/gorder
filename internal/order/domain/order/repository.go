package order

import (
	"context"
	"fmt"
)

type Repository interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrder(ctx context.Context, id, customerID string) (*Order, error)
	UpdateOrder(
		ctx context.Context,
		order *Order,
		updateFunc func(context.Context, *Order) (*Order, error),
	) error
}

type NotFoundError struct {
	OrderID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.OrderID)
}
