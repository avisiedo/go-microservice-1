package logger

import (
	"log/slog"

	"github.com/avisiedo/go-microservice-1/internal/api/header"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	request_id := req.Header.Get(header.HdrRequestID)
	if request_id == "" {
		request_id = res.Header().Get(header.HdrRequestID)
	}

	logAttr = append(logAttr,
		slog.String("request_id", request_id),
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
