package main

import (
	"log"

	"github.com/SInITRS/gorder/common/config"
	"github.com/SInITRS/gorder/common/logging"
	"github.com/SInITRS/gorder/common/server"
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

	servertype := viper.GetString("payment.server-to-run")

	paymentHandler := NewPaymentHandler()

	switch servertype {
	case "http":
		server.RunHTTPServer(viper.GetString("payment.service-name"), paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("grpc is unsupported")
	default:
		logrus.Panicf("unknown server type: %s", servertype)
	}
}
