package interactor

import (
	"context"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
)

type Todo interface {
	Create(ctx context.Context, todo *model.Todo) (*model.Todo, error)
	GetAll(ctx context.Context) ([]model.Todo, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Todo, error)
	Update(ctx context.Context, todo *model.Todo) error
	Patch(ctx context.Context, todo *model.Todo) (*model.Todo, error)
}
