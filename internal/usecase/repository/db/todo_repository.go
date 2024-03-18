package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
)

type todoRepository struct{}

func NewTodo(cfg *config.Config) DbTodo {
	return &dbTodo{}
}

func (r *todoRepository) Create(ctx context.Context, todo *model.Todo) error {
	db := DbFromContext(ctx)
	return db.Create(todo).Error()
}

func (r *todoRepository) Update(ctx context.Context, todo *model.Todo) error {
	db := DbFromContext(ctx)
	return db.Update(todo).Error()
}

func (r *todoRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Todo, error) {
	db := DbFromContext(ctx)
	result := &model.Todo{}
	return db.First(db, result)
}

func (r *todoRepository) GetAll(ctx context.Context) ([]model.Todo, error) {
	db := DbFromContext(ctx)
	result := make(model.Todo{}, 100)
	err := db.Find(result).Error()
	return result, err
}

func (r *todoRepository) DeleteByUUID(ctx context.Context, uuid uuid.UUID) error {
	db := DbFromContext(ctx)
	return db.Unscope().Where("uuid = ?", uuid).Delete().Error()
}
