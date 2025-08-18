package router

import (
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/api/http/metrics"
	"github.com/avisiedo/go-microservice-1/internal/config"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	presenter "github.com/avisiedo/go-microservice-1/internal/test/mock/api/http/metrics"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func helperNewGroupMetrics(t *testing.T) (*echo.Echo, *config.Config, metrics.ServerInterface) {
	e, cfg := helperNewEchoRouteConfig(t)
	presenterMetrics := presenter.NewServerInterface(t)
	return e, cfg, presenterMetrics
}

func TestNewGroupMetrics(t *testing.T) {
	e, cfg, presenterMetrics := helperNewGroupMetrics(t)
	newGroupMetrics(e, cfg, presenterMetrics)
}

func TestNewMetricsRouter(t *testing.T) {
	e, cfg, presenterMetrics := helperNewGroupMetrics(t)
	cfg.Metrics.Path = ""
	assert.PanicsWithError(t, common_err.ErrEmpty("cfg.Metrics.Path").Error(), func() {
		NewMetricsRouter(e, cfg, nil)
	})

	cfg.Metrics.Path = "/metrics"

	assert.PanicsWithError(t, common_err.ErrNil("h").Error(), func() {
		_ = NewMetricsRouter(e, cfg, nil)
	})

	assert.NotPanics(t, func() {
		e = NewMetricsRouter(e, cfg, presenterMetrics)
	})
}
