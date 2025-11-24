package domain

import (
	"context"

	"github.com/SInITRS/gorder/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
}
