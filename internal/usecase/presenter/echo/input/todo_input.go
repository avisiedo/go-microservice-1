package input

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/http/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
)

// TodoInput is the input adapter for Todo resource
type TodoInput struct{}

// Create input adapter for CreateTodo operation
func (i TodoInput) Create(ctx echo.Context, dataInput *model.Todo) error {
	var apiInput public.ToDo
	if err := ctx.Bind(&apiInput); err != nil {
		return fmt.Errorf("binding data: %w", err)
	}
	*dataInput = model.Todo{
		Title:       &apiInput.Title,
		Description: apiInput.Description,
	}
	return nil
}

func (i TodoInput) GetAll(ctx echo.Context) error {
	return nil
}
