package todo

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type Repository interface {
	Create(ctx context.Context, t *Todo) error
	List(ctx context.Context, f ListFilter) ([]Todo, error)
	FindByID(ctx context.Context, id int64) (Todo, error)
	Update(ctx context.Context, t *Todo) error
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, t *Todo) error {
	_, err := r.db.NewInsert().Model(t).Returning("*").Exec(ctx)
	return err
}

func (r *repository) List(ctx context.Context, f ListFilter) ([]Todo, error) {
	var todos []Todo
	q := r.db.NewSelect().Model(&todos).Order("id ASC")

	if f.Q != "" {
		q = q.Where("title ILIKE ? OR description ILIKE ?", "%"+f.Q+"%", "%"+f.Q+"%")
	}
	if f.IsCompleted != nil {
		q = q.Where("is_completed = ?", *f.IsCompleted)
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	if f.Offset > 0 {
		q = q.Offset(f.Offset)
	}

	err := q.Scan(ctx)
	return todos, err
}

func (r *repository) FindByID(ctx context.Context, id int64) (Todo, error) {
	var t Todo
	err := r.db.NewSelect().Model(&t).Where("id = ?", id).Scan(ctx)
	if err == sql.ErrNoRows {
		return Todo{}, sql.ErrNoRows
	}
	return t, err
}

func (r *repository) Update(ctx context.Context, t *Todo) error {
	_, err := r.db.NewUpdate().Model(t).Where("id = ?", t.ID).Returning("*").Exec(ctx)
	return err
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.NewDelete().Model((*Todo)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
