package discovery

import (
	"context"
	"time"

	"github.com/SInITRS/gorder/common/discovery/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return func() error { return nil }, err
	}
	instanceID := GenerateInstanceID(serviceName)
	hostPort := viper.Sub(serviceName).GetString("grpc-addr")
	err = registry.Register(ctx, instanceID, serviceName, hostPort)
	if err != nil {
		return func() error { return nil }, err
	}
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logrus.Panicf("no heartbeat from %s to registry, err = %v", serviceName, err)
			}
			time.Sleep(1 * time.Second)

		}

	}()
	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"addr":        hostPort,
	}).Info("Registered to Consul")
	return func() error {
		return registry.Deregister(ctx, instanceID, serviceName)
	}, nil
}
