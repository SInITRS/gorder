package ports

import (
	context "context"

	"github.com/SInITRS/gorder/common/genproto/stockpb"
	"github.com/SInITRS/gorder/stock/app"
)

type GRPCService struct {
	app app.Application
}

func NewGRPCService(app app.Application) *GRPCService {
	return &GRPCService{app: app}
}

func (G GRPCService) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCService) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	//TODO implement me
	panic("implement me")
}
