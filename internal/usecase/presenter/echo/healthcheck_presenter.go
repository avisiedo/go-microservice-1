package presenter

import (
	"net/http"

	"github.com/labstack/echo/v4"
	interactor "github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	presenter "github.com/podengo-project/idmsvc-backend/internal/interface/presenter/echo"
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
	if err := p.GetLivez(ctx); err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}
	return ctx.JSON(http.StatusOK, "Livez")
}

// Readiness kubernetes probe endpoint
// (GET /readyz)
func (p *healthcheckPresenter) GetReadyz(ctx echo.Context) error {
	if err := p.GetReadyz(ctx); err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}
	return ctx.JSON(http.StatusOK, "Readyz")
}
