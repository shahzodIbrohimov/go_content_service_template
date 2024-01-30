package storage

import (
	"github.com/golanguzb70/go_content_service/pkg/db"
	"github.com/golanguzb70/go_content_service/pkg/logger"
	"github.com/golanguzb70/go_content_service/storage/postgres"
	"github.com/golanguzb70/go_content_service/storage/repo"
)

type StorageI interface {
	Position() repo.PositionI
	Staff() repo.StaffI
}

type storagePg struct {
	position repo.PositionI
	staff    repo.StaffI
}

func New(db *db.Postgres, log logger.Logger) StorageI {
	return &storagePg{
		position: postgres.NewPositionRepo(db, log),
		staff:    postgres.NewStaffRepo(db, log),
	}
}

func (s *storagePg) Position() repo.PositionI {
	return s.position
}

func (s *storagePg) Staff() repo.StaffI {
	return s.staff
}
