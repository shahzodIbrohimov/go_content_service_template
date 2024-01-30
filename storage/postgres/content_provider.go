package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	pb "github.com/golanguzb70/go_content_service/genproto/content_service"
	"github.com/golanguzb70/go_content_service/pkg/db"
	"github.com/golanguzb70/go_content_service/pkg/logger"
	"github.com/golanguzb70/go_content_service/storage/repo"
)

type contentProvider struct {
	db  *db.Postgres
	log logger.Logger
}

func NewContentProviderRepo(db *db.Postgres, log logger.Logger) repo.ContentProviderI {
	return &contentProvider{
		db:  db,
		log: log,
	}
}

func (r *contentProvider) Create(ctx context.Context, req *pb.ContentProvider) (*pb.ContentProvider, error) {
	query := r.db.Builder.Insert("content_providers").
		Columns(`
			id, slug, title_uz, title_ru, title_en, active, logo
		`).
		Values(
			req.Id, req.Slug, req.Title.Uz, req.Title.Ru,
			req.Title.En, req.Active, req.Logo,
		).
		Suffix("RETURNING created_at, updated_at")

	err := query.RunWith(r.db.Db).QueryRow().Scan(&req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return req, HandleDatabaseError(err, r.log, "create content_provider")
	}

	return req, nil
}

func (r *contentProvider) Get(ctx context.Context, req *pb.Id) (*pb.ContentProvider, error) {
	res := &pb.ContentProvider{
		Title: &pb.MultiLanguage{},
	}

	query := r.db.Builder.Select(`
		id, slug, title_uz, title_ru, title_en, active, logo, sort_order, created_at, updated_at
	`).From("content_providers")

	if req.Id != "" {
		query = query.Where(squirrel.Eq{"id": req.Id})
	} else {
		return nil, fmt.Errorf("id is required")
	}

	err := query.RunWith(r.db.Db).QueryRow().Scan(
		&res.Id, &res.Slug, &res.Title.Uz, &res.Title.Ru,
		&res.Title.En, &res.Active, &res.Logo, &res.Order, &res.CreatedAt, &res.UpdatedAt,
	)

	return res, HandleDatabaseError(err, r.log, "getting content_provider")
}

func (r *contentProvider) Find(ctx context.Context, req *pb.GetListFilter) (*pb.ContentProviders, error) {
	var (
		res            = &pb.ContentProviders{}
		whereCondition = PrepareWhere(req.Filters)
		orderBy        = PrepareOrder(req.Sorts)
	)

	query := r.db.Builder.Select(`
		id, slug, title_uz, title_ru, title_en, active, logo, sort_order, created_at, updated_at
	`).From("content_providers")

	if len(req.Filters) > 0 {
		query = query.Where(whereCondition)

	}
	if len(req.Sorts) > 0 {
		query = query.OrderBy(orderBy)
	}

	query = query.Offset(uint64((req.Page - 1) * req.Limit)).Limit(uint64(req.Limit))

	rows, err := query.RunWith(r.db.Db).Query()
	if err != nil {
		return nil, HandleDatabaseError(err, r.log, "error while finding content_provider")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &pb.ContentProvider{
			Title: &pb.MultiLanguage{},
		}
		err = rows.Scan(
			&temp.Id, &temp.Slug, &temp.Title.Uz, &temp.Title.Ru,
			&temp.Title.En, &temp.Active, &temp.Logo, &temp.Order, &temp.CreatedAt, &temp.UpdatedAt,
		)
		if err != nil {
			return nil, HandleDatabaseError(err, r.log, "error while scanning content_provider")
		}

		res.Items = append(res.Items, temp)
	}

	count := r.db.Builder.Select("count(1)").
		From("content_providers").
		Where(whereCondition)

	err = count.RunWith(r.db.Db).Scan(&res.Count)

	return res, HandleDatabaseError(err, r.log, "getting content_provider count")
}

func (r *contentProvider) Update(ctx context.Context, req *pb.ContentProvider) (*pb.ContentProvider, error) {
	var (
		mp             = make(map[string]interface{})
		whereCondition = squirrel.And{squirrel.Eq{"id": req.Id}}
	)

	mp["title_uz"] = req.Title.Uz
	mp["title_ru"] = req.Title.Ru
	mp["title_en"] = req.Title.En
	mp["active"] = req.Active
	mp["slug"] = req.Slug
	mp["logo"] = req.Logo
	mp["updated_at"] = time.Now()

	query := r.db.Builder.Update("content_providers").SetMap(mp).
		Where(whereCondition).
		Suffix("RETURNING updated_at, created_at")

	err := query.RunWith(r.db.Db).QueryRow().Scan(
		&req.CreatedAt, &req.UpdatedAt,
	)

	return req, HandleDatabaseError(err, r.log, "update content_provider")
}

func (r *contentProvider) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	query := r.db.Builder.Delete("content_providers").Where(squirrel.Eq{"id": req.Id})
	_, err := query.RunWith(r.db.Db).Exec()
	return &pb.Empty{}, HandleDatabaseError(err, r.log, "delete from content_provider")
}
