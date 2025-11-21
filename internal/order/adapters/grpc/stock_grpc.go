package grpc

import (
	"context"

	"github.com/SInITRS/gorder/common/genproto/orderpb"
	"github.com/SInITRS/gorder/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
)

// StockGRPC is a gRPC adapter for the StockService.
type StockGRPC struct {
	client stockpb.StockServiceClient
}

// NewStockGRPC creates a new StockGRPC adapter.
func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
	return &StockGRPC{client: client}
}

// CheckIfItemsInStock implements the StockServiceServer interface.
func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
	logrus.Info("stock_grpc response", resp)
	return resp, err
}

// GetItems implements the StockServiceServer interface.
func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
	if err != nil {
		return nil, err
	} else {
		return resp.Items, nil
	}
}
