package echo

import (
	"log/slog"
	"net/http"
	"testing"

	helper_http_echo "github.com/avisiedo/go-microservice-1/internal/test/helper/http/echo"
	"github.com/avisiedo/go-microservice-1/internal/test/mock/interface/interactor"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHealthcheck(t *testing.T) {
	assert.PanicsWithValue(t, "the interactor is nil", func() {
		NewHealthcheck(nil)
	})

	i := interactor.NewHealthcheckInteractor(t)
	p := NewHealthcheck(i)
	assert.NotNil(t, p)
	i.AssertExpectations(t)
}

func TestGetLivez(t *testing.T) {
	const path = "/livez"
	const interactorMethod = "IsLive"
	var ctx echo.Context
	i := interactor.NewHealthcheckInteractor(t)
	p := NewHealthcheck(i)
	require.NotNil(t, p)

	e := echo.New()
	err := p.GetLivez(nil)
	require.EqualError(t, err, echo.ErrInternalServerError.Error())

	ctx = helper_http_echo.NewContext(e, http.MethodGet, path, http.Header{}, nil, slog.Default())
	require.NotNil(t, ctx)
	i.On(interactorMethod).Return(nil)
	err = p.GetLivez(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, ctx.Response().Status)
	i.AssertExpectations(t)

	ctx = helper_http_echo.NewContext(e, http.MethodGet, path, http.Header{}, nil, slog.Default())
	i = interactor.NewHealthcheckInteractor(t)
	require.NotNil(t, ctx)
	i.On(interactorMethod).Return(echo.ErrServiceUnavailable)
	p = NewHealthcheck(i)
	err = p.GetLivez(ctx)
	require.EqualError(t, err, "code=503, message=Service Unavailable")
	require.Equal(t, http.StatusServiceUnavailable, ctx.Response().Status)
	i.AssertExpectations(t)
}

func TestReadyz(t *testing.T) {
	const path = "/readyz"
	const interactorMethod = "IsReady"
	i := interactor.NewHealthcheckInteractor(t)
	p := NewHealthcheck(i)
	require.NotNil(t, p)

	e := echo.New()
	err := p.GetReadyz(nil)
	require.EqualError(t, err, echo.ErrInternalServerError.Error())

	ctx := helper_http_echo.NewContext(e, http.MethodGet, path, http.Header{}, nil, slog.Default())
	require.NotNil(t, ctx)
	i = interactor.NewHealthcheckInteractor(t)
	i.On(interactorMethod).Return(nil)
	p = NewHealthcheck(i)
	err = p.GetReadyz(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, ctx.Response().Status)

	ctx = helper_http_echo.NewContext(e, http.MethodGet, path, http.Header{}, nil, slog.Default())
	i = interactor.NewHealthcheckInteractor(t)
	require.NotNil(t, ctx)
	i.On(interactorMethod).Return(echo.ErrServiceUnavailable)
	p = NewHealthcheck(i)
	err = p.GetReadyz(ctx)
	require.EqualError(t, err, "code=503, message=Service Unavailable")
	require.Equal(t, http.StatusServiceUnavailable, ctx.Response().Status)
	i.AssertExpectations(t)
}
