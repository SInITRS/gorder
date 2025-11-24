package command

import (
	"context"

	"github.com/SInITRS/gorder/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
