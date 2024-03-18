package presenter

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/http/public"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	presenter "github.com/podengo-project/idmsvc-backend/internal/interface/presenter/echo"
	"github.com/podengo-project/idmsvc-backend/internal/usecase/presenter/echo/input"
	"github.com/podengo-project/idmsvc-backend/internal/usecase/presenter/echo/output"
	"gorm.io/gorm"
)

type todoPresenter struct {
	db         *gorm.DB
	interactor interactor.Todo
	input      input.TodoInput
	output     output.TodoOutput
}

var ErrNotImplemented = errors.New("not implemented")

func NewTodo(cfg *config.Config, i interactor.Todo, db *gorm.DB) presenter.Todo {
	return &todoPresenter{
		db:         db,
		interactor: i,
	}
}

// Retrieve all ToDo items
// (GET /todos)
func (p *todoPresenter) GetAllTodos(ctx echo.Context) error {
	var (
		result    []model.Todo
		err       error
		apiInput  public.ToDo
		dataInput model.Todo
	)
	log, ok := ctx.Get("log").(*log.Logger)
	if !ok || log == nil {
		return errors.New("'log' is undefined in the echo context")
	}
	if err = p.input.GetAll(ctx); err != nil {
		return err
	}
	if err = p.db.Transaction(func(tx *gorm.DB) error {
		var err error
		c := context.WithValue(
			context.WithValue(
				ctx.Request().Context(), "db", tx),
			"log", log,
		)
		if result, err = p.interactor.GetAll(c); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, result)
}

// Create a new ToDo item
// (POST /todos)
func (p *todoPresenter) CreateTodo(ctx echo.Context) error {
	var (
		result []model.Todo
		err    error
	)
	log, ok := ctx.Get("log").(*log.Logger)
	if !ok || log == nil {
		return errors.New("'log' is undefined in the echo context")
	}
	if err = p.db.Transaction(func(tx *gorm.DB) error {
		var err error
		c := context.WithValue(
			context.WithValue(
				ctx.Request().Context(), "db", tx),
			"log", log,
		)
		if result, err = p.interactor.Create(c); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, result)

	tx := p.db.Begin()
	defer tx.Rollback()
	c := context.WithValue(ctx.Request().Context(), "db", tx)
	result, err := p.interactor.Create(c)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	tx.Commit()
	return ctx.JSON(http.StatusCreated, result)
}

// Remove item by ID
// (DELETE /todos/{todoId})
func (p *todoPresenter) DeleteTodo(ctx echo.Context, todoId openapi_types.UUID) error {
	return ErrNotImplemented
}

// Get item by ID
// (GET /todos/{todoId})
func (p *todoPresenter) GetTodo(ctx echo.Context, todoId openapi_types.UUID) error {
	return ErrNotImplemented
}

// Patch an existing ToDo item
// (PATCH /todos/{todoId})
func (p *todoPresenter) PatchTodo(ctx echo.Context, todoId openapi_types.UUID) error {
	return ErrNotImplemented
}

// Substitute an existing ToDo item
// (PUT /todos/{todoId})
func (p *todoPresenter) UpdateTodo(ctx echo.Context, todoId openapi_types.UUID) error {
	return ErrNotImplemented
}
