package query

import (
	"context"

	"github.com/SInITRS/gorder/common/genproto/orderpb"
)

type StockService interface {
	CheckIfItemsInStock(context.Context, []*orderpb.ItemWithQuantity) error
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}
