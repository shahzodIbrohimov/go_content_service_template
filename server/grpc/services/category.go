package services

import (
	"context"

	pb "github.com/golanguzb70/go_content_service/genproto/content_service"
	l "github.com/golanguzb70/go_content_service/pkg/logger"
	grpclient "github.com/golanguzb70/go_content_service/server/grpc/client"
	"github.com/golanguzb70/go_content_service/storage"
)

type CategoryService struct {
	storage  storage.StorageI
	logger   l.Logger
	services grpclient.ServiceManager
	pb.UnimplementedCategoryServiceServer
}

// New Category Service
func NewCategoryService(stroge storage.StorageI, log l.Logger, services grpclient.ServiceManager) *CategoryService {
	return &CategoryService{
		storage:  stroge,
		logger:   log,
		services: services,
	}
}

func (s *CategoryService) Create(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	return s.storage.Category().Create(ctx, req)
}

func (s *CategoryService) Get(ctx context.Context, req *pb.Id) (*pb.Category, error) {
	return s.storage.Category().Get(ctx, req)
}

func (s *CategoryService) Find(ctx context.Context, req *pb.GetListFilter) (*pb.Categories, error) {
	return s.storage.Category().Find(ctx, req)
}

func (s *CategoryService) Update(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	return s.storage.Category().Update(ctx, req)
}

func (s *CategoryService) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	return s.storage.Category().Delete(ctx, req)
}
