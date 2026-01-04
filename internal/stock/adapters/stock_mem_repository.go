package adapters

import (
	"context"
	"sync"

	domain "github.com/SInITRS/gorder/stock/domain/stock"
	"github.com/SInITRS/gorder/stock/entity"
)

type StockRepository struct {
	lock  *sync.RWMutex
	store map[string]*entity.Item
}

var stub = map[string]*entity.Item{
	"item_id": {
		ID:       "foo_item",
		Name:     "stub item",
		Quantity: 1000,
		PriceID:  "stub_item_price_id",
	},
	"item1": {
		ID:       "foo_item1",
		Name:     "stub item1",
		Quantity: 1000,
		PriceID:  "stub_item_price_id1",
	},
	"item2": {
		ID:       "foo_item2",
		Name:     "stub item2",
		Quantity: 1000,
		PriceID:  "stub_item_price_id2",
	},
	"item3": {
		ID:       "foo_item3",
		Name:     "stub item3",
		Quantity: 1000,
		PriceID:  "stub_item_price_id3",
	},
}

func NewMemoryStockRepository() *StockRepository {
	return &StockRepository{
		lock:  &sync.RWMutex{},
		store: stub,
	}
}

func (s StockRepository) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var (
		res     []*entity.Item
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
