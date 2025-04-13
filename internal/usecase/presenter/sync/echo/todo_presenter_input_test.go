package echo

import (
	"net/http"
	"testing"

	. "github.com/avisiedo/go-microservice-1/internal/interface/presenter/sync/echo"
	helper_echo "github.com/avisiedo/go-microservice-1/internal/test/helper/http/echo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTodoInput(t *testing.T) {
	var r TodoInput
	assert.NotPanics(t, func() {
		r = newTodoInput()
	})
	assert.NotNil(t, r)

	assert.NotPanics(t, func() {
		r = NewTodoInput()
	})
	assert.NotNil(t, r)
}

func TestInputCreate(t *testing.T) {
	e := echo.New()

	input := newTodoInput()
	require.NotNil(t, input)

	input.Create(helper_echo.NewDummyContext(e))
}

func TestInputGetAll(t *testing.T) {
	e := echo.New()

	input := newTodoInput()
	require.NotNil(t, input)
	input.GetAll(helper_echo.NewDummyContext(e))

	input = newTodoInput()
	require.NotNil(t, input)
	input.GetAll(helper_echo.NewContext(e, http.MethodGet, "/todos/v1/todo?param=value", http.Header{}, nil))
}
