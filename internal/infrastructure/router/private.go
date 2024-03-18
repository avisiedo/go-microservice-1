package router

import (
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/http/private"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/usecase/interactor"
	presenter "github.com/podengo-project/idmsvc-backend/internal/usecase/presenter/echo"
)

func newGroupPrivate(e *echo.Group, cfg *config.Config) *echo.Group {
	ph := presenter.NewHealthcheck(interactor.NewHealthcheck(cfg))
	private.RegisterHandlers(e)
	return e
}
