package middleware

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	app_context "github.com/avisiedo/go-microservice-1/internal/infrastructure/context"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestSLogMiddlewareWithConfig(t *testing.T) {
	require.PanicsWithError(t, common_err.ErrNil("cfg").Error(), func() {
		_ = SLogMiddlewareWithConfig(nil)
	})

	cfg := &SLogMiddlewareConfig{
		Log: nil,
	}
	m := SLogMiddlewareWithConfig(cfg)
	require.NotNil(t, m)

	cfg.Log = slog.Default()
	m = SLogMiddlewareWithConfig(cfg)
	require.NotNil(t, m)

	cfg.Skipper = func(c echo.Context) bool {
		return true
	}
	m = SLogMiddlewareWithConfig(cfg)
	require.NotNil(t, m)
	e := helperNewEchoNooperation(http.MethodGet, "/", m)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	e.ServeHTTP(res, req)
	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "Ok", res.Body.String())

	cfg.Skipper = nil
	e = helperNewEchoNooperation(http.MethodGet, "/", m)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	e.ServeHTTP(res, req)
	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "Ok", res.Body.String())
}

func TestSLogRequest(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level:     logger.LevelInfo,
		AddSource: false,
	}))
	ctx := app_context.WithLog(context.TODO(), l)
	SLogRequest(ctx, nil, http.MethodGet, "/test", http.StatusOK)
	timeRegexp := `time=[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}T[[:digit:]]{1,2}:[[:digit:]]{2}:[[:digit:]]{2}\.([[:digit:]]*)(\+[[:digit:]]{2}:[[:digit:]]{2}|Z|)`
	regexp := timeRegexp + " level=INFO " + "msg=\"success request\" method=GET path=/test status=200"
	buffString := buf.String()
	require.Regexp(t, regexp, buffString)

	buf = bytes.NewBuffer([]byte{})
	l = slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level:     logger.LevelInfo,
		AddSource: false,
	}))
	ctx = app_context.WithLog(context.TODO(), l)
	SLogRequest(ctx, errors.New("test error"), http.MethodGet, "/test", http.StatusInternalServerError)
	regexp = timeRegexp + " level=ERROR msg=\"test error\" method=GET path=/test status=500"
	buffString = buf.String()
	require.Regexp(t, regexp, buffString)
}
