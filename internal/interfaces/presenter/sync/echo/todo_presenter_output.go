package echo

import (
	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	"github.com/labstack/echo/v4"
)

type TodoOutput interface {
	Create(ctx echo.Context, data *model.Todo) (*public.ToDo, error)
	GetAll(ctx echo.Context, data []model.Todo) ([]public.ToDo, error)
	Get(ctx echo.Context, data *model.Todo) (*public.ToDo, error)
}
