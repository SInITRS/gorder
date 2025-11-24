package service

import (
	"context"

	"github.com/SInITRS/gorder/common/broker"
	grpcClient "github.com/SInITRS/gorder/common/client"
	"github.com/SInITRS/gorder/common/metrics"
	"github.com/SInITRS/gorder/order/adapters"
	"github.com/SInITRS/gorder/order/adapters/grpc"
	"github.com/SInITRS/gorder/order/app"
	"github.com/SInITRS/gorder/order/app/command"
	"github.com/SInITRS/gorder/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)

	channel, closech := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)

	return newApplication(ctx, stockGRPC, channel), func() {
		_ = closeStockClient()
		_ = channel.Close()
		_ = closech()
	}

}

func newApplication(_ context.Context, stockGRPC query.StockService, channel *amqp.Channel) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, channel, logger, metricsClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricsClient),
		},
	}
}
