package db

import (
	"context"

	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	app_context "github.com/avisiedo/go-microservice-1/internal/infrastructure/context"
	repository "github.com/avisiedo/go-microservice-1/internal/interface/repository/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type todoRepository struct{}

func NewTodo(cfg *config.Config) repository.TodoRepository {
	return &todoRepository{}
}

func (r *todoRepository) Create(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	var (
		db  *gorm.DB
		err error
	)
	if db, err = app_context.DBFromContext(ctx); err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, common_err.ErrNil("todo")
	}
	if (todo.UUID == uuid.UUID{}) {
		todo.UUID = uuid.New()
	}
	if err := db.Create(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	var (
		db  *gorm.DB
		err error
	)
	if db, err = app_context.DBFromContext(ctx); err != nil {
		return nil, err
	}
	if err := db.Updates(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *todoRepository) GetByUUID(ctx context.Context, id uuid.UUID) (*model.Todo, error) {
	var (
		db  *gorm.DB
		err error
	)
	if db, err = app_context.DBFromContext(ctx); err != nil {
		return nil, err
	}
	if (id == uuid.UUID{}) {
		return nil, common_err.ErrEmpty("id")
	}
	result := &model.Todo{}
	if err := db.First(result, "uuid = ?", id).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *todoRepository) GetAll(ctx context.Context) ([]model.Todo, error) {
	// TODO refactor to support paginated results
	//      results must be ordered
	//      results must be limited
	//      results should start on an initial item
	var (
		db    *gorm.DB
		err   error
		count int64
	)
	if db, err = app_context.DBFromContext(ctx); err != nil {
		return nil, err
	}
	if err = db.Model(&model.Todo{}).Count(&count).Error; err != nil {
		return []model.Todo{}, err
	}
	if count > 0 {
		result := make([]model.Todo, count)
		err = db.Find(&result).Error
		return result, err
	}
	return []model.Todo{}, nil
}

func (r *todoRepository) DeleteByUUID(ctx context.Context, todo_uuid uuid.UUID) error {
	var (
		db  *gorm.DB
		err error
	)
	if db, err = app_context.DBFromContext(ctx); err != nil {
		return err
	}
	if (todo_uuid == uuid.UUID{}) {
		return common_err.ErrEmpty("todo_uuid")
	}
	return db.Unscoped().Delete(&model.Todo{}, "uuid = ?", todo_uuid).Error
}
