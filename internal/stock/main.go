package main

import (
	"github.com/SInITRS/gorder/common/genproto/stockpb"
	"github.com/SInITRS/gorder/common/server"
	"github.com/SInITRS/gorder/stock/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")
	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCService()
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
		//not used
	default:
		panic("invalid stock server type")
	}

}
