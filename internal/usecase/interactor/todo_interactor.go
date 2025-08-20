package interactor

import (
	"context"

	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	"github.com/avisiedo/go-microservice-1/internal/interface/interactor"
	db "github.com/avisiedo/go-microservice-1/internal/interface/repository/db"
	"github.com/google/uuid"
)

type todoInteractor struct {
	todoDB db.TodoRepository
}

func newTodo(repo db.TodoRepository) *todoInteractor {
	if repo == nil {
		panic(common_err.ErrNil("repo"))
	}
	return &todoInteractor{
		todoDB: repo,
	}
}

// NewTodo instantiate a new Todo interactor
// repo is a not nil instance that implement db.TodoRepository.
// Return an instance that accomplish the interactor.Todo interface.
func NewTodo(repo db.TodoRepository) interactor.Todo {
	return newTodo(repo)
}

func (i *todoInteractor) checkCtx(ctx context.Context) error {
	if ctx == nil {
		return common_err.ErrNil("ctx")
	}
	return nil
}

func (i *todoInteractor) checkCtxAndTodo(ctx context.Context, todo *model.Todo) error {
	if err := i.checkCtx(ctx); err != nil {
		return err
	}
	if todo == nil {
		return common_err.ErrNil("todo")
	}
	return nil
}

func (i *todoInteractor) checkCtxAndId(ctx context.Context, id uuid.UUID) error {
	if err := i.checkCtx(ctx); err != nil {
		return err
	}
	if (id == uuid.UUID{}) {
		return common_err.ErrEmpty("id")
	}
	return nil
}

func (i *todoInteractor) checkCreateTodo(ctx context.Context, todo *model.Todo) error {
	if err := i.checkCtxAndTodo(ctx, todo); err != nil {
		return err
	}
	if todo.Description == "" {
		return common_err.ErrEmpty("todo.Description")
	}
	return nil
}

func (i *todoInteractor) Create(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	if err := i.checkCreateTodo(ctx, todo); err != nil {
		return nil, err
	}
	return i.todoDB.Create(ctx, todo)
}

func (i *todoInteractor) Update(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	if err := i.checkCtxAndTodo(ctx, todo); err != nil {
		return nil, err
	}
	return i.todoDB.Update(ctx, todo)
}

func (i *todoInteractor) GetByUUID(ctx context.Context, id uuid.UUID) (*model.Todo, error) {
	if err := i.checkCtxAndId(ctx, id); err != nil {
		return nil, err
	}
	return i.todoDB.GetByUUID(ctx, id)
}

func (i *todoInteractor) GetAll(ctx context.Context) ([]model.Todo, error) {
	var (
		err    error
		result []model.Todo
	)
	if err = i.checkCtx(ctx); err != nil {
		return nil, err
	}
	if result, err = i.todoDB.GetAll(ctx); err != nil {
		return nil, err
	}
	return result, nil
}

func (i *todoInteractor) DeleteByUUID(ctx context.Context, id uuid.UUID) error {
	if err := i.checkCtxAndId(ctx, id); err != nil {
		return err
	}
	return i.todoDB.DeleteByUUID(ctx, id)
}

func (i *todoInteractor) Patch(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	var (
		newTodo *model.Todo
		err     error
	)
	if err = i.checkCtxAndTodo(ctx, todo); err != nil {
		return nil, err
	}
	if newTodo, err = i.GetByUUID(ctx, todo.UUID); err != nil {
		return nil, err
	}
	newTodo.Title = todo.Title
	newTodo.Description = todo.Description
	if todo.DueDate != nil {
		newTodo.DueDate = todo.DueDate
	}
	if newTodo, err = i.todoDB.Update(ctx, newTodo); err != nil {
		return nil, err
	}
	return newTodo, nil
}
