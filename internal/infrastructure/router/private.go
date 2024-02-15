package router

import (
	"github.com/avisiedo/go-microservice-1/internal/api/http/private"
	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/usecase/interactor"
	presenter "github.com/avisiedo/go-microservice-1/internal/usecase/presenter/echo"
	"github.com/labstack/echo/v4"
)

func newGroupPrivate(e *echo.Group, cfg *config.Config) *echo.Group {
	i := interactor.NewHealthcheck(cfg)
	ph := presenter.NewHealthcheck(i)
	private.RegisterHandlers(e, ph)
	return e
}
