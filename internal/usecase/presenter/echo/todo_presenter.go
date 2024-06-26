package echo

import (
	"net/http"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	app_context "github.com/avisiedo/go-microservice-1/internal/infrastructure/context"
	"github.com/avisiedo/go-microservice-1/internal/interface/interactor"
	presenter "github.com/avisiedo/go-microservice-1/internal/interface/presenter/echo"
	"github.com/avisiedo/go-microservice-1/internal/usecase/presenter/echo/input"
	"github.com/avisiedo/go-microservice-1/internal/usecase/presenter/echo/output"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

type todoPresenter struct {
	db         *gorm.DB
	interactor interactor.Todo
	input      input.TodoInput
	output     output.TodoOutput
}

func NewTodo(cfg *config.Config, i interactor.Todo, db *gorm.DB) presenter.Todo {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if i == nil {
		panic("interactor is nil")
	}
	if db == nil {
		panic("'db' is nil")
	}
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
	if err = p.input.GetAll(c); err != nil {
		app_context.LogFromContext(ctx).
			ErrorContext(ctx, err.Error())
		return err
	}
	if err = p.db.Transaction(func(tx *gorm.DB) error {
		var err error
		c := app_context.WithDB(ctx, tx)
		if todos, err = p.interactor.GetAll(c); err != nil {
			app_context.LogFromContext(ctx).
				ErrorContext(ctx, err.Error())
			return err
		}
		return nil
	}); err != nil {
		app_context.LogFromContext(ctx).
			ErrorContext(ctx, err.Error())
		return err
	}
	if output, err = p.output.GetAll(c, todos); err != nil {
		app_context.LogFromContext(ctx).
			ErrorContext(ctx, err.Error())
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

	if data, err = p.input.Create(ctx); err != nil {
		return err
	}
	if data == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "empty todo data")
	}
	if data.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "title is an empty string")
	}
	if data.Description == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "description is an empty string")
	}
	if err = p.db.Transaction(func(tx *gorm.DB) error {
		c := app_context.WithDB(ctx.Request().Context(), tx)
		if data, err = p.interactor.Create(c, data); err != nil {
			app_context.LogFromContext(c).
				ErrorContext(c, err.Error())
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	if output, err = p.output.Create(ctx, data); err != nil {
		c := ctx.Request().Context()
		app_context.LogFromContext(c).
			ErrorContext(c, err.Error())
		return err
	}
	return ctx.JSON(http.StatusCreated, output)
}

// Remove item by ID
// (DELETE /todos/{todoId})
func (p *todoPresenter) DeleteTodo(c echo.Context, todoId openapi_types.UUID) error {
	ctx := c.Request().Context()
	app_context.LogFromContext(ctx).
		ErrorContext(ctx, "not implemented")
	return echo.ErrNotImplemented
}

// Get item by ID
// (GET /todos/{todoId})
func (p *todoPresenter) GetTodo(c echo.Context, todoId openapi_types.UUID) error {
	var (
		todo   *model.Todo
		output *public.ToDo
		err    error
	)
	ctx := c.Request().Context()
	if err = p.db.Transaction(func(tx *gorm.DB) error {
		var err error
		c := app_context.WithDB(ctx, tx)
		if todo, err = p.interactor.GetByUUID(c, todoId); err != nil {
			app_context.LogFromContext(ctx).
				ErrorContext(ctx, err.Error())
			return err
		}
		return nil
	}); err != nil {
		app_context.LogFromContext(ctx).
			ErrorContext(ctx, err.Error())
		return err
	}
	if output, err = p.output.Get(c, todo); err != nil {
		app_context.LogFromContext(ctx).
			ErrorContext(ctx, err.Error())
		return err
	}
	return c.JSON(http.StatusOK, output)
}

// Patch an existing ToDo item
// (PATCH /todos/{todoId})
func (p *todoPresenter) PatchTodo(c echo.Context, todoId openapi_types.UUID) error {
	ctx := c.Request().Context()
	app_context.LogFromContext(ctx).
		ErrorContext(ctx, "not implemented")
	return echo.ErrNotImplemented
}

// Substitute an existing ToDo item
// (PUT /todos/{todoId})
func (p *todoPresenter) UpdateTodo(c echo.Context, todoId openapi_types.UUID) error {
	ctx := c.Request().Context()
	app_context.LogFromContext(ctx).
		ErrorContext(ctx, "not implemented")
	return echo.ErrNotImplemented
}
