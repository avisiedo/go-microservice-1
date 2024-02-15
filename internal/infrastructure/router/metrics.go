package router

import (
	"github.com/avisiedo/go-microservice-1/internal/api/http/metrics"
	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/labstack/echo/v4"
)

func newGroupMetrics(e *echo.Echo, cfg *config.Config, handlers metrics.ServerInterface) *echo.Echo {
	metrics.RegisterHandlersWithBaseURL(e, handlers, cfg.Metrics.Path)
	return e
}
