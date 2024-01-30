package services

import (
	"context"

	pb "github.com/golanguzb70/go_content_service/genproto/content_service"
	l "github.com/golanguzb70/go_content_service/pkg/logger"
	grpclient "github.com/golanguzb70/go_content_service/server/grpc/client"
	"github.com/golanguzb70/go_content_service/storage"
)

type PositionService struct {
	storage  storage.StorageI
	logger   l.Logger
	services grpclient.ServiceManager
	pb.UnimplementedPositionServiceServer
}

// New Category Service
func NewPositionService(stroge storage.StorageI, log l.Logger, services grpclient.ServiceManager) *PositionService {
	return &PositionService{
		storage:  stroge,
		logger:   log,
		services: services,
	}
}

func (s *PositionService) Create(ctx context.Context, req *pb.Position) (*pb.Position, error) {
	return s.storage.Position().Create(ctx, req)
}

func (s *PositionService) Get(ctx context.Context, req *pb.Id) (*pb.Position, error) {
	return s.storage.Position().Get(ctx, req)
}

func (s *PositionService) Find(ctx context.Context, req *pb.GetListFilter) (*pb.Positions, error) {
	return s.storage.Position().Find(ctx, req)
}

func (s *PositionService) Update(ctx context.Context, req *pb.Position) (*pb.Position, error) {
	return s.storage.Position().Update(ctx, req)
}

func (s *PositionService) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	return s.storage.Position().Delete(ctx, req)
}
