package router

import (
	"github.com/avisiedo/go-microservice-1/internal/api/http/private"
	handler "github.com/avisiedo/go-microservice-1/internal/api/http/private"
	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/labstack/echo/v4"
)

func newPrivate(e *echo.Group, cfg *config.Config, app handler.ServerInterface) *echo.Group {
	private.RegisterHandlers(e.Group(privatePath), app)
	return e
}
