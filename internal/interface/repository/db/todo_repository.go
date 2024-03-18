package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *model.Todo) (*model.Todo, error)
	Update(ctx context.Context, todo *model.Todo) (*model.Todo, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Todo, error)
	GetAll(ctx context.Context) ([]model.Todo, error)
	DeleteByUUID(ctx context.Context, uuid uuid.UUID) error
}
