package main

import (
	"log"

	"github.com/SInITRS/gorder/common/broker"
	"github.com/SInITRS/gorder/common/config"
	"github.com/SInITRS/gorder/common/logging"
	"github.com/SInITRS/gorder/common/server"
	"github.com/SInITRS/gorder/payment/infrastructure/consumer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatalf("viper.ReadInConfig() failed: %v", err)
	}
	logging.Init()

}

func main() {
	serverType := viper.GetString("payment.server-to-run")

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

	go consumer.NewConsumer().Listen(ch)

	paymentHandler := NewPaymentHandler()
	switch serverType {
	case "http":
		server.RunHTTPServer(viper.GetString("payment.service-name"), paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported server type: grpc")
	default:
		logrus.Panic("unreachable code")
	}
}
