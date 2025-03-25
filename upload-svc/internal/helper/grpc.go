package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Registry interface {
	RegisterService(ctx context.Context, serviceName, serviceID, serviceAddress string, servicePort int, tags []string) error
	DeregisterService(ctx context.Context, serviceID string) error
	GetService(ctx context.Context, serviceName string) ([]*consul.ServiceEntry, error)
	HealthCheck(serviceID, serviceName string) error
}

func GenerateServiceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

func ServiceConnection(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
	service, err := registry.GetService(ctx, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	if len(service) == 0 {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	serviceEntry := service[rand.Intn(len(service))]

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", serviceEntry.Service.Address, serviceEntry.Service.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return conn, nil
}
