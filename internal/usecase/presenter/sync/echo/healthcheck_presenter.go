package echo

import (
	"log/slog"
	"net/http"

	api_healthcheck "github.com/avisiedo/go-microservice-1/internal/api/http/healthcheck"
	"github.com/avisiedo/go-microservice-1/internal/errors/common"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
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
		panic(common_err.ErrNil("i"))
	}
	return &healthcheckPresenter{
		interactor: i,
	}
}

// Liveness kubernetes probe endpoint
// (GET /livez)
func (p *healthcheckPresenter) GetLivez(ctx echo.Context) error {
	if ctx == nil {
		slog.Error(common_err.ErrNil("ctx").Error())
		return echo.ErrInternalServerError
	}
	if err := p.interactor.IsLive(); err != nil {
		newCtx := ctx.Request().Context()
		app_context.LogFromContext(newCtx).
			ErrorContext(newCtx, err.Error())
		if err2 := ctx.String(http.StatusServiceUnavailable, api_healthcheck.NotHealthy); err2 != nil {
			slog.ErrorContext(newCtx, err2.Error())
		}
		return echo.ErrServiceUnavailable
	}
	return ctx.String(http.StatusOK, api_healthcheck.Healthy)
}

// Readiness kubernetes probe endpoint
// (GET /readyz)
func (p *healthcheckPresenter) GetReadyz(ctx echo.Context) error {
	if ctx == nil {
		slog.Default().Error(common.ErrNil("ctx").Error())
		return echo.ErrInternalServerError
	}
	if err := p.interactor.IsReady(); err != nil {
		newCtx := ctx.Request().Context()
		app_context.LogFromContext(newCtx).
			ErrorContext(newCtx, err.Error())
		if err2 := ctx.String(http.StatusServiceUnavailable, api_healthcheck.NotReady); err2 != nil {
			slog.ErrorContext(newCtx, err2.Error())
		}
		return echo.ErrServiceUnavailable
	}
	return ctx.String(http.StatusOK, api_healthcheck.Ready)
}
