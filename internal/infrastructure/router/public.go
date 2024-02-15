package router

import (
	"fmt"
	"strings"

	"github.com/avisiedo/go-microservice-1/internal/api/http/openapi"
	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

const (
	headerRequestID = "X-Request-Id"
	basePath        = "/api/todo"
)

func newGroupPublic(e *echo.Group, cfg *config.Config, publicHandler public.ServerInterface, openapiHandler openapi.ServerInterface, metrics *metrics.Metrics) *echo.Group {
	if e == nil {
		panic("echo group is nil")
	}

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
				TargetHeader: headerRequestID,
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

	// Wire the middlewares
	e.Use(middlewares...)

	// Setup routes
	public.RegisterHandlersWithBaseURL(e, publicHandler, "")
	openapi.RegisterHandlersWithBaseURL(e, openapiHandler, "")
	return e
}

func getOpenapiPaths(cfg *config.Config, swagger *openapi3.T) func() []string {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if swagger == nil {
		panic("'swagger' is nil")
	}
	version := swagger.Info.Version
	if version == "" {
		panic(fmt.Errorf("'Info.Version' at public api is empty"))
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

// skipperOpenapi skip /api/idmsvc/v*/openapi.json path
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
