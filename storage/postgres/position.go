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

type position struct {
	db  *db.Postgres
	log logger.Logger
}

func NewPositionRepo(db *db.Postgres, log logger.Logger) repo.PositionI {
	return &position{
		db:  db,
		log: log,
	}
}

func (r *position) Create(ctx context.Context, req *pb.Position) (*pb.Position, error) {
	query := r.db.Builder.Insert("positions").
		Columns(`
			id, title_uz, title_ru, title_en, active
		`).
		Values(
			req.Id, req.Title.Uz, req.Title.Ru, req.Title.En, req.Active,
		).
		Suffix("RETURNING created_at, updated_at")

	err := query.RunWith(r.db.Db).QueryRow().Scan(&req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return req, HandleDatabaseError(err, r.log, "create position")
	}

	return req, nil
}

func (r *position) Get(ctx context.Context, req *pb.Id) (*pb.Position, error) {
	res := &pb.Position{
		Title: &pb.MultiLanguage{},
	}

	query := r.db.Builder.Select(`
		id, title_uz, title_ru, title_en, active, created_at, updated_at
	`).From("positions")

	if req.Id != "" {
		query = query.Where(squirrel.Eq{"id": req.Id})
	} else {
		return nil, fmt.Errorf("id is required")
	}

	err := query.RunWith(r.db.Db).QueryRow().Scan(
		&res.Id, &res.Title.Uz, &res.Title.Ru, &res.Title.En, &res.Active,
		&res.CreatedAt, &res.UpdatedAt,
	)

	return res, HandleDatabaseError(err, r.log, "getting position")
}

func (r *position) Find(ctx context.Context, req *pb.GetListFilter) (*pb.Positions, error) {
	var (
		res            = &pb.Positions{}
		whereCondition = PrepareWhere(req.Filters)
		orderBy        = PrepareOrder(req.Sorts)
	)

	query := r.db.Builder.Select(`
		id, title_uz, title_ru, title_en, active, created_at, updated_at
	`).From("positions")

	if len(req.Filters) > 0 {
		query = query.Where(whereCondition)

	}
	if len(req.Sorts) > 0 {
		query = query.OrderBy(orderBy)
	}

	query = query.Offset(uint64((req.Page - 1) * req.Limit)).Limit(uint64(req.Limit))

	rows, err := query.RunWith(r.db.Db).Query()
	if err != nil {
		return nil, HandleDatabaseError(err, r.log, "error while finding position")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &pb.Position{
			Title: &pb.MultiLanguage{},
		}
		err = rows.Scan(
			&temp.Id, &temp.Title.Uz, &temp.Title.Ru, &temp.Title.En, &temp.Active,
			&temp.CreatedAt, &temp.UpdatedAt,
		)
		if err != nil {
			return nil, HandleDatabaseError(err, r.log, "error while scanning position")
		}

		res.Items = append(res.Items, temp)
	}

	count := r.db.Builder.Select("count(1)").
		From("positions").
		Where(whereCondition)

	err = count.RunWith(r.db.Db).Scan(&res.Count)

	return res, HandleDatabaseError(err, r.log, "getting position count")
}

func (r *position) Update(ctx context.Context, req *pb.Position) (*pb.Position, error) {
	var (
		mp             = make(map[string]interface{})
		whereCondition = squirrel.And{squirrel.Eq{"id": req.Id}}
	)
	mp["title_uz"] = req.Title.Uz
	mp["title_ru"] = req.Title.Ru
	mp["title_en"] = req.Title.En
	mp["active"] = req.Active
	mp["updated_at"] = time.Now()

	query := r.db.Builder.Update("positions").SetMap(mp).
		Where(whereCondition).
		Suffix("RETURNING updated_at, created_at")

	err := query.RunWith(r.db.Db).QueryRow().Scan(
		&req.CreatedAt, &req.UpdatedAt,
	)

	return req, HandleDatabaseError(err, r.log, "update position")
}

func (r *position) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	query := r.db.Builder.Delete("positions").Where(squirrel.Eq{"id": req.Id})
	_, err := query.RunWith(r.db.Db).Exec()
	return &pb.Empty{}, HandleDatabaseError(err, r.log, "delete from position")
}
