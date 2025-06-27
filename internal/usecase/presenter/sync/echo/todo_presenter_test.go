package echo

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	presenter "github.com/avisiedo/go-microservice-1/internal/interface/presenter/sync/echo"
	"github.com/avisiedo/go-microservice-1/internal/test"
	http_builder "github.com/avisiedo/go-microservice-1/internal/test/builder/api/http"
	model_builder "github.com/avisiedo/go-microservice-1/internal/test/builder/model"
	echo_helper "github.com/avisiedo/go-microservice-1/internal/test/helper/http/echo"
	"github.com/avisiedo/go-microservice-1/internal/test/mock/interface/interactor"
	presenter_mock "github.com/avisiedo/go-microservice-1/internal/test/mock/interface/presenter/echo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewTodo(t *testing.T) {
	assert.PanicsWithError(t, common_err.ErrNil("cfg").Error(), func() {
		NewTodo(nil, nil, nil)
	})

	cfg := config.Get()
	assert.PanicsWithError(t, common_err.ErrNil("interactor").Error(), func() {
		NewTodo(cfg, nil, nil)
	})

	i := interactor.NewTodo(t)
	assert.PanicsWithError(t, common_err.ErrNil("db").Error(), func() {
		NewTodo(cfg, i, nil)
	})

	var p presenter.Todo
	dbMock, db, err := test.NewSqlMock(nil)
	require.NotNil(t, dbMock)
	require.NotNil(t, db)
	require.NoError(t, err)
	assert.NotPanics(t, func() {
		p = NewTodo(cfg, i, db)
	})
	assert.NotNil(t, p)
}

func TestGetAllTodos(t *testing.T) {
	const path = "/todo/v1/todos"
	var (
		db      *gorm.DB
		sqlMock sqlmock.Sqlmock
		err     error
	)
	// inbound error
	cfg := config.Get()
	i := interactor.NewTodo(t)
	inputMock := presenter_mock.NewTodoInput(t)
	outputMock := presenter_mock.NewTodoOutput(t)
	if sqlMock, db, err = test.NewSqlMock(&gorm.Session{}); err != nil {
		require.NoError(t, err)
	}
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()
	expectedErr := "inbound data error"
	inputMock.On("GetAll", mock.Anything).Return(errors.New(expectedErr))
	e := echo.New()
	require.NotNil(t, e)
	p := newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)
	ctx := echo_helper.NewContext(e, http.MethodGet, path, http.Header{}, nil, slog.Default())
	require.NotNil(t, ctx)
	err = p.GetAllTodos(ctx)
	assert.EqualError(t, err, expectedErr)

	// interactor error
	cfg = config.Get()
	i = interactor.NewTodo(t)
	inputMock = presenter_mock.NewTodoInput(t)
	outputMock = presenter_mock.NewTodoOutput(t)
	if sqlMock, db, err = test.NewSqlMock(&gorm.Session{}); err != nil {
		require.NoError(t, err)
	}
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()
	expectedErr = "interactor error"
	inputMock.On("GetAll", mock.Anything).Return(nil)
	i.On("GetAll", mock.Anything).Return(nil, errors.New(expectedErr))
	p = newTodo(cfg, inputMock, outputMock, i, db)
	ctx = echo_helper.NewContext(e, http.MethodGet, path, http.Header{}, nil, slog.Default())
	require.NotNil(t, ctx)
	err = p.GetAllTodos(ctx)
	assert.EqualError(t, err, expectedErr)

	// outbound error
	cfg = config.Get()
	i = interactor.NewTodo(t)
	inputMock = presenter_mock.NewTodoInput(t)
	outputMock = presenter_mock.NewTodoOutput(t)
	if sqlMock, db, err = test.NewSqlMock(&gorm.Session{}); err != nil {
		require.NoError(t, err)
	}
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	expectedErr = "outbound error"
	inputMock.On("GetAll", mock.Anything).Return(nil)
	i.On("GetAll", mock.Anything).Return([]model.Todo{}, nil)
	outputMock.On("GetAll", mock.Anything, []model.Todo{}).Return(nil, errors.New(expectedErr))
	p = newTodo(cfg, inputMock, outputMock, i, db)
	ctx = echo_helper.NewContext(e, http.MethodGet, path, http.Header{}, nil, slog.Default())
	require.NotNil(t, ctx)
	err = p.GetAllTodos(ctx)
	assert.EqualError(t, err, expectedErr)

	// Success case
	cfg = config.Get()
	i = interactor.NewTodo(t)
	inputMock = presenter_mock.NewTodoInput(t)
	outputMock = presenter_mock.NewTodoOutput(t)
	if sqlMock, db, err = test.NewSqlMock(&gorm.Session{}); err != nil {
		require.NoError(t, err)
	}
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	inputMock.On("GetAll", mock.Anything).Return(nil)
	i.On("GetAll", mock.Anything).Return([]model.Todo{}, nil)
	outputMock.On("GetAll", mock.Anything, []model.Todo{}).Return([]public.ToDo{}, nil)
	p = newTodo(cfg, inputMock, outputMock, i, db)
	ctx = echo_helper.NewContext(e, http.MethodGet, path, http.Header{}, nil, slog.Default())
	require.NotNil(t, ctx)
	err = p.GetAllTodos(ctx)
	assert.NoError(t, err)
}

func TestCreateTodo(t *testing.T) {
	var (
		db            *gorm.DB
		sqlMock       sqlmock.Sqlmock
		err           error
		dataInput     *public.ToDo
		data          *model.Todo
		processedData *model.Todo
		dataOutput    *public.ToDo
	)

	const (
		path = "/todo/v1/todos"
	)

	// inbound error
	cfg := config.Get()
	i := interactor.NewTodo(t)
	inputMock := presenter_mock.NewTodoInput(t)
	outputMock := presenter_mock.NewTodoOutput(t)
	if _, db, err = test.NewSqlMock(&gorm.Session{}); err != nil {
		require.NoError(t, err)
	}
	expectedErr := common_err.ErrNil("ctx")
	expectedErrStr := expectedErr.Error()
	expectedHttpErr := echo.NewHTTPError(http.StatusBadRequest, expectedErrStr)
	inputMock.On("Create", mock.Anything).Return(nil, expectedErr)
	e := echo.New()
	require.NotNil(t, e)
	p := newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)
	dataInput = http_builder.NewToDo().Build()
	ctx := echo_helper.NewContext(e, http.MethodPost, path, http.Header{}, dataInput, slog.Default())
	require.NotNil(t, ctx)

	err = p.CreateTodo(ctx)
	require.Error(t, err)
	assert.EqualError(t, err, expectedHttpErr.Error())

	// interactor error
	cfg = config.Get()
	i = interactor.NewTodo(t)
	inputMock = presenter_mock.NewTodoInput(t)
	outputMock = presenter_mock.NewTodoOutput(t)
	if sqlMock, db, err = test.NewSqlMock(&gorm.Session{}); err != nil {
		require.NoError(t, err)
	}
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()
	expectedErrStr = "creating todo failed"
	expectedErr = fmt.Errorf("%s", expectedErrStr)
	expectedHttpErr = echo.NewHTTPError(http.StatusInternalServerError, expectedErrStr)
	dataInput = http_builder.NewToDo().WithTitle("title").WithDescription("description").Build()
	data = model_builder.NewTodo().WithTitle(dataInput.Title).WithDescription(dataInput.Description).Build()
	e = echo.New()
	require.NotNil(t, e)
	p = newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)
	ctx = echo_helper.NewContext(e, http.MethodPost, path, http.Header{}, dataInput, slog.Default())
	require.NotNil(t, ctx)
	inputMock.On("Create", mock.Anything).Return(data, nil)
	i.On("Create", mock.Anything, data).Return(nil, expectedErr)

	err = p.CreateTodo(ctx)
	assert.EqualError(t, err, expectedHttpErr.Error())

	// outbound error
	cfg = config.Get()
	i = interactor.NewTodo(t)
	inputMock = presenter_mock.NewTodoInput(t)
	outputMock = presenter_mock.NewTodoOutput(t)
	if sqlMock, db, err = test.NewSqlMock(&gorm.Session{}); err != nil {
		require.NoError(t, err)
	}
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	expectedErrStr = "failure on formatting output response"
	expectedErr = fmt.Errorf("%s", expectedErrStr)
	expectedHttpErr = echo.NewHTTPError(http.StatusInternalServerError, expectedErrStr)
	dataInput = http_builder.NewToDo().WithTitle("title").WithDescription("description").Build()
	data = model_builder.NewTodo().WithTitle(dataInput.Title).WithDescription(dataInput.Description).Build()
	e = echo.New()
	require.NotNil(t, e)
	p = newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)
	ctx = echo_helper.NewContext(e, http.MethodPost, path, http.Header{}, dataInput, slog.Default())
	require.NotNil(t, ctx)
	inputMock.On("Create", mock.Anything).Return(data, nil)
	i.On("Create", mock.Anything, data).Return(data, nil)
	outputMock.On("Create", mock.Anything, data).Return(nil, expectedErr)

	err = p.CreateTodo(ctx)
	require.Error(t, err)
	assert.EqualError(t, err, expectedHttpErr.Error())

	// Success scenario
	cfg = config.Get()
	i = interactor.NewTodo(t)
	inputMock = presenter_mock.NewTodoInput(t)
	outputMock = presenter_mock.NewTodoOutput(t)
	if sqlMock, db, err = test.NewSqlMock(&gorm.Session{}); err != nil {
		require.NoError(t, err)
	}
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	// expectedErrStr = "failure on formatting output response"
	// expectedErr = fmt.Errorf("%s", expectedErrStr)
	// expectedHttpErr = echo.NewHTTPError(http.StatusInternalServerError, expectedErrStr)
	dataInput = http_builder.NewToDo().
		WithTitle("title").
		WithDescription("description").
		Build()
	data = model_builder.NewTodo().
		WithTitle(dataInput.Title).
		WithDescription(dataInput.Description).
		Build()
	processedData = model_builder.NewTodo().
		WithTitle(dataInput.Title).
		WithDescription(dataInput.Description).
		Build()
	processedUUID := processedData.UUID
	dataOutput = http_builder.NewToDo().
		WithID(&processedUUID).
		WithDueDate(processedData.DueDate).
		WithTitle(processedData.Title).
		WithDescription(processedData.Description).
		WithDueDate(processedData.DueDate).
		Build()
	e = echo.New()
	require.NotNil(t, e)

	p = newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)
	ctx = echo_helper.NewContext(e, http.MethodPost, path, http.Header{}, dataInput, slog.Default())
	require.NotNil(t, ctx)
	inputMock.On("Create", mock.Anything).Return(data, nil)
	i.On("Create", mock.Anything, data).Return(processedData, nil)
	outputMock.On("Create", mock.Anything, processedData).Return(dataOutput, nil)

	err = p.CreateTodo(ctx)
	require.NoError(t, err)
}
