package services

import (
	"context"

	pb "github.com/golanguzb70/go_content_service/genproto/content_service"
	l "github.com/golanguzb70/go_content_service/pkg/logger"
	grpclient "github.com/golanguzb70/go_content_service/server/grpc/client"
	"github.com/golanguzb70/go_content_service/storage"
)

type CountryService struct {
	storage  storage.StorageI
	logger   l.Logger
	services grpclient.ServiceManager
	pb.UnimplementedCountryServiceServer
}

// New Country Service
func NewCountryService(stroge storage.StorageI, log l.Logger, services grpclient.ServiceManager) *CountryService {
	return &CountryService{
		storage:  stroge,
		logger:   log,
		services: services,
	}
}

func (s *CountryService) Create(ctx context.Context, req *pb.Country) (*pb.Country, error) {
	return s.storage.Country().Create(ctx, req)
}

func (s *CountryService) Get(ctx context.Context, req *pb.Id) (*pb.Country, error) {
	return s.storage.Country().Get(ctx, req)
}

func (s *CountryService) Find(ctx context.Context, req *pb.GetListFilter) (*pb.Countries, error) {
	return s.storage.Country().Find(ctx, req)
}

func (s *CountryService) Update(ctx context.Context, req *pb.Country) (*pb.Country, error) {
	return s.storage.Country().Update(ctx, req)
}

func (s *CountryService) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	return s.storage.Country().Delete(ctx, req)
}
