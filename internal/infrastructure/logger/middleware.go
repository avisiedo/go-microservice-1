package logger

import (
	"log/slog"

	"github.com/avisiedo/go-microservice-1/internal/api/header"
	// "github.com/avisiedo/go-microservice-1/internal/infrastructure/logger"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// This requires the following values to be set in
// middleware.RequestLoggerWithConfig:
//
// LogError:  true,
// LogMethod: true,
// LogStatus: true,
// LogURI:    true,
func MiddlewareLogValues(c echo.Context, v middleware.RequestLoggerValues) error {
	var logLevel slog.Level
	logAttr := []slog.Attr{}

	req := c.Request()
	res := c.Response()

	requestID := req.Header.Get(header.HdrRequestID)
	if requestID == "" {
		requestID = res.Header().Get(header.HdrRequestID)
	}

	logAttr = append(logAttr,
		slog.String("request-id", requestID),
		slog.String("method", v.Method),
		slog.String("uri", v.URI),
		slog.Int("status", v.Status),
	)
	if v.Error == nil {
		logLevel = slog.LevelInfo
	} else {
		logLevel = slog.LevelError
		logAttr = append(logAttr, slog.String("err", v.Error.Error()))
	}

	slog.LogAttrs(
		req.Context(),
		logLevel,
		"log request",
		logAttr...,
	)

	return nil
}
