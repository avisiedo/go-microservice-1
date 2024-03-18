package router

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	handler "github.com/podengo-project/idmsvc-backend/internal/handler/http"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
)

type RouterConfig struct {
	Handlers           handler.Application
	PublicPath         string
	PrivatePath        string
	Version            string
	MetricsPath        string
	IsFakeEnabled      bool
	EnableAPIValidator bool
	Metrics            *metrics.Metrics
}

const (
	privatePath    = "/internal"
	basepublicPath = "/api/todo/v1"
)

func getMajorVersion(version string) string {
	if version == "" {
		return ""
	}
	return strings.Split(version, ".")[0]
}

func loggerSkipperWithPaths(paths ...string) middleware.Skipper {
	return func(c echo.Context) bool {
		path := c.Path()
		for _, item := range paths {
			if item == path {
				return true
			}
		}
		return false
	}
}

func configCommonMiddlewares(e *echo.Echo, cfg *config.Config) {
	e.Pre(middleware.RemoveTrailingSlash())

	skipperPaths := []string{
		privatePath + "/readyz",
		privatePath + "/livez",
		cfg.Metrics.Path,
	}

	middlewares := make([]echo.MiddlewareFunc, 10)
	middlewares = append(middlewares,
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			// Request logger values for middleware.RequestLoggerValues
			LogError:  true,
			LogMethod: true,
			LogStatus: true,
			LogURI:    true,

			// Forwards error to the global error handler, so it can decide
			// appropriate status code.
			HandleError: true,

			Skipper: loggerSkipperWithPaths(skipperPaths...),

			LogValuesFunc: logger.MiddlewareLogValues,
		}),
	)
	middlewares = append(middlewares, middleware.Recover())

	e.Use(middlewares...)
}

// NewRouterWithConfig fill the router configuration for the given echo instance,
// providing routes for the public endpoints, the private paths (includes the healthcheck),
// and the /metrics path
// e is the echo instance where to add the routes.
// c is the router configuration.
// metrics is the reference to the metrics storage.
// Return the echo instance set up; is something fails it panics.
func NewRouterWithConfig(e *echo.Echo, cfg *config.Config, public *openapi3.T) *echo.Echo {
	if e == nil {
		panic("'e' is nil")
	}
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if public == nil {
		panic("'public' is nil")
	}

	configCommonMiddlewares(e, cfg)

	newGroupPrivate(e.Group(privatePath), cfg)
	newGroupPublic(e.Group(basepublicPath+"/v"+public.Info.Version), cfg)
	newGroupPublic(e.Group(basepublicPath+"/v"+getMajorVersion(c.Version)), c)
	return e
}

// NewRouterForMetrics fill the routing information for /metrics endpoint.
// e is the echo instance
// c is the router configuration
// Return the echo instance configured for the metrics for success execution,
// else raise any panic.
func NewRouterForMetrics(e *echo.Echo, cfg *config.Config, handlers metrics.ServiceInterface) *echo.Echo {
	if e == nil {
		panic("'e' is nil")
	}
	if c.MetricsPath == "" {
		panic(fmt.Errorf("MetricsPath cannot be an empty string"))
	}

	configCommonMiddlewares(e, cfg)

	// Register handlers
	return newGroupMetrics(e, cfg)
}
