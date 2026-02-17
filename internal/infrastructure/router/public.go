package router

import (
	"fmt"
	"strings"

	"github.com/avisiedo/go-microservice-1/internal/api/header"
	"github.com/avisiedo/go-microservice-1/internal/api/http/openapi"
	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/config"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

const (
	basePath = "/api/todo"
)

func createMiddlewares(cfg *config.Config, metrics *metrics.Metrics) []echo.MiddlewareFunc {
	middlewares := []echo.MiddlewareFunc{}

	// Initialize middlewares
	middlewares = append(middlewares,
		middleware.MetricsMiddlewareWithConfig(
			&middleware.MetricsConfig{
				Metrics: metrics,
			},
		),
	)
	middlewares = append(middlewares,
		echo_middleware.RequestIDWithConfig(
			echo_middleware.RequestIDConfig{
				TargetHeader: header.HdrRequestID,
			},
		),
	)
	if cfg.Application.ValidateAPI {
		middleware.InitOpenAPIFormats()
		middlewares = append(middlewares,
			middleware.RequestResponseValidatorWithConfig(
				// FIXME Get the values from the application config
				&middleware.RequestResponseValidatorConfig{
					Skipper:          nil,
					ValidateRequest:  true,
					ValidateResponse: false,
				},
			),
		)
	}
	return middlewares
}

func newPublic(e *echo.Group, cfg *config.Config, publicHandler public.ServerInterface, openapiHandler openapi.ServerInterface, metricsHandler *metrics.Metrics) *echo.Group {
	if e == nil {
		panic(common_err.ErrNil("e"))
	}
	if cfg == nil {
		panic(common_err.ErrNil("cfg"))
	}
	if publicHandler == nil {
		panic(common_err.ErrNil("publicHandler"))
	}
	if openapiHandler == nil {
		panic(common_err.ErrNil("openapiHandler"))
	}
	if metricsHandler == nil {
		panic(common_err.ErrNil("metricsHandler"))
	}

	// Wire the middlewares
	middlewares := createMiddlewares(cfg, metricsHandler)
	e.Use(middlewares...)

	// Setup routes
	public.RegisterHandlers(e, publicHandler)
	openapi.RegisterHandlers(e, openapiHandler)
	return e
}

func getOpenapiPaths(cfg *config.Config, swagger *openapi3.T) func() []string {
	if cfg == nil {
		panic(common_err.ErrNil("cfg"))
	}
	if swagger == nil {
		panic(common_err.ErrNil("swagger"))
	}
	version := swagger.Info.Version
	if version == "" {
		panic(common_err.ErrEmpty("swagger.Info.Version"))
	}
	majorVersion := strings.Split(version, ".")[0]
	majorMinorVersion := fmt.Sprintf("%s.%s", majorVersion, strings.Split(version, ".")[1])
	cachedPaths := []string{
		fmt.Sprintf("%s/v%s/openapi.json", basePath, majorVersion),
		fmt.Sprintf("%s/v%s/openapi.json", basePath, majorMinorVersion),
	}
	return func() []string {
		return cachedPaths
	}
}

// newSkipperOpenapi skip /api/todo/v*/openapi.json path
func newSkipperOpenapi(cfg *config.Config) echo_middleware.Skipper {
	var (
		swagger *openapi3.T
		err     error
	)
	if swagger, err = public.GetSwagger(); err != nil {
		panic(err)
	}
	paths := getOpenapiPaths(cfg, swagger)()
	return func(ctx echo.Context) bool {
		route := ctx.Path()
		for i := range paths {
			if paths[i] == route {
				return true
			}
		}
		return false
	}
}
