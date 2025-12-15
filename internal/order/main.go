package main

import (
	"context"
	"log"

	"github.com/SInITRS/gorder/common/broker"
	"github.com/SInITRS/gorder/common/config"
	"github.com/SInITRS/gorder/common/discovery"
	"github.com/SInITRS/gorder/common/genproto/orderpb"
	"github.com/SInITRS/gorder/common/logging"
	"github.com/SInITRS/gorder/common/server"
	"github.com/SInITRS/gorder/common/tracing"
	"github.com/SInITRS/gorder/order/infrastructure/consumer"
	"github.com/SInITRS/gorder/order/ports"
	"github.com/SInITRS/gorder/order/service"
	"github.com/gin-gonic/gin"
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
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)
	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()
	go consumer.NewConsumer(application).Listen(ch)

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		router.StaticFile("/success", "../../public/success.html")
		ports.RegisterHandlersWithOptions(
			router,
			HTTPServer{
				app: application,
			},
			ports.GinServerOptions{
				BaseURL:      "/api",
				Middlewares:  nil,
				ErrorHandler: nil,
			})
	})

}
