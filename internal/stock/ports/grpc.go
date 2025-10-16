package ports

import (
	context "context"

	"github.com/SInITRS/gorder/common/genproto/stockpb"
)

type GRPCService struct {
}

func NewGRPCService() *GRPCService {
	return &GRPCService{}
}

func (G GRPCService) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCService) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	//TODO implement me
	panic("implement me")
}
