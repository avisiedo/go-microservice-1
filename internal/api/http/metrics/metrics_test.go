package metrics

import (
	"log/slog"
	"net/http"
	"testing"

	echo_helper "github.com/avisiedo/go-microservice-1/internal/test/helper/http/echo"
	http_metrics "github.com/avisiedo/go-microservice-1/internal/test/mock/api/http/metrics"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetMetricsWithNoError(t *testing.T) {
	handlerMock := http_metrics.NewServerInterface(t)
	wrapper := &ServerInterfaceWrapper{
		Handler: handlerMock,
	}

	e := echo.New()
	ctx := echo_helper.NewContext(e, http.MethodGet, "/metrics", nil, nil, slog.Default())
	handlerMock.On("GetMetrics", ctx).Return(nil)
	assert.NoError(t, wrapper.GetMetrics(ctx))
}

func TestGetMetricsWithError(t *testing.T) {
	handlerMock := http_metrics.NewServerInterface(t)
	wrapper := &ServerInterfaceWrapper{
		Handler: handlerMock,
	}

	e := echo.New()
	ctx := echo_helper.NewContext(e, http.MethodGet, "/metrics", nil, nil, slog.Default())
	handlerMock.On("GetMetrics", ctx).Return(echo.NewHTTPError(http.StatusBadRequest, "Bad Request"))
	err := wrapper.GetMetrics(ctx)
	assert.EqualError(t, err, "code=400, message=Bad Request")
}

func TestRegisterHandlersWithBaseURL(t *testing.T) {
	e := http_metrics.NewEchoRouter(t)
	w := http_metrics.NewServerInterface(t)
	e.On("GET", "/root", mock.AnythingOfType("echo.HandlerFunc")).Return(nil)
	require.NotPanics(t, func() {
		RegisterHandlersWithBaseURL(e, w, "/root")
	})
	w.AssertExpectations(t)
}

func TestRegisterHandlers(t *testing.T) {
	e := http_metrics.NewEchoRouter(t)
	w := http_metrics.NewServerInterface(t)
	e.On("GET", "/metrics", mock.AnythingOfType("echo.HandlerFunc")).Return(nil)
	require.NotPanics(t, func() {
		RegisterHandlers(e, w)
	})
	w.AssertExpectations(t)
}
