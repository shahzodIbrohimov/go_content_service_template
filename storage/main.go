package storage

import (
	"github.com/golanguzb70/go_content_service/pkg/db"
	"github.com/golanguzb70/go_content_service/pkg/logger"
	"github.com/golanguzb70/go_content_service/storage/postgres"
	"github.com/golanguzb70/go_content_service/storage/repo"
)

type StorageI interface {
	Position() repo.PositionI
}

type storagePg struct {
	position repo.PositionI
}

func New(db *db.Postgres, log logger.Logger) StorageI {
	return &storagePg{
		position: postgres.NewPositionRepo(db, log),
	}
}

func (s *storagePg) Position() repo.PositionI {
	return s.position
}
