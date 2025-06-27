package echo

import (
	"fmt"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	. "github.com/avisiedo/go-microservice-1/internal/interface/presenter/sync/echo"
	"github.com/labstack/echo/v4"
)

// TodoInput is the input adapter for Todo resource
type todoInput struct{}

// NewTodoInput create a new input adapter.
func NewTodoInput() TodoInput {
	return newTodoInput()
}

func newTodoInput() *todoInput {
	return &todoInput{}
}

// Create input adapter for CreateTodo operation
func (i *todoInput) Create(ctx echo.Context) (*model.Todo, error) {
	var apiInput public.ToDo
	if ctx.Request().Body == nil {
		return nil, common_err.ErrNil("Body")
	}
	if err := ctx.Bind(&apiInput); err != nil {
		return nil, fmt.Errorf("binding request data")
	}
	data := &model.Todo{
		Title:       apiInput.Title,
		Description: apiInput.Description,
		DueDate:     apiInput.DueDate,
	}
	if data.Title == "" {
		return nil, common_err.ErrEmpty("data.Title")
	}
	if data.Description == "" {
		return nil, common_err.ErrEmpty("data.Description")
	}

	return data, nil
}

// GetAll input adapter for GetAll operation
func (i *todoInput) GetAll(ctx echo.Context) error {
	if len(ctx.QueryParams()) > 0 {
		return fmt.Errorf("No query parameters expected for " + ctx.Path())
	}
	return nil
}
