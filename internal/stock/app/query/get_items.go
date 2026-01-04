package query

import (
	"context"

	"github.com/SInITRS/gorder/common/decorator"
	domain "github.com/SInITRS/gorder/stock/domain/stock"
	"github.com/SInITRS/gorder/stock/entity"
	"github.com/sirupsen/logrus"
)

type GetItems struct {
	ItemsIDs []string
}

type GetItemsHandler decorator.QueryHandler[GetItems, []*entity.Item]

type getItemsHandler struct {
	stockRepo domain.Repository
}

func (g getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*entity.Item, error) {
	items, err := g.stockRepo.GetItems(ctx, query.ItemsIDs)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func NewGetItemsHandler(
	stockRepo domain.Repository,
	logger *logrus.Entry,
	client decorator.MetricsClient,
) GetItemsHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	return decorator.ApplyQueryDecorators(
		getItemsHandler{stockRepo: stockRepo},
		logger,
		client,
	)
}
