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
	Tag() repo.TagI
	Genre() repo.GenreI
	Category() repo.CategoryI
	Country() repo.CountryI
}

type storagePg struct {
	position repo.PositionI
	staff    repo.StaffI
	tag      repo.TagI
	genre    repo.GenreI
	category repo.CategoryI
	country  repo.CountryI
}

func New(db *db.Postgres, log logger.Logger) StorageI {
	return &storagePg{
		position: postgres.NewPositionRepo(db, log),
		staff:    postgres.NewStaffRepo(db, log),
		tag:      postgres.NewTagRepo(db, log),
		genre:    postgres.NewGenreRepo(db, log),
		category: postgres.NewCategoryRepo(db, log),
		country:  postgres.NewCountryRepo(db, log),
	}
}

func (s *storagePg) Position() repo.PositionI {
	return s.position
}

func (s *storagePg) Staff() repo.StaffI {
	return s.staff
}

func (s *storagePg) Tag() repo.TagI {
	return s.tag
}

func (s *storagePg) Genre() repo.GenreI {
	return s.genre
}

func (s *storagePg) Category() repo.CategoryI {
	return s.category
}

func (s *storagePg) Country() repo.CountryI {
	return s.country
}
