package server

import (
	"net"

	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func init() {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	grpcLogrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
}

// RunGRPCServer starts a gRPC server for the given service name.
func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
	addr := viper.Sub(serviceName).GetString("grpc-addr")
	if addr == "" {
		// TODO: Warning log
		addr = viper.GetString("fallback-grpc-addr")
	}

	RunGRPCServerOnAddr(addr, registerServer)
}

// RunGRPCServerOnAddr starts a gRPC server on the specified address.
func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			grpcTags.UnaryServerInterceptor(grpcTags.WithFieldExtractor(grpcTags.CodeGenRequestFieldExtractor)),
			grpcLogrus.UnaryServerInterceptor(logrusEntry),
		),
		grpc.ChainStreamInterceptor(
			grpcTags.StreamServerInterceptor(grpcTags.WithFieldExtractor(grpcTags.CodeGenRequestFieldExtractor)),
			grpcLogrus.StreamServerInterceptor(logrusEntry),
		),
	)
	registerServer(grpcServer)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("Starting gRPC server, Listening: %s", addr)
	if err := grpcServer.Serve(listen); err != nil {
		logrus.Panic(err)
	}
}
