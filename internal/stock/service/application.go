package service

import (
	"context"

	metrics "github.com/SInITRS/gorder/common/metrics"
	"github.com/SInITRS/gorder/stock/adapters"
	"github.com/SInITRS/gorder/stock/app"
	"github.com/SInITRS/gorder/stock/app/query"
	"github.com/SInITRS/gorder/stock/infrastructure/integration"
	"github.com/sirupsen/logrus"
)

func NewApplication(_ context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	stripeAPI := integration.NewStripeAPI()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
		},
	}
}
