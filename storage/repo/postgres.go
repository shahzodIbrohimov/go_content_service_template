package repo

import (
	"context"

	pb "github.com/golanguzb70/go_content_service/genproto/content_service"
)

type PositionI interface {
	Create(ctx context.Context, req *pb.Position) (*pb.Position, error)
	Get(ctx context.Context, req *pb.Id) (*pb.Position, error)
	Find(ctx context.Context, req *pb.GetListFilter) (*pb.Positions, error)
	Update(ctx context.Context, req *pb.Position) (*pb.Position, error)
	Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error)
}
