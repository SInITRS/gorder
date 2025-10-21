package service

import (
	"context"

	grpcClient "github.com/SInITRS/gorder/common/client"
	"github.com/SInITRS/gorder/common/metrics"
	"github.com/SInITRS/gorder/order/adapters"
	"github.com/SInITRS/gorder/order/adapters/grpc"
	"github.com/SInITRS/gorder/order/app"
	"github.com/SInITRS/gorder/order/app/command"
	"github.com/SInITRS/gorder/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {

	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	stockGRPC := grpc.NewStockGRPC(stockClient)
	if err != nil {
		panic(err)
	}
	return newApplication(ctx, stockGRPC), func() {
		_ = closeStockClient()
	}

}

func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricsClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricsClient),
		},
	}
}
