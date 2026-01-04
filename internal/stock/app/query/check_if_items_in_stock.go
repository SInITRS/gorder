package query

import (
	"context"

	"github.com/SInITRS/gorder/common/decorator"
	domain "github.com/SInITRS/gorder/stock/domain/stock"
	"github.com/SInITRS/gorder/stock/entity"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
}

// Deprecated: this is a stub implementation.
var stub = map[string]string{
	"1": "price_1SSH8TBK7fImhh4SQD4QSQJJ",
	"2": "price_1SXNBBBK7fImhh4SHTMTlY11",
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
	var res []*entity.Item
	for _, item := range query.Items {
		priceID, ok := stub[item.ID]
		if !ok {
			priceID = stub["1"]
		}
		res = append(res, &entity.Item{
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
