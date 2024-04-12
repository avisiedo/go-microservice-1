package middleware

import (
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/logger/slogctx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
)

type SLogMiddlewareConfig struct {
	Skipper middleware.Skipper
	Log     *slog.Logger
}

// SLogMiddlewareWithConfig create a middleware for logging
func SLogMiddlewareWithConfig(config *SLogMiddlewareConfig) echo.MiddlewareFunc {
	if config == nil {
		panic("'config' is nil")
	}
	if config.Log == nil {
		config.Log = slog.Default()
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}
			req := c.Request()
			c.SetRequest(req.WithContext(slogctx.NewCtx(req.Context(), config.Log)))

			// Invoke next middleware
			err := next(c)

			// Log status
			l := slogctx.FromCtx(c.Request().Context())
			if err != nil {
				l.ErrorContext(
					c.Request().Context(),
					err.Error(),
					slog.String("method", c.Request().Method),
					slog.String("path", c.Request().URL.RequestURI()),
					slog.Int("status", c.Response().Status),
				)
			} else {
				l.InfoContext(
					c.Request().Context(),
					"success request",
					slog.String("method", c.Request().Method),
					slog.String("path", c.Request().URL.RequestURI()),
					slog.Int("status", c.Response().Status),
				)
			}
			return err
		}
	}
}
