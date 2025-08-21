package echo

import (
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	"github.com/labstack/echo/v4"
)

type TodoInput interface {
	Create(ctx echo.Context) (*model.Todo, error)
	GetAll(ctx echo.Context) error
}
