package adapters

import (
	"context"
	"sync"

	"github.com/SInITRS/gorder/common/genproto/orderpb"
	domain "github.com/SInITRS/gorder/stock/domain/stock"
)

type StockRepository struct {
	lock  *sync.RWMutex
	store map[string]*orderpb.Item
}

var stub = map[string]*orderpb.Item{
	"item_id": {
		ID:       "foo_item",
		Name:     "stub item",
		Quantity: 1000,
		PriceID:  "stub_item_price_id",
	},
}

func NewMemoryStockRepository() *StockRepository {
	return &StockRepository{
		lock:  &sync.RWMutex{},
		store: stub,
	}
}

func (s StockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var (
		res     []*orderpb.Item
		missing []string
	)
	for _, id := range ids {
		if item, exist := s.store[id]; exist {
			res = append(res, item)
		} else {
			missing = append(missing, id)
		}
	}
	if len(res) == len(ids) {
		return res, nil
	}
	return res, domain.NotFoundError{Items: missing}
}
