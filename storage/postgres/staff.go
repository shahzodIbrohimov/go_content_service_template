package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	pb "github.com/golanguzb70/go_content_service/genproto/content_service"
	"github.com/golanguzb70/go_content_service/pkg/db"
	"github.com/golanguzb70/go_content_service/pkg/logger"
	"github.com/golanguzb70/go_content_service/storage/repo"
	"github.com/google/uuid"
)

type staff struct {
	db  *db.Postgres
	log logger.Logger
}

func NewStaffRepo(db *db.Postgres, log logger.Logger) repo.StaffI {
	return &staff{
		db:  db,
		log: log,
	}
}

func (r *staff) Create(ctx context.Context, req *pb.Staff) (*pb.Staff, error) {
	req.Lang = "uz" // Har doim Uzbek tilini birinchiga ko'tar

	queryCreate := r.db.Builder.Insert("staffs").
		Columns(`
			id, photo, position_id, slug
		`).
		Values(
			req.Id, req.Photo, req.PositionId, req.Slug,
		).
		Suffix("RETURNING created_at, updated_at")

	queryTranslate := r.db.Builder.Insert("staff_translates").
		Columns(`
			id, lang, first_name, last_name, biography, staff_id
		`).
		Values(uuid.New().String(), req.Lang, req.FirstName, req.LastName, req.Biography, req.Id).
		Values(uuid.New().String(), "ru", req.FirstName, req.LastName, req.Biography, req.Id).
		Values(uuid.New().String(), "en", req.FirstName, req.LastName, req.Biography, req.Id)

	tr := r.db.Db.MustBegin()
	query, args, _ := queryCreate.ToSql()

	err := tr.QueryRow(query, args...).Scan(&req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		tr.Rollback()
		return nil, HandleDatabaseError(err, r.log, "creating staff")
	}

	query, args, _ = queryTranslate.ToSql()
	_, err = tr.Exec(query, args...)
	if err != nil {
		tr.Rollback()
		return nil, HandleDatabaseError(err, r.log, "creating staff translations")
	}

	err = tr.Commit()
	if err != nil {
		return nil, HandleDatabaseError(err, r.log, "committing transaction in create staff")
	}

	return req, nil
}

func (r *staff) Get(ctx context.Context, req *pb.Id) (*pb.Staff, error) {
	res := &pb.Staff{}

	query := r.db.Builder.Select(`
		s.id, photo, position_id, slug, lang, first_name, 
		last_name, biography, s.created_at, s.updated_at
	`).From("staffs s").Join("staff_translates st ON st.staff_id=s.id")

	if req.Id != "" {
		query = query.Where(squirrel.Eq{"s.id": req.Id, "lang": req.Lang})
	} else if req.Slug != "" {
		query = query.Where(squirrel.Eq{"s.slug": req.Id, "lang": req.Lang})
	} else {
		return nil, fmt.Errorf("id is required")
	}

	err := query.RunWith(r.db.Db).QueryRow().Scan(
		&res.Id, &res.Photo, &res.PositionId, &res.Slug, &res.Lang,
		&res.FirstName, &res.LastName, &res.Biography,
		&res.CreatedAt, &res.UpdatedAt,
	)

	return res, HandleDatabaseError(err, r.log, "getting staff")
}

func (r *staff) Find(ctx context.Context, req *pb.GetListFilter) (*pb.Staffs, error) {
	var (
		res            = &pb.Staffs{}
		whereCondition = PrepareWhere(req.Filters)
		orderBy        = PrepareOrder(req.Sorts)
	)

	query := r.db.Builder.Select(`
	s.id, photo, position_id, slug, lang, first_name, 
	last_name, biography, s.created_at, s.updated_at
	`).From("staffs s").
		Join(fmt.Sprintf("staff_translates st ON (st.staff_id=s.id AND lang='%s')", req.Lang))

	if len(req.Filters) > 0 {
		query = query.Where(whereCondition)

	}
	if len(req.Sorts) > 0 {
		query = query.OrderBy(orderBy)
	}

	query = query.Offset(uint64((req.Page - 1) * req.Limit)).Limit(uint64(req.Limit))

	rows, err := query.RunWith(r.db.Db).Query()
	if err != nil {
		return nil, HandleDatabaseError(err, r.log, "error while finding staff")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &pb.Staff{}
		err = rows.Scan(
			&temp.Id, &temp.Photo, &temp.PositionId, &temp.Slug, &temp.Lang,
			&temp.FirstName, &temp.LastName, &temp.Biography,
			&temp.CreatedAt, &temp.UpdatedAt,
		)
		if err != nil {
			return nil, HandleDatabaseError(err, r.log, "error while scanning staff")
		}

		res.Items = append(res.Items, temp)
	}

	count := r.db.Builder.Select("count(1)").
		From("staffs s").
		Join(fmt.Sprintf("staff_translates st ON (st.staff_id=s.id AND lang='%s')", req.Lang)).
		Where(whereCondition)

	err = count.RunWith(r.db.Db).Scan(&res.Count)

	return res, HandleDatabaseError(err, r.log, "getting staff count")
}

func (r *staff) Update(ctx context.Context, req *pb.Staff) (*pb.Staff, error) {
	var (
		mp            = make(map[string]interface{})
		mpTransaltion = make(map[string]interface{})

		whereCondition = squirrel.And{squirrel.Eq{"id": req.Id}}
	)
	mp["photo"] = req.Photo
	mp["position_id"] = req.PositionId
	mp["slug"] = req.Slug
	mp["updated_at"] = time.Now()

	mpTransaltion["updated_at"] = time.Now()
	mpTransaltion["first_name"] = req.FirstName
	mpTransaltion["last_name"] = req.LastName
	mpTransaltion["biography"] = req.Biography

	queryStaff := r.db.Builder.Update("staffs").SetMap(mp).
		Where(whereCondition).
		Suffix("RETURNING updated_at, created_at")

	queryTranslate := r.db.Builder.Update("staff_translates").SetMap(mpTransaltion).
		Where(squirrel.Eq{"staff_id": req.Id, "lang": req.Lang})

	tr := r.db.Db.MustBegin()

	query, args, _ := queryStaff.ToSql()

	err := tr.QueryRow(query, args...).Scan(&req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		tr.Rollback()
		return nil, HandleDatabaseError(err, r.log, "updating staff")
	}

	query, args, _ = queryTranslate.ToSql()
	_, err = tr.Exec(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			queryTranslate := r.db.Builder.Insert("staff_translates").
				Columns(`
				id, lang, first_name, last_name, biography, staff_id
			`).Values(uuid.New().String(), req.Lang, req.FirstName, req.LastName, req.Biography, req.Id)

			query, args, _ = queryTranslate.ToSql()
			_, err = tr.Exec(query, args...)
			if err != nil {
				tr.Rollback()
				return nil, HandleDatabaseError(err, r.log, "creating staff translation")
			}
		} else {
			if err != nil {
				tr.Rollback()
				return nil, HandleDatabaseError(err, r.log, "updating staff translation")
			}
		}
	}

	err = tr.Commit()
	if err != nil {
		return nil, HandleDatabaseError(err, r.log, "committing transaction in create staff")
	}

	return req, HandleDatabaseError(err, r.log, "update staff")
}

func (r *staff) Delete(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	queryTranslate := r.db.Builder.Delete("staff_translates").Where(squirrel.Eq{"staff_id": req.Id})
	queryStaff := r.db.Builder.Delete("staffs").Where(squirrel.Eq{"id": req.Id})

	tr := r.db.Db.MustBegin()
	query, args, _ := queryTranslate.ToSql()
	_, err := tr.Exec(query, args...)
	if err != nil {
		tr.Rollback()
		return nil, HandleDatabaseError(err, r.log, "delete staff translations")
	}

	query, args, _ = queryStaff.ToSql()
	_, err = tr.Exec(query, args...)
	if err != nil {
		tr.Rollback()
		return nil, HandleDatabaseError(err, r.log, "delete staff")
	}

	tr.Commit()
	if err != nil {
		return nil, HandleDatabaseError(err, r.log, "committing transaction in delete staff")
	}

	return &pb.Empty{}, nil
}
