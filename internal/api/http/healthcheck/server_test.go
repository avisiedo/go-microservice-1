package healthcheck

import (
	"net/http"
	"testing"

	mock_healthcheck "github.com/avisiedo/go-microservice-1/internal/test/mock/api/http/healthcheck"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetLivez(t *testing.T) {
	const methodName = "GetLivez"
	var (
		m *mock_healthcheck.MockServerInterface
		h ServerInterfaceWrapper
	)

	m = mock_healthcheck.NewMockServerInterface(t)
	h = ServerInterfaceWrapper{Handler: m}
	m.On(methodName, nil).Return(nil)
	r := h.GetLivez(nil)
	require.NoError(t, r)
	mock.AssertExpectationsForObjects(t, m)

	m = mock_healthcheck.NewMockServerInterface(t)
	h = ServerInterfaceWrapper{Handler: m}
	m.On(methodName, nil).Return(echo.NewHTTPError(http.StatusInternalServerError, "internal server error"))
	r = h.GetLivez(nil)
	require.EqualError(t, r, "code=500, message=internal server error")
	mock.AssertExpectationsForObjects(t, m)
}

func TestGetReadyz(t *testing.T) {
	const methodName = "GetReadyz"
	var (
		m *mock_healthcheck.MockServerInterface
		h ServerInterfaceWrapper
	)

	m = mock_healthcheck.NewMockServerInterface(t)
	h = ServerInterfaceWrapper{Handler: m}
	m.On(methodName, nil).Return(nil)
	r := h.GetReadyz(nil)
	require.NoError(t, r)
	mock.AssertExpectationsForObjects(t, m)

	m = mock_healthcheck.NewMockServerInterface(t)
	h = ServerInterfaceWrapper{Handler: m}
	m.On(methodName, nil).Return(echo.NewHTTPError(http.StatusInternalServerError, "internal server error"))
	r = h.GetReadyz(nil)
	require.EqualError(t, r, "code=500, message=internal server error")
	mock.AssertExpectationsForObjects(t, m)
}

func TestRegisterHandlersWithBaseURL(t *testing.T) {
	e := mock_healthcheck.NewMockEchoRouter(t)
	w := mock_healthcheck.NewMockServerInterface(t)
	e.On("GET", "/root/livez", mock.AnythingOfType("echo.HandlerFunc")).Return(nil)
	e.On("GET", "/root/readyz", mock.AnythingOfType("echo.HandlerFunc")).Return(nil)
	require.NotPanics(t, func() {
		RegisterHandlersWithBaseURL(e, w, "/root")
	})
	w.AssertExpectations(t)
}

func TestRegisterHandlers(t *testing.T) {
	e := mock_healthcheck.NewMockEchoRouter(t)
	w := mock_healthcheck.NewMockServerInterface(t)
	e.On("GET", "/livez", mock.AnythingOfType("echo.HandlerFunc")).Return(nil)
	e.On("GET", "/readyz", mock.AnythingOfType("echo.HandlerFunc")).Return(nil)
	require.NotPanics(t, func() {
		RegisterHandlers(e, w)
	})
	w.AssertExpectations(t)
}
