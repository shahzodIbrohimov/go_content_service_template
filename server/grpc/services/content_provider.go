package services

import (
	"context"

	pb "github.com/golanguzb70/go_content_service/genproto/content_service"
	l "github.com/golanguzb70/go_content_service/pkg/logger"
	grpclient "github.com/golanguzb70/go_content_service/server/grpc/client"
	"github.com/golanguzb70/go_content_service/storage"
)

type ContentProviderService struct {
	storage  storage.StorageI
	logger   l.Logger
	services grpclient.ServiceManager
	pb.UnimplementedContentProviderServiceServer
}

// New ContentProvider Service
func NewContentProviderService(stroge storage.StorageI, log l.Logger, services grpclient.ServiceManager) *ContentProviderService {
	return &ContentProviderService{
		storage:  stroge,
		logger:   log,
		services: services,
	}
}

func (s *ContentProviderService) Create(ctx context.Context, req *pb.ContentProvider) (*pb.ContentProvider, error) {
	return s.storage.ContentProvider().Create(ctx, req)
}

func (s *ContentProviderService) Get(ctx context.Context, req *pb.Id) (*pb.ContentProvider, error) {
	return s.storage.ContentProvider().Get(ctx, req)
}

func (s *ContentProviderService) Find(ctx context.Context, req *pb.GetListFilter) (*pb.ContentProviders, error) {
	return s.storage.ContentProvider().Find(ctx, req)
}

func (s *ContentProviderService) Update(ctx context.Context, req *pb.ContentProvider) (*pb.ContentProvider, error) {
	return s.storage.ContentProvider().Update(ctx, req)
}

func (s *ContentProviderService) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	return s.storage.ContentProvider().Delete(ctx, req)
}
