package stock

import (
	"context"
	"fmt"
	"strings"

	"github.com/SInITRS/gorder/stock/entity"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*entity.Item, error)
}

type NotFoundError struct {
	Items []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("These Items not found in stock: %s", strings.Join(e.Items, ","))
}
