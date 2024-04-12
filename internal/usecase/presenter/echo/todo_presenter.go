package echo

import (
	"context"
	"errors"
	"net/http"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/logger/slogctx"
	"github.com/avisiedo/go-microservice-1/internal/interface/interactor"
	presenter "github.com/avisiedo/go-microservice-1/internal/interface/presenter/echo"
	"github.com/avisiedo/go-microservice-1/internal/usecase/presenter/echo/input"
	"github.com/avisiedo/go-microservice-1/internal/usecase/presenter/echo/output"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/exp/slog"
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
func (p *todoPresenter) GetAllTodos(c echo.Context) error {
	var (
		todos  []model.Todo
		output []public.ToDo
		err    error
	)
	ctx := c.Request().Context()
	l := slogctx.FromCtx(ctx)
	// l := slog.Default()
	if err = p.input.GetAll(c); err != nil {
		l.ErrorContext(ctx, "presenter input adapter error at GetAll(): %s", err.Error())
		return err
	}
	if err = p.db.Transaction(func(tx *gorm.DB) error {
		var err error
		c := context.WithValue(ctx, "db", tx)
		if todos, err = p.interactor.GetAll(c); err != nil {
			l.ErrorContext(ctx, "presenter error at GetAll(): %s", err.Error())
			return err
		}
		return nil
	}); err != nil {
		l.ErrorContext(ctx, "transaction error at GetAll(): %s", err.Error())
		return err
	}
	if output, err = p.output.GetAll(c, todos); err != nil {
		l.ErrorContext(ctx, "presenter output adapter error at GetAll(): %s", err.Error())
		return err
	}
	return c.JSON(http.StatusOK, output)
}

// Create a new ToDo item
// (POST /todos)
func (p *todoPresenter) CreateTodo(ctx echo.Context) error {
	var (
		output *public.ToDo
		data   *model.Todo
		err    error
	)
	l, ok := ctx.Get("log").(*slog.Logger)
	if !ok || l == nil {
		return errors.New("'log' is undefined in the echo context")
	}
	if data, err = p.input.Create(ctx); err != nil {
		return err
	}
	if err = p.db.Transaction(func(tx *gorm.DB) error {
		var err error
		c := context.WithValue(
			context.WithValue(
				ctx.Request().Context(), "db", tx),
			"log", l,
		)
		if data, err = p.interactor.Create(c, data); err != nil {
			l.ErrorContext(ctx.Request().Context(), "presenter error on invoking interactor.Create: %s", err.Error())
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	if output, err = p.output.Create(ctx, data); err != nil {
		l.ErrorContext(ctx.Request().Context(), "presenter output adapter error at GetAll(): %s", err.Error())
		return err
	}
	return ctx.JSON(http.StatusCreated, output)
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
