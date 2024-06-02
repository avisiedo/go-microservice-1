package middleware

import (
	"context"
	"log/slog"

	app_context "github.com/avisiedo/go-microservice-1/internal/infrastructure/context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SLogMiddlewareConfig hold the configuration for the echo slog middleware
type SLogMiddlewareConfig struct {
	Skipper middleware.Skipper
	Log     *slog.Logger
}

// SLogMiddlewareWithConfig create a middleware for logging
// config provides
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
			ctx := c.Request().Context()
			ctx = app_context.WithLog(ctx, config.Log)
			c.SetRequest(req.WithContext(ctx))

			// Invoke next middleware
			err := next(c)

			// Log status
			SLogRequest(
				c.Request().Context(),
				err,
				c.Request().Method,
				c.Request().URL.RequestURI(),
				c.Response().Status,
			)
			return err
		}
	}
}

// SLogRequest write a log entry with the request information
func SLogRequest(ctx context.Context, err error, method, path string, status int) {
	l := app_context.LogFromContext(ctx)
	if err != nil {
		l.ErrorContext(
			ctx,
			err.Error(),
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
		)
	}
	l.InfoContext(
		ctx,
		"success request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status", status),
	)
}
