package query

import (
	"context"

	"github.com/SInITRS/gorder/common/decorator"
	"github.com/SInITRS/gorder/common/genproto/orderpb"
	domain "github.com/SInITRS/gorder/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*orderpb.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
}

var stub = map[string]string{
	"1": "price_1SSH8TBK7fImhh4SQD4QSQJJ",
	"2": "price_1SXNBBBK7fImhh4SHTMTlY11",
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var res []*orderpb.Item
	for _, item := range query.Items {
		priceID, ok := stub[item.ID]
		if !ok {
			priceID = stub["1"]
		}
		res = append(res, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
			PriceID:  priceID,
		})
	}
	return res, nil
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	return decorator.ApplyQueryDecorators(
		checkIfItemsInStockHandler{stockRepo: stockRepo},
		logger,
		metricClient,
	)
}
