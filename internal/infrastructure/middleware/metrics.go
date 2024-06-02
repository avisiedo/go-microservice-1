package middleware

import (
	"strconv"
	"time"

	"github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsConfig struct {
	Skipper echo_middleware.Skipper
	Metrics *metrics.Metrics
}

var defaultConfig MetricsConfig = MetricsConfig{
	Skipper: echo_middleware.DefaultSkipper,
	Metrics: metrics.NewMetrics(prometheus.NewRegistry()),
}

func MetricsMiddlewareWithConfig(config *MetricsConfig) echo.MiddlewareFunc {
	if config == nil {
		config = &defaultConfig
	}
	if config.Skipper == nil {
		config.Skipper = echo_middleware.DefaultSkipper
	}
	if config.Metrics == nil {
		panic("config.Metrics can not be nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()
			if config.Skipper(ctx) {
				return next(ctx)
			}
			method := ctx.Request().Method
			path := MatchedRoute(ctx)
			err := next(ctx)
			status := strconv.Itoa(ctx.Response().Status)
			config.Metrics.HttpStatusHistogram.WithLabelValues(status, method, path).Observe(time.Since(start).Seconds())
			return err
		}
	}
}

func CreateMetricsMiddleware(metrics *metrics.Metrics) echo.MiddlewareFunc {
	return MetricsMiddlewareWithConfig(
		&MetricsConfig{
			Metrics: metrics,
		})
}
