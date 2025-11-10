package consul

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

type Registry struct {
	client *api.Client
}

var (
	consulClient *Registry
	once         sync.Once
	initErr      error
)

// AI好恐怖啊，以下内容我就写了个func New，剩下的全是Tab补出来的
// New creates a new Consul client instance.
func New(consulAddr string) (*Registry, error) {
	once.Do(func() {
		client, err := api.NewClient(&api.Config{
			Address: consulAddr,
		})
		if err != nil {
			initErr = err
			return
		}
		consulClient = &Registry{client: client}
	})
	if initErr != nil {
		return nil, initErr
	}
	return consulClient, nil
}

func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("Invalid host:port format")
	}
	host := parts[0]
	port, _ := strconv.Atoi(parts[1])

	return r.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: host,
		Port:    port,
		Check: &api.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  false,
			TTL:                            "5s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	logrus.WithFields(logrus.Fields{
		"instance ID":  instanceID,
		"service Name": serviceName,
	}).Info("deregister from consul")
	return r.client.Agent().ServiceDeregister(instanceID)
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)

	if err != nil {
		return nil, err
	}

	var results []string
	for _, entry := range entries {
		address := entry.Service.Address
		port := entry.Service.Port
		results = append(results, address+":"+strconv.Itoa(port))
	}
	return results, nil
}

func (r *Registry) HealthCheck(instanceID string, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
}
