package db

import (
	"context"

	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	repository "github.com/avisiedo/go-microservice-1/internal/interface/repository/db"
	"github.com/google/uuid"
)

type todoRepository struct{}

func NewTodo(cfg *config.Config) repository.TodoRepository {
	return &todoRepository{}
}

func (r *todoRepository) Create(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	db := DbFromContext(ctx)
	if err := db.Create(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	db := DbFromContext(ctx)
	if err := db.Updates(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *todoRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Todo, error) {
	db := DbFromContext(ctx)
	result := &model.Todo{}
	if err := db.First(db, result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *todoRepository) GetAll(ctx context.Context) ([]model.Todo, error) {
	// TODO refactor to support paginated results
	var count int64
	db := DbFromContext(ctx)
	if err := db.Model(&model.Todo{}).Count(&count).Error; err != nil {
		return []model.Todo{}, err
	}
	if count > 0 {
		result := make([]model.Todo, count)
		err := db.Count(&count).Find(&model.Todo{}).Error
		return result, err
	}
	return []model.Todo{}, nil
}

func (r *todoRepository) DeleteByUUID(ctx context.Context, uuid uuid.UUID) error {
	db := DbFromContext(ctx)
	return db.Unscoped().Delete(&model.Todo{}, "uuid = ?", uuid).Error
}
