package echo

import (
	"net/http"
	"testing"

	infra_metrics "github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	helper_http_echo "github.com/avisiedo/go-microservice-1/internal/test/helper/http/echo"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstrumentation(t *testing.T) {
	assert.PanicsWithValue(t, "'metrics' is nil", func() {
		_ = NewInstrumentation(nil)
	})

	m := &infra_metrics.Metrics{}
	p := NewInstrumentation(m)
	require.NotNil(t, p)
}

func TestGetMetrics(t *testing.T) {
	const path = "/metrics"

	reg := prometheus.NewRegistry()
	m := infra_metrics.NewMetrics(reg)
	p := NewInstrumentation(m)
	e := echo.New()
	ctx := helper_http_echo.NewContext(e, http.MethodGet, path, http.Header{}, nil)
	err := p.GetMetrics(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, ctx.Response().Status)
}
