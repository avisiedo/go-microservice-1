package interactor

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	db "github.com/podengo-project/idmsvc-backend/internal/interface/repository/db"
)

type todoInteractor struct {
	r db.TodoRepository
}

var ErrNotImplemented = errors.New("not implemented")

func NewTodo(repo db.TodoRepository) interactor.Todo {
	return &todoInteractor{
		r: repo,
	}
}

func (i *todoInteractor) Create(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	if ctx == nil {
		return nil, errors.New("'ctx' is nil")
	}
	if todo == nil {
		return nil, errors.New("'todo' is nil")
	}
	return i.r.Create(ctx, todo)
}

func (i *todoInteractor) Update(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	if ctx == nil {
		return nil, errors.New("'ctx' is nil")
	}
	if todo == nil {
		return nil, errors.New("'todo' is nil")
	}
	return i.r.Create(ctx, todo)
}

func (i *todoInteractor) GetByUUID(ctx context.Context, id uuid.UUID) (*model.Todo, error) {
	if ctx == nil {
		return nil, errors.New("'ctx' is nil")
	}
	if (id == uuid.UUID{}) {
		return nil, errors.New("'uuid' is empty")
	}
	return i.r.GetByUUID(ctx, id)
}

func (i *todoInteractor) GetAll(ctx context.Context) ([]model.Todo, error) {
	var (
		err    error
		result []model.Todo
	)
	if ctx == nil {
		return nil, errors.New("'ctx' is nil")
	}
	if result, err = i.r.GetAll(ctx); err != nil {
		return nil, err
	}
	return result, nil
}

func (i *todoInteractor) DeleteByUUID(ctx context.Context, id uuid.UUID) error {
	if ctx == nil {
		return errors.New("'ctx' is nil")
	}
	if (id == uuid.UUID{}) {
		return errors.New("'id' is nil")
	}
	return nil
}

func (i *todoInteractor) Patch(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	var (
		newTodo *model.Todo
		err     error
	)
	if newTodo, err = i.GetByUUID(ctx, todo.UUID); err != nil {
		return nil, err
	}
	if todo.Description != nil {
		newTodo.Description = todo.Description
	}
	if todo.DueDate != nil {
		newTodo.DueDate = todo.DueDate
	}
	if todo.Title != nil {
		newTodo.Title = todo.Title
	}
	if newTodo, err = i.r.Update(ctx, newTodo); err != nil {
		return nil, err
	}
	return newTodo, nil
}
