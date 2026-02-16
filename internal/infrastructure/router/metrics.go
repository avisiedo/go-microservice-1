package router

import (
	"github.com/avisiedo/go-microservice-1/internal/api/http/metrics"
	metrics_handler "github.com/avisiedo/go-microservice-1/internal/api/http/metrics"
	"github.com/avisiedo/go-microservice-1/internal/config"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	"github.com/labstack/echo/v4"
)

func newGroupMetrics(e *echo.Echo, cfg *config.Config, handlers metrics.ServerInterface) *echo.Echo {
	metrics.RegisterHandlersWithBaseURL(e, handlers, cfg.Metrics.Path)
	return e
}

// NewMetricsRouter fill the routing information for /metrics endpoint.
// e is the echo instance
// cfg is the router configuration
// h is the handler to retrieve the metrics.
// Return the echo instance configured for the metrics for success execution,
// else raise any panic.
func NewMetricsRouter(e *echo.Echo, cfg *config.Config, h metrics_handler.ServerInterface) *echo.Echo {
	newRouterWithConfigCommonGuards(e, cfg)
	if cfg.Metrics.Path == "" {
		panic(common_err.ErrEmpty("cfg.Metrics.Path"))
	}
	if h == nil {
		panic(common_err.ErrNil("h"))
	}

	configCommonMiddlewares(e, cfg)

	// Register handlers
	return newGroupMetrics(e, cfg, h)
}
