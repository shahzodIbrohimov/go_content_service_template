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

type tag struct {
	db  *db.Postgres
	log logger.Logger
}

func NewTagRepo(db *db.Postgres, log logger.Logger) repo.TagI {
	return &tag{
		db:  db,
		log: log,
	}
}

func (r *tag) Create(ctx context.Context, req *pb.Tag) (*pb.Tag, error) {
	query := r.db.Builder.Insert("tags").
		Columns(`
			id, title_uz, title_ru, title_en, active, color
		`).
		Values(
			req.Id, req.Title.Uz, req.Title.Ru, req.Title.En, req.Active, req.Color,
		).
		Suffix("RETURNING created_at, updated_at")

	err := query.RunWith(r.db.Db).QueryRow().Scan(&req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return req, HandleDatabaseError(err, r.log, "create tag")
	}

	return req, nil
}

func (r *tag) Get(ctx context.Context, req *pb.Id) (*pb.Tag, error) {
	res := &pb.Tag{
		Title: &pb.MultiLanguage{},
	}

	query := r.db.Builder.Select(`
		id, title_uz, title_ru, title_en, active, color, created_at, updated_at
	`).From("tags")

	if req.Id != "" {
		query = query.Where(squirrel.Eq{"id": req.Id})
	} else {
		return nil, fmt.Errorf("id is required")
	}

	err := query.RunWith(r.db.Db).QueryRow().Scan(
		&res.Id, &res.Title.Uz, &res.Title.Ru, &res.Title.En, &res.Active,
		&res.Color, &res.CreatedAt, &res.UpdatedAt,
	)

	return res, HandleDatabaseError(err, r.log, "getting tag")
}

func (r *tag) Find(ctx context.Context, req *pb.GetListFilter) (*pb.Tags, error) {
	var (
		res            = &pb.Tags{}
		whereCondition = PrepareWhere(req.Filters)
		orderBy        = PrepareOrder(req.Sorts)
	)

	query := r.db.Builder.Select(`
		id, title_uz, title_ru, title_en, active, color, created_at, updated_at
	`).From("tags")

	if len(req.Filters) > 0 {
		query = query.Where(whereCondition)

	}
	if len(req.Sorts) > 0 {
		query = query.OrderBy(orderBy)
	}

	query = query.Offset(uint64((req.Page - 1) * req.Limit)).Limit(uint64(req.Limit))

	rows, err := query.RunWith(r.db.Db).Query()
	if err != nil {
		return nil, HandleDatabaseError(err, r.log, "error while finding tag")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &pb.Tag{
			Title: &pb.MultiLanguage{},
		}
		err = rows.Scan(
			&temp.Id, &temp.Title.Uz, &temp.Title.Ru, &temp.Title.En, &temp.Active,
			&temp.Color, &temp.CreatedAt, &temp.UpdatedAt,
		)
		if err != nil {
			return nil, HandleDatabaseError(err, r.log, "error while scanning tag")
		}

		res.Items = append(res.Items, temp)
	}

	count := r.db.Builder.Select("count(1)").
		From("tags").
		Where(whereCondition)

	err = count.RunWith(r.db.Db).Scan(&res.Count)

	return res, HandleDatabaseError(err, r.log, "getting tag count")
}

func (r *tag) Update(ctx context.Context, req *pb.Tag) (*pb.Tag, error) {
	var (
		mp             = make(map[string]interface{})
		whereCondition = squirrel.And{squirrel.Eq{"id": req.Id}}
	)
	mp["title_uz"] = req.Title.Uz
	mp["title_ru"] = req.Title.Ru
	mp["title_en"] = req.Title.En
	mp["active"] = req.Active
	mp["color"] = req.Color
	mp["updated_at"] = time.Now()

	query := r.db.Builder.Update("tags").SetMap(mp).
		Where(whereCondition).
		Suffix("RETURNING updated_at, created_at")

	err := query.RunWith(r.db.Db).QueryRow().Scan(
		&req.CreatedAt, &req.UpdatedAt,
	)

	return req, HandleDatabaseError(err, r.log, "update tag")
}

func (r *tag) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	query := r.db.Builder.Delete("tags").Where(squirrel.Eq{"id": req.Id})
	_, err := query.RunWith(r.db.Db).Exec()
	return &pb.Empty{}, HandleDatabaseError(err, r.log, "delete from tag")
}
