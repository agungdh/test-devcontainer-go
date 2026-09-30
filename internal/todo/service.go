package todo

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("todo not found")
var ErrEmptyUpdate = errors.New("minimal satu field harus diisi")

type Service interface {
	Create(ctx context.Context, req CreateTodoRequest) (Todo, error)
	List(ctx context.Context, f ListFilter) ([]Todo, error)
	Get(ctx context.Context, id int64) (Todo, error)
	Update(ctx context.Context, id int64, req UpdateTodoRequest) (Todo, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateTodoRequest) (Todo, error) {
	t := Todo{
		Title:       req.Title,
		Description: req.Description,
	}
	if req.IsCompleted != nil {
		t.IsCompleted = *req.IsCompleted
	}
	if err := s.repo.Create(ctx, &t); err != nil {
		return Todo{}, err
	}
	return t, nil
}

func (s *service) List(ctx context.Context, f ListFilter) ([]Todo, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	return s.repo.List(ctx, f)
}

func (s *service) Get(ctx context.Context, id int64) (Todo, error) {
	t, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Todo{}, ErrNotFound
	}
	return t, err
}

func (s *service) Update(ctx context.Context, id int64, req UpdateTodoRequest) (Todo, error) {
	if req.Empty() {
		return Todo{}, ErrEmptyUpdate
	}
	t, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Todo{}, ErrNotFound
	}
	if err != nil {
		return Todo{}, err
	}

	if req.Title != nil {
		t.Title = *req.Title
	}
	if req.Description != nil {
		t.Description = *req.Description
	}
	if req.IsCompleted != nil {
		t.IsCompleted = *req.IsCompleted
	}

	if err := s.repo.Update(ctx, &t); err != nil {
		return Todo{}, err
	}
	return t, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
