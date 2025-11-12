package main

import (
	"context"
	"log"

	"github.com/SInITRS/gorder/common/config"
	"github.com/SInITRS/gorder/common/discovery"
	"github.com/SInITRS/gorder/common/genproto/stockpb"
	"github.com/SInITRS/gorder/common/logging"
	"github.com/SInITRS/gorder/common/server"
	"github.com/SInITRS/gorder/stock/ports"
	"github.com/SInITRS/gorder/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatalf("viper.ReadInConfig() failed: %v", err)
	}
	logging.Init()

}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := service.NewApplication(ctx)

	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCService(application)
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
		//not used
	default:
		panic("invalid stock server type")
	}

}
