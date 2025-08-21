package repository

import (
	"context"

	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	"github.com/google/uuid"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *model.Todo) (*model.Todo, error)
	Update(ctx context.Context, todo *model.Todo) (*model.Todo, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Todo, error)
	GetAll(ctx context.Context) ([]model.Todo, error)
	DeleteByUUID(ctx context.Context, uuid uuid.UUID) error
}
