package logger

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"testing"

	test_echo "github.com/avisiedo/go-microservice-1/internal/test/helper/http/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareLogValues(t *testing.T) {
	slogDefault := slog.Default()
	defer slog.SetDefault(slogDefault)

	b := bytes.NewBufferString("")
	slog.SetDefault(slog.New(slog.NewTextHandler(b, &slog.HandlerOptions{})))
	e := echo.New()
	ctx := test_echo.NewContext(e, http.MethodGet, "/api/todos/v1/todo", nil, nil, slogDefault)
	err := MiddlewareLogValues(ctx, middleware.RequestLoggerValues{
		Method:    http.MethodGet,
		URIPath:   "/api/todos/v1/todo",
		RoutePath: "/api/todos/v1/todo",
		Status:    http.StatusOK,
	})
	require.NoError(t, err)
	// http://101regexp
	regExp := "^" + timeExp + " " + levelExp + " " + `msg="(.*)" request-id="(.*)" method=(.*) uri="(.*)" status=(.*)\n$`
	require.Regexp(t, regExp, b.String())
}

func TestMiddlewareLogValuesError(t *testing.T) {
	slogDefault := slog.Default()
	defer slog.SetDefault(slogDefault)

	b := bytes.NewBufferString("")
	slog.SetDefault(slog.New(slog.NewTextHandler(b, &slog.HandlerOptions{})))
	e := echo.New()
	ctx := test_echo.NewContext(e, http.MethodGet, "/api/todos/v1/todo", nil, nil, slogDefault)
	err := MiddlewareLogValues(ctx, middleware.RequestLoggerValues{
		Method:    http.MethodGet,
		URIPath:   "/api/todos/v1/todo",
		RoutePath: "/api/todos/v1/todo",
		Status:    http.StatusOK,
		Error:     errors.New("an error happened"),
	})
	require.NoError(t, err)
	// http://101regexp
	regExp := "^" + timeExp + " " + levelExp + " " + `msg="(.*)" request-id="(.*)" method=(.*) uri="(.*)" status=(.*) err="(.*)"\n$`
	require.Regexp(t, regExp, b.String())
}
