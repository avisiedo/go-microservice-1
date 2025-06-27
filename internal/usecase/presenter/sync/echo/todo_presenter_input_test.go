package echo

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	. "github.com/avisiedo/go-microservice-1/internal/interface/presenter/sync/echo"
	builder "github.com/avisiedo/go-microservice-1/internal/test/builder/model"
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

	// When Body is nil
	data, err := input.Create(helper_echo.NewContext(e, http.MethodPost, "/todos/v1/todo", http.Header{}, nil, slog.Default()))
	assert.Nil(t, data)
	require.Error(t, err)
	require.NotErrorIs(t, err, &echo.HTTPError{})
	assert.EqualError(t, err, common_err.ErrNil("Body").Error())

	// When binding fails because wrong json format
	data, err = input.Create(helper_echo.NewContext(e, http.MethodPost, "/todos/v1/todo", http.Header{}, "{", slog.Default()))
	assert.Nil(t, data)
	require.Error(t, err)
	require.NotErrorIs(t, err, &echo.HTTPError{})
	assert.EqualError(t, err, "binding request data")

	// When title is an empty string
	dataModel := builder.NewTodo().WithTitle("").Build()
	data, err = input.Create(helper_echo.NewContext(e, http.MethodPost, "/todos/v1/todo", http.Header{}, &model.Todo{}, slog.Default()))
	assert.Nil(t, data)
	require.Error(t, err)
	require.NotErrorIs(t, err, &echo.HTTPError{})
	assert.EqualError(t, err, common_err.ErrEmpty("data.Title").Error())

	// When description is an empty string
	dataModel = builder.NewTodo().WithDescription("").Build()
	data, err = input.Create(helper_echo.NewContext(e, http.MethodPost, "/todos/v1/todo", http.Header{}, dataModel, slog.Default()))
	assert.Nil(t, data)
	require.Error(t, err)
	require.NotErrorIs(t, err, &echo.HTTPError{})
	assert.EqualError(t, err, common_err.ErrEmpty("data.Description").Error())

	// Success case
	dataModel = builder.NewTodo().Build()
	data, err = input.Create(helper_echo.NewContext(e, http.MethodPost, "/todos/v1/todo", http.Header{}, dataModel, slog.Default()))
	require.NotNil(t, data)
	assert.NoError(t, err)
}

func TestInputGetAll(t *testing.T) {
	e := echo.New()

	input := newTodoInput()
	require.NotNil(t, input)
	input.GetAll(helper_echo.NewDummyContext(e))

	input = newTodoInput()
	require.NotNil(t, input)
	input.GetAll(helper_echo.NewContext(e, http.MethodGet, "/todos/v1/todo?param=value", http.Header{}, nil, slog.Default()))
}
