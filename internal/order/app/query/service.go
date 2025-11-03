package query

import (
	"context"

	"github.com/SInITRS/gorder/common/genproto/orderpb"
	"github.com/SInITRS/gorder/common/genproto/stockpb"
)

type StockService interface {
	CheckIfItemsInStock(context.Context, []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}
