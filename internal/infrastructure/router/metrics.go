package router

import (
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/http/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/config"
)

func newGroupMetrics(e *echo.Echo, cfg *config.Config, handlers metrics.ServerInterface) *echo.Echo {
	metrics.RegisterHandlersWithBaseURL(e, handlers, cfg.Metrics.Path)
	return e
}
