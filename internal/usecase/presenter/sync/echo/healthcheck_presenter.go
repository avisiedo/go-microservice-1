package echo

import (
	"log/slog"
	"net/http"

	api_healthcheck "github.com/avisiedo/go-microservice-1/internal/api/http/healthcheck"
	app_context "github.com/avisiedo/go-microservice-1/internal/infrastructure/context"
	interactor "github.com/avisiedo/go-microservice-1/internal/interface/interactor"
	presenter "github.com/avisiedo/go-microservice-1/internal/interface/presenter/sync/echo"
	"github.com/labstack/echo/v4"
)

type healthcheckPresenter struct {
	interactor interactor.HealthcheckInteractor
}

func NewHealthcheck(i interactor.HealthcheckInteractor) presenter.Healthcheck {
	if i == nil {
		panic("the interactor is nil")
	}
	return &healthcheckPresenter{
		interactor: i,
	}
}

// Liveness kubernetes probe endpoint
// (GET /livez)
func (p *healthcheckPresenter) GetLivez(c echo.Context) error {
	if c == nil {
		slog.Error("echo context is nil")
		return echo.ErrInternalServerError
	}
	if err := p.interactor.IsLive(); err != nil {
		ctx := c.Request().Context()
		app_context.LogFromContext(ctx).
			ErrorContext(ctx, err.Error())
		if err2 := c.String(http.StatusServiceUnavailable, api_healthcheck.NotHealthy); err2 != nil {
			slog.ErrorContext(ctx, err2.Error())
		}
		return echo.ErrServiceUnavailable
	}
	return c.String(http.StatusOK, api_healthcheck.Healthy)
}

// Readiness kubernetes probe endpoint
// (GET /readyz)
func (p *healthcheckPresenter) GetReadyz(c echo.Context) error {
	if c == nil {
		slog.Error("echo context is nil")
		return echo.ErrInternalServerError
	}
	if err := p.interactor.IsReady(); err != nil {
		ctx := c.Request().Context()
		app_context.LogFromContext(ctx).
			ErrorContext(ctx, err.Error())
		if err2 := c.String(http.StatusServiceUnavailable, api_healthcheck.NotReady); err2 != nil {
			slog.ErrorContext(ctx, err2.Error())
		}
		return echo.ErrServiceUnavailable
	}
	return c.String(http.StatusOK, api_healthcheck.Ready)
}
