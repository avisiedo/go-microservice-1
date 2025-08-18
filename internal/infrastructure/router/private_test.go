package router

import (
	"testing"

	handler "github.com/avisiedo/go-microservice-1/internal/api/http/private"
	"github.com/avisiedo/go-microservice-1/internal/config"
	presenter "github.com/avisiedo/go-microservice-1/internal/test/mock/api/http/private"
	"github.com/labstack/echo/v4"
)

func helperNewGroupPrivate(t *testing.T) (*echo.Group, *config.Config, handler.ServerInterface) {
	e, cfg := helperNewEchoRouteConfig(t)
	handlers := presenter.NewServerInterface(t)
	return e.Group(privatePath), cfg, handlers
}

func TestNewPrivate(t *testing.T) {
	e, cfg, handlers := helperNewGroupPrivate(t)
	newPrivate(e, cfg, handlers)
}
