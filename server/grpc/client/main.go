package grpclient

import (
	"github.com/golanguzb70/go_content_service/config"
)

type ServiceManager interface {
}

type grpcClients struct {
}

func NewGrpcClients(conf *config.Config) (ServiceManager, error) {
	// connect to external clients here
	return &grpcClients{}, nil
}
