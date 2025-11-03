package ports

import (
	"context"

	"github.com/SInITRS/gorder/common/genproto/stockpb"
	"github.com/SInITRS/gorder/stock/app"
	"github.com/SInITRS/gorder/stock/app/query"
)

type GRPCService struct {
	app app.Application
}

func NewGRPCService(app app.Application) *GRPCService {
	return &GRPCService{app: app}
}

func (G GRPCService) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	items, err := G.app.Queries.GetItems.Handle(ctx, query.GetItems{ItemsIDs: request.ItemIDs})
	if err != nil {
		return nil, err
	}
	return &stockpb.GetItemsResponse{Items: items}, nil
}

func (G GRPCService) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: request.Items})
	if err != nil {
		return nil, err
	}
	return &stockpb.CheckIfItemsInStockResponse{
		InStock: 1,
		Items:   items,
	}, nil
}
