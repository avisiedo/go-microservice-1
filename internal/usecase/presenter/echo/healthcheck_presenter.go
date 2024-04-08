package echo

import (
	"net/http"

	interactor "github.com/avisiedo/go-microservice-1/internal/interface/interactor"
	presenter "github.com/avisiedo/go-microservice-1/internal/interface/presenter/echo"
	"github.com/labstack/echo/v4"
)

type healthcheckPresenter struct {
	interactor interactor.HealthcheckInteractor
}

func NewHealthcheck(i interactor.HealthcheckInteractor) presenter.Healthcheck {
	return &healthcheckPresenter{
		interactor: i,
	}
}

// Liveness kubernetes probe endpoint
// (GET /livez)
func (p *healthcheckPresenter) GetLivez(ctx echo.Context) error {
	if err := p.interactor.IsLive(); err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}
	return ctx.JSON(http.StatusOK, "Livez")
}

// Readiness kubernetes probe endpoint
// (GET /readyz)
func (p *healthcheckPresenter) GetReadyz(ctx echo.Context) error {
	if err := p.interactor.IsReady(); err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}
	return ctx.JSON(http.StatusOK, "Readyz")
}
