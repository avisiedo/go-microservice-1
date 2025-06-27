package router

import (
	"log/slog"
	"strings"

	"github.com/avisiedo/go-microservice-1/internal/api/http/healthcheck"
	metrics_handler "github.com/avisiedo/go-microservice-1/internal/api/http/metrics"
	"github.com/avisiedo/go-microservice-1/internal/config"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
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
	middlewares = append(middlewares,
		app_middleware.SLogMiddlewareWithConfig(&app_middleware.SLogMiddlewareConfig{
			Skipper: loggerSkipperWithPaths(skipperPaths...),
			Log:     slog.Default(),
		}),
	)

	middlewares = append(middlewares, middleware.Recover())

	e.Use(middlewares...)
}

func newRouterWithConfigCommonGuards(e *echo.Echo, cfg *config.Config) {
	if e == nil {
		panic(common_err.ErrNil("e"))
	}
	if cfg == nil {
		panic(common_err.ErrNil("cfg"))
	}
}

func newRouterWithConfigGuards(e *echo.Echo, cfg *config.Config, public *openapi3.T) {
	newRouterWithConfigCommonGuards(e, cfg)
	if public == nil {
		panic(common_err.ErrNil("public"))
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

	healthcheck.RegisterHandlers(e, h)
	newPrivate(e.Group(privatePath), cfg, h)
	newPublic(e.Group(cfg.Application.PathPrefix), cfg, h, h, m)
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
