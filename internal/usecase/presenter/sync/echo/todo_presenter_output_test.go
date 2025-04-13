package echo

import (
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	builder "github.com/avisiedo/go-microservice-1/internal/test/builder/model"
	helper_echo "github.com/avisiedo/go-microservice-1/internal/test/helper/http/echo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTodoOutput(t *testing.T) {
	sut := NewTodoOutput()
	require.NotNil(t, sut)
}

func TestOutputCreateGuards(t *testing.T) {
	var err error
	sut := &todoOutput{}

	e := echo.New()
	require.NotNil(t, e)

	err = sut.createGuards(nil, nil)
	require.EqualError(t, err, "ctx is nil")

	err = sut.createGuards(helper_echo.NewDummyContext(e), nil)
	require.EqualError(t, err, "dataOutput is nil")

	err = sut.createGuards(helper_echo.NewDummyContext(e), &model.Todo{})
	require.NoError(t, err)
}

func TestOutputCreate(t *testing.T) {
	var (
		err  error
		data *public.ToDo
	)

	sut := &todoOutput{}

	e := echo.New()

	data, err = sut.Create(helper_echo.NewDummyContext(e), nil)
	require.EqualError(t, err, "dataOutput is nil")
	assert.Nil(t, data)

	data, err = sut.Create(
		helper_echo.NewDummyContext(e),
		builder.NewTodo().
			WithTitle("todo title test").
			WithDescription("todo description test").
			Build())
	require.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "todo title test", data.Title)
	assert.Equal(t, "todo description test", data.Description)
}

func TestOutputGetAllGuards(t *testing.T) {
	var err error

	sut := &todoOutput{}

	e := echo.New()
	require.NotNil(t, e)

	err = sut.getAllGuards(nil, nil)
	require.EqualError(t, err, "ctx is nil")

	err = sut.getAllGuards(helper_echo.NewDummyContext(e), nil)
	require.EqualError(t, err, "dataOutput is nil")

	err = sut.getAllGuards(helper_echo.NewDummyContext(e), []model.Todo{*builder.NewTodo().Build()})
	require.NoError(t, err)
}

func TestOutputGetAll(t *testing.T) {
	var (
		err    error
		output []public.ToDo
	)

	sut := &todoOutput{}

	e := echo.New()
	require.NotNil(t, e)

	output, err = sut.GetAll(nil, nil)
	require.EqualError(t, err, "ctx is nil")
	assert.Nil(t, output)

	input := []model.Todo{*builder.NewTodo().Build()}
	output, err = sut.GetAll(helper_echo.NewDummyContext(e), input)
	require.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, len(input), len(output))
	assert.Equal(t, input[0].Title, output[0].Title)
	assert.Equal(t, input[0].Description, output[0].Description)
	assert.Equal(t, input[0].DueDate, output[0].DueDate)
}

func TestOutputGetGuards(t *testing.T) {
	var err error

	sut := &todoOutput{}

	e := echo.New()
	require.NotNil(t, e)

	err = sut.getGuards(nil, nil)
	require.EqualError(t, err, "'ctx' is nil")

	err = sut.getGuards(helper_echo.NewDummyContext(e), nil)
	require.EqualError(t, err, "'data' is nil")

	err = sut.getGuards(helper_echo.NewDummyContext(e), builder.NewTodo().WithID(uuid.UUID{}).Build())
	require.EqualError(t, err, "'data.UUID' is empty")

	err = sut.getGuards(helper_echo.NewDummyContext(e), builder.NewTodo().Build())
	require.NoError(t, err)
}

func TestOutputGet(t *testing.T) {
	var (
		err    error
		output *public.ToDo
	)

	sut := &todoOutput{}

	e := echo.New()
	require.NotNil(t, e)

	output, err = sut.Get(nil, nil)
	require.EqualError(t, err, "'ctx' is nil")
	assert.Nil(t, output)

	input := builder.NewTodo().Build()
	output, err = sut.Get(helper_echo.NewDummyContext(e), input)
	require.NoError(t, err)
	require.NotNil(t, output.TodoId)
	assert.Equal(t, input.UUID, *output.TodoId)
	assert.Equal(t, input.Title, output.Title)
	assert.Equal(t, input.Description, output.Description)
	assert.Equal(t, input.DueDate, output.DueDate)
	assert.Equal(t, input.Title, output.Title)
}
