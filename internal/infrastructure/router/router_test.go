package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/config"
	handler_impl "github.com/avisiedo/go-microservice-1/internal/handler/http/impl"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	"github.com/avisiedo/go-microservice-1/internal/test"
	presenter "github.com/avisiedo/go-microservice-1/internal/test/mock/api/http/metrics"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGetMajorVersion(t *testing.T) {
	assert.Equal(t, "", getMajorVersion(""))
	assert.Equal(t, "1", getMajorVersion("1.0"))
	assert.Equal(t, "1", getMajorVersion("1.0.3"))
	assert.Equal(t, "1", getMajorVersion("1."))
	assert.Equal(t, "a", getMajorVersion("a.b.c"))
}

func TestLoggerSkipperWithPaths(t *testing.T) {
	var skipper middleware.Skipper

	// Empty path does not panic
	assert.NotPanics(t, func() {
		skipper = loggerSkipperWithPaths()
	})
	assert.NotNil(t, skipper)

	// Only one path does not panic
	assert.NotPanics(t, func() {
		skipper = loggerSkipperWithPaths("/test")
	})
	assert.NotNil(t, skipper)

	// Check several paths
	assert.NotPanics(t, func() {
		skipper = loggerSkipperWithPaths("/test", "/anothertest")
	})
	assert.NotNil(t, skipper)

	// Check skipped paths
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/test")
	assert.True(t, skipper(ctx))

	req = httptest.NewRequest(http.MethodGet, "/anothertest", nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetPath("/anothertest")
	assert.True(t, skipper(ctx))

	// Check no skipped paths
	req = httptest.NewRequest(http.MethodGet, "/noskipped", nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetPath("/noskipped")
	assert.False(t, skipper(ctx))
}

func TestConfigCommonMiddlewares(t *testing.T) {
	cfg := &config.Config{}
	_ = config.Load(cfg)
	e := echo.New()
	configCommonMiddlewares(e, cfg)
}

func TestNewRouterWithConfig(t *testing.T) {
	assert.Panics(t, func() {
		NewRouterWithConfig(nil, nil, nil, nil, nil)
	}, "'e' is nil")

	e := echo.New()
	assert.Panics(t, func() {
		NewRouterWithConfig(e, nil, nil, nil, nil)
	})

	cfg := test.GetTestConfig()
	reg := prometheus.NewRegistry()
	m := metrics.NewMetrics(reg)
	_, db, _ := test.NewSqlMock(&gorm.Session{SkipHooks: true})

	// Create application handlers
	app := handler_impl.NewHandler(cfg, db, m)
	swagger, err := public.GetSwagger()
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		_ = NewRouterWithConfig(e, cfg, swagger, app, m)
	})
}

func TestNewRouterForMetrics(t *testing.T) {
	presenterMetrics := presenter.NewServerInterface(t)
	assert.Panics(t, func() {
		NewRouterForMetrics(nil, nil, nil)
	})

	e := echo.New()
	assert.Panics(t, func() {
		NewRouterForMetrics(e, nil, nil)
	})

	cfg := &config.Config{}
	_ = config.Load(cfg)

	assert.Panics(t, func() {
		_ = NewRouterForMetrics(e, cfg, nil)
	}, "'e' is nil")

	assert.NotPanics(t, func() {
		e = NewRouterForMetrics(e, cfg, presenterMetrics)
	}, "MetricsPath cannot be an empty string")
}
