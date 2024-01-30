package services

import (
	"context"

	pb "github.com/golanguzb70/go_content_service/genproto/content_service"
	l "github.com/golanguzb70/go_content_service/pkg/logger"
	grpclient "github.com/golanguzb70/go_content_service/server/grpc/client"
	"github.com/golanguzb70/go_content_service/storage"
)

type StaffService struct {
	storage  storage.StorageI
	logger   l.Logger
	services grpclient.ServiceManager
	pb.UnimplementedStaffServiceServer
}

// New Category Service
func NewStaffService(stroge storage.StorageI, log l.Logger, services grpclient.ServiceManager) *StaffService {
	return &StaffService{
		storage:  stroge,
		logger:   log,
		services: services,
	}
}

func (s *StaffService) Create(ctx context.Context, req *pb.Staff) (*pb.Staff, error) {
	return s.storage.Staff().Create(ctx, req)
}

func (s *StaffService) Get(ctx context.Context, req *pb.Id) (*pb.Staff, error) {
	return s.storage.Staff().Get(ctx, req)
}

func (s *StaffService) Find(ctx context.Context, req *pb.GetListFilter) (*pb.Staffs, error) {
	return s.storage.Staff().Find(ctx, req)
}

func (s *StaffService) Update(ctx context.Context, req *pb.Staff) (*pb.Staff, error) {
	return s.storage.Staff().Update(ctx, req)
}

func (s *StaffService) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	return s.storage.Staff().Delete(ctx, req)
}
