package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	"github.com/labstack/echo/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const URLPrefix = "/api/" + config.DefaultAppName

func TestCreateMetricsMiddleware(t *testing.T) {
	var (
		m          *metrics.Metrics
		middleware echo.MiddlewareFunc
	)
	m = metrics.NewMetrics(prometheus.NewRegistry())
	middleware = CreateMetricsMiddleware(m)

	assert.NotNil(t, middleware)
}

func TestMetricsMiddlewareWithConfigCreation(t *testing.T) {
	var (
		reg    *prometheus.Registry
		config *MetricsConfig
	)

	config = &MetricsConfig{
		Metrics: nil,
		Skipper: nil,
	}
	assert.Panics(t, func() {
		MetricsMiddlewareWithConfig(config)
	})

	reg = prometheus.NewRegistry()
	config = &MetricsConfig{
		Metrics: metrics.NewMetrics(reg),
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/ping"
		},
	}

	require.NotPanics(t, func() {
		MetricsMiddlewareWithConfig(config)
	})

	assert.NotPanics(t, func() {
		MetricsMiddlewareWithConfig(nil)
	})

	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	}

	e := echo.New()
	m := MetricsMiddlewareWithConfig(config)
	e.Use(m)
	path := "/api/todo/v1/todos/"
	e.Add(http.MethodGet, path, h)

	// Check normal request
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	e.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "Ok", resp.Body.String())

	// Check skipper
	resp = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/ping", nil)
	e.ServeHTTP(resp, req)
}
