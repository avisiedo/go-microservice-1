package router

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slog"

	metrics_handler "github.com/avisiedo/go-microservice-1/internal/api/http/metrics"
	"github.com/avisiedo/go-microservice-1/internal/config"
	handler "github.com/avisiedo/go-microservice-1/internal/handler/http"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	app_middleware "github.com/avisiedo/go-microservice-1/internal/infrastructure/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	privatePath = "/private"
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

	middlewares := []echo.MiddlewareFunc{}
	middlewares = append(middlewares,
		app_middleware.SLogMiddlewareWithConfig(&app_middleware.SLogMiddlewareConfig{
			Skipper: loggerSkipperWithPaths(skipperPaths...),
			Log:     slog.Default(),
		}),
	)
	// middlewares = append(middlewares,
	// 	middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	// 		// Request logger values for middleware.RequestLoggerValues
	// 		LogError:  true,
	// 		LogMethod: true,
	// 		LogStatus: true,
	// 		LogURI:    true,

	// 		// Forwards error to the global error handler, so it can decide
	// 		// appropriate status code.
	// 		HandleError: true,

	// 		Skipper: loggerSkipperWithPaths(skipperPaths...),

	// 		LogValuesFunc: logger.MiddlewareLogValues,
	// 	}),
	// )
	middlewares = append(middlewares, middleware.Recover())

	e.Use(middlewares...)
}

func newRouterWithConfigGuards(e *echo.Echo, cfg *config.Config, public *openapi3.T) {
	if e == nil {
		panic("'e' echo server is nil")
	}
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if public == nil {
		panic("'public' is nil")
	}
}

// NewRouterWithConfig fill the router configuration for the given echo instance,
// providing routes for the public endpoints, the private paths (includes the healthcheck),
// and the /metrics path
// e is the echo instance where to add the routes.
// cfg is the application configuration.
// public is the openapi specification.
// h is the application handler.
// m is the reference to the metrics storage.
// Return the echo instance set up; is something fails it panics.
func NewRouterWithConfig(e *echo.Echo, cfg *config.Config, public *openapi3.T, h handler.Application, m *metrics.Metrics) *echo.Echo {
	newRouterWithConfigGuards(e, cfg, public)
	configCommonMiddlewares(e, cfg)

	majorVersion := strings.Split(public.Info.Version, ".")[0]
	majorMinorVersion := strings.Join([]string{majorVersion, strings.Split(public.Info.Version, ".")[1]}, ".")

	newGroupPrivate(e.Group(privatePath), cfg)
	newGroupPublic(e.Group(cfg.Application.PathPrefix+"/v"+majorVersion), cfg, h, h, m)
	newGroupPublic(e.Group(cfg.Application.PathPrefix+"/v"+majorMinorVersion), cfg, h, h, m)
	return e
}

// NewRouterForMetrics fill the routing information for /metrics endpoint.
// e is the echo instance
// cfg is the router configuration
// h is the handler to retrieve the metrics.
// Return the echo instance configured for the metrics for success execution,
// else raise any panic.
func NewRouterForMetrics(e *echo.Echo, cfg *config.Config, h metrics_handler.ServerInterface) *echo.Echo {
	if e == nil {
		panic("'e' is nil")
	}
	if cfg.Metrics.Path == "" {
		panic(fmt.Errorf("'cfg.Metrics.Path' is an empty string"))
	}
	if h == nil {
		panic(fmt.Errorf("'h' is nil"))
	}

	configCommonMiddlewares(e, cfg)

	// Register handlers
	return newGroupMetrics(e, cfg, h)
}
