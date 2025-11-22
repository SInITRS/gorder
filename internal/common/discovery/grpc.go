package discovery

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/SInITRS/gorder/common/discovery/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// RegisterToConsul registers the service to Consul and starts a heartbeat for health checking.
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
	// Log registration success
	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"addr":        hostPort,
	}).Info("Registered to Consul")
	// Return deregistration function
	return func() error {
		return registry.Deregister(ctx, instanceID, serviceName)
	}, nil
}

func GetServiceAddr(ctx context.Context, serviceName string) (string, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return "", err
	}
	addrs, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return "", err
	}
	if len(addrs) == 0 {
		return "", errors.New("no address found")
	}
	i := rand.Intn(len(addrs))
	logrus.Infof("Discovered %d instance of %s, using %s", len(addrs), serviceName, addrs[i])
	return addrs[i], nil
}
