package grpc

import (
	"fmt"

	"net"

	"github.com/golanguzb70/go_content_service/config"
	pb "github.com/golanguzb70/go_content_service/genproto/content_service"
	"github.com/golanguzb70/go_content_service/pkg/db"
	l "github.com/golanguzb70/go_content_service/pkg/logger"
	grpclient "github.com/golanguzb70/go_content_service/server/grpc/client"
	"github.com/golanguzb70/go_content_service/server/grpc/services"
	"github.com/golanguzb70/go_content_service/storage"

	_ "github.com/lib/pq"

	"google.golang.org/grpc"
)

type GRPCService struct {
	PositionService *services.PositionService
}

func New(cfg *config.Config, log l.Logger) (*GRPCService, error) {
	psql, err := db.New(*cfg)
	if err != nil {
		return nil, fmt.Errorf("Error while connecting to database: %v", err)
	}

	storageObj := storage.New(psql, log)

	grpcClient, err := grpclient.NewGrpcClients(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error while connecting with grpc clients: %v", err)
	}

	return &GRPCService{
		PositionService: services.NewPositionService(storageObj, log, grpcClient),
	}, nil
}

func (service *GRPCService) Run(logger l.Logger, cfg *config.Config) {
	server := grpc.NewServer()

	pb.RegisterPositionServiceServer(server, service.PositionService)

	listener, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		panic("Error while creating listener")
	}

	logger.Info("GRPC server is running on port " + cfg.RPCPort)

	err = server.Serve(listener)
	if err != nil {
		panic("error while serving gRPC server on port " + cfg.RPCPort)
	}
}
