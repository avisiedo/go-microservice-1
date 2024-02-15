package input

import (
	"fmt"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	"github.com/labstack/echo/v4"
)

// TodoInput is the input adapter for Todo resource
type TodoInput struct{}

// Create input adapter for CreateTodo operation
func (i TodoInput) Create(ctx echo.Context) (*model.Todo, error) {
	var apiInput public.ToDo
	if err := ctx.Bind(&apiInput); err != nil {
		return nil, fmt.Errorf("binding data: %w", err)
	}
	data := &model.Todo{
		Title:       &apiInput.Title,
		Description: apiInput.Description,
	}
	return data, nil
}

// GetAll input adapter for GetAll operation
func (i TodoInput) GetAll(ctx echo.Context) error {
	return nil
}
