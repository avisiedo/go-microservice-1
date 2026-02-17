package echo

import (
	"log/slog"
	"net/http"
	"testing"

	helper_http_echo "github.com/avisiedo/go-microservice-1/internal/test/helper/http/echo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpenAPI(t *testing.T) {
	p := NewOpenAPI()
	assert.NotNil(t, p)
}

func TestGetOpenapi(t *testing.T) {
	const path = "/api/todos/v1/openapi.json"
	p := NewOpenAPI()
	require.NotNil(t, p)

	e := echo.New()
	ctx := helper_http_echo.NewContext(e, http.MethodGet, path, http.Header{}, nil, slog.Default())
	err := p.GetOpenapi(ctx)
	require.NoError(t, err)
}
