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

type category struct {
	db  *db.Postgres
	log logger.Logger
}

func NewCategoryRepo(db *db.Postgres, log logger.Logger) repo.CategoryI {
	return &category{
		db:  db,
		log: log,
	}
}

func (r *category) Create(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	query := r.db.Builder.Insert("categories").
		Columns(`
			id, slug, title_uz, title_ru, title_en, active, photo
		`).
		Values(
			req.Id, req.Slug, req.Title.Uz, req.Title.Ru,
			req.Title.En, req.Active, req.Photo,
		).
		Suffix("RETURNING created_at, updated_at")

	err := query.RunWith(r.db.Db).QueryRow().Scan(&req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return req, HandleDatabaseError(err, r.log, "create category")
	}

	return req, nil
}

func (r *category) Get(ctx context.Context, req *pb.Id) (*pb.Category, error) {
	res := &pb.Category{
		Title: &pb.MultiLanguage{},
	}

	query := r.db.Builder.Select(`
		id, slug, title_uz, title_ru, title_en, active, photo, sort_order, created_at, updated_at
	`).From("categories")

	if req.Id != "" {
		query = query.Where(squirrel.Eq{"id": req.Id})
	} else {
		return nil, fmt.Errorf("id is required")
	}

	err := query.RunWith(r.db.Db).QueryRow().Scan(
		&res.Id, &res.Slug, &res.Title.Uz, &res.Title.Ru,
		&res.Title.En, &res.Active, &res.Photo, &res.Order, &res.CreatedAt, &res.UpdatedAt,
	)

	return res, HandleDatabaseError(err, r.log, "getting category")
}

func (r *category) Find(ctx context.Context, req *pb.GetListFilter) (*pb.Categories, error) {
	var (
		res            = &pb.Categories{}
		whereCondition = PrepareWhere(req.Filters)
		orderBy        = PrepareOrder(req.Sorts)
	)

	query := r.db.Builder.Select(`
		id, slug, title_uz, title_ru, title_en, active, photo, sort_order, created_at, updated_at
	`).From("categories")

	if len(req.Filters) > 0 {
		query = query.Where(whereCondition)

	}
	if len(req.Sorts) > 0 {
		query = query.OrderBy(orderBy)
	}

	query = query.Offset(uint64((req.Page - 1) * req.Limit)).Limit(uint64(req.Limit))

	rows, err := query.RunWith(r.db.Db).Query()
	if err != nil {
		return nil, HandleDatabaseError(err, r.log, "error while finding category")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &pb.Category{
			Title: &pb.MultiLanguage{},
		}
		err = rows.Scan(
			&temp.Id, &temp.Slug, &temp.Title.Uz, &temp.Title.Ru,
			&temp.Title.En, &temp.Active, &temp.Photo, &temp.Order, &temp.CreatedAt, &temp.UpdatedAt,
		)
		if err != nil {
			return nil, HandleDatabaseError(err, r.log, "error while scanning category")
		}

		res.Items = append(res.Items, temp)
	}

	count := r.db.Builder.Select("count(1)").
		From("categories").
		Where(whereCondition)

	err = count.RunWith(r.db.Db).Scan(&res.Count)

	return res, HandleDatabaseError(err, r.log, "getting category count")
}

func (r *category) Update(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	var (
		mp             = make(map[string]interface{})
		whereCondition = squirrel.And{squirrel.Eq{"id": req.Id}}
	)

	mp["title_uz"] = req.Title.Uz
	mp["title_ru"] = req.Title.Ru
	mp["title_en"] = req.Title.En
	mp["active"] = req.Active
	mp["slug"] = req.Slug
	mp["photo"] = req.Photo
	mp["updated_at"] = time.Now()

	query := r.db.Builder.Update("categories").SetMap(mp).
		Where(whereCondition).
		Suffix("RETURNING updated_at, created_at")

	err := query.RunWith(r.db.Db).QueryRow().Scan(
		&req.CreatedAt, &req.UpdatedAt,
	)

	return req, HandleDatabaseError(err, r.log, "update category")
}

func (r *category) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	query := r.db.Builder.Delete("categories").Where(squirrel.Eq{"id": req.Id})
	_, err := query.RunWith(r.db.Db).Exec()
	return &pb.Empty{}, HandleDatabaseError(err, r.log, "delete from category")
}
