package echo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

type AdditionalContext interface {
	SetLog(log *slog.Logger)
	Log() *slog.Logger
}

type ExtendedContext interface {
	echo.Context
	AdditionalContext
}

type extendedContext struct {
	logger *slog.Logger
	echo.Context
	AdditionalContext
}

func NewExtendedContext(ctx echo.Context) ExtendedContext {
	extContext := &extendedContext{
		Context: ctx,
	}
	extContext.AdditionalContext = extContext
	return extContext
}

func (ctx *extendedContext) SetLog(log *slog.Logger) {
	ctx.logger = log
}

func checkHttpMethod(value string) {
	switch value {
	case http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace:
		return
	default:
		panic(fmt.Sprintf("method '%s' is not valid", value))
	}
}

// NewContextWithContext create an echo.Context related with go context.Context.
func NewContextWithContext(ctx context.Context, e *echo.Echo, method, path string, headers http.Header, body any) ExtendedContext {
	if e == nil {
		panic("echo instance is nil")
	}
	checkHttpMethod(method)
	var bodyReader io.Reader = nil
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}
	req, err := http.NewRequestWithContext(ctx, method, path, bodyReader)
	if err != nil {
		panic(err)
	}

	if headers == nil && body != nil {
		headers = http.Header{}
	}
	req.Header = headers
	if body != nil {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	res := httptest.NewRecorder()

	echoCtx := NewExtendedContext(e.NewContext(req, res))
	echoCtx.SetLog(slog.Default())
	return echoCtx
}

// NewContext create a new echo context ready for a test request.
func NewContext(e *echo.Echo, method, path string, headers http.Header, body any, logger *slog.Logger) ExtendedContext {
	return NewContextWithContext(context.Background(), e, method, path, headers, body)
}

// NewHandler create a demo handler
func NewHandler(status int, response any, err error) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err != nil {
			return err
		}
		return c.JSON(status, response)
	}
}

func NewDummyContext(e *echo.Echo) echo.Context {
	return NewContextWithContext(context.Background(), e, http.MethodGet, "/", http.Header{}, nil)
}
