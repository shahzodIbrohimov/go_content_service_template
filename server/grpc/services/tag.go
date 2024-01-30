package services

import (
	"context"

	pb "github.com/golanguzb70/go_content_service/genproto/content_service"
	l "github.com/golanguzb70/go_content_service/pkg/logger"
	grpclient "github.com/golanguzb70/go_content_service/server/grpc/client"
	"github.com/golanguzb70/go_content_service/storage"
)

type TagService struct {
	storage  storage.StorageI
	logger   l.Logger
	services grpclient.ServiceManager
	pb.UnimplementedTagServiceServer
}

// New Category Service
func NewTagService(stroge storage.StorageI, log l.Logger, services grpclient.ServiceManager) *TagService {
	return &TagService{
		storage:  stroge,
		logger:   log,
		services: services,
	}
}

func (s *TagService) Create(ctx context.Context, req *pb.Tag) (*pb.Tag, error) {
	return s.storage.Tag().Create(ctx, req)
}

func (s *TagService) Get(ctx context.Context, req *pb.Id) (*pb.Tag, error) {
	return s.storage.Tag().Get(ctx, req)
}

func (s *TagService) Find(ctx context.Context, req *pb.GetListFilter) (*pb.Tags, error) {
	return s.storage.Tag().Find(ctx, req)
}

func (s *TagService) Update(ctx context.Context, req *pb.Tag) (*pb.Tag, error) {
	return s.storage.Tag().Update(ctx, req)
}

func (s *TagService) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	return s.storage.Tag().Delete(ctx, req)
}
