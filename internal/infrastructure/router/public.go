package router

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/api/http/openapi"
	"github.com/podengo-project/idmsvc-backend/internal/api/http/public"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
)

const (
	headerRequestID = "X-Request-Id"
	basePath        = "/api/todo"
)

func getOpenapiPaths(cfg *config.Config) func() []string {
	swagger, err := public.GetSwagger()
	if err != nil {
		panic("error calling public.GetSwagger()")
	}
	version := swagger.Info.Version
	if version == "" {
		panic(fmt.Errorf("'Info.Version' at public api is empty"))
	}
	majorVersion := strings.Split(version, ".")[0]
	fullVersion := version
	cachedPaths := []string{
		fmt.Sprintf("%s/v%s/openapi.json", basePath, fullVersion),
		fmt.Sprintf("%s/v%s/openapi.json", basePath, majorVersion),
	}
	return func() []string {
		return cachedPaths
	}
}

func newGroupPublic(e *echo.Group, cfg *config.Config, publicHandler public.ServerInterface, openapiHanlder openapi.ServerInterface, metrics *metrics.Metrics) *echo.Group {
	if e == nil {
		panic("echo group is nil")
	}

	middlewares := make([]echo.MiddlewareFunc, 10)

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
	openapi.RegisterHandlersWithBaseURL(e, openapiHanlder, "")
	return e
}

// skipperOpenapi skip /api/idmsvc/v*/openapi.json path
func newSkipperOpenapi(cfg *config.Config) echo_middleware.Skipper {
	paths := getOpenapiPaths(cfg)()
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
