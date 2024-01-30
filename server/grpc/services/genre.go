package services

import (
	"context"

	pb "github.com/golanguzb70/go_content_service/genproto/content_service"
	l "github.com/golanguzb70/go_content_service/pkg/logger"
	grpclient "github.com/golanguzb70/go_content_service/server/grpc/client"
	"github.com/golanguzb70/go_content_service/storage"
)

type GenreService struct {
	storage  storage.StorageI
	logger   l.Logger
	services grpclient.ServiceManager
	pb.UnimplementedGenreServiceServer
}

// New Category Service
func NewGenreService(stroge storage.StorageI, log l.Logger, services grpclient.ServiceManager) *GenreService {
	return &GenreService{
		storage:  stroge,
		logger:   log,
		services: services,
	}
}

func (s *GenreService) Create(ctx context.Context, req *pb.Genre) (*pb.Genre, error) {
	return s.storage.Genre().Create(ctx, req)
}

func (s *GenreService) Get(ctx context.Context, req *pb.Id) (*pb.Genre, error) {
	return s.storage.Genre().Get(ctx, req)
}

func (s *GenreService) Find(ctx context.Context, req *pb.GetListFilter) (*pb.Genres, error) {
	return s.storage.Genre().Find(ctx, req)
}

func (s *GenreService) Update(ctx context.Context, req *pb.Genre) (*pb.Genre, error) {
	return s.storage.Genre().Update(ctx, req)
}

func (s *GenreService) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	return s.storage.Genre().Delete(ctx, req)
}
