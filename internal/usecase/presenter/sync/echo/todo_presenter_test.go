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
	interactor_mock "github.com/avisiedo/go-microservice-1/internal/test/mock/interface/interactor"
	presenter_mock "github.com/avisiedo/go-microservice-1/internal/test/mock/interface/presenter/echo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func helperNewTodo(t *testing.T) (*config.Config, *presenter_mock.TodoInput, *presenter_mock.TodoOutput, *interactor_mock.Todo, sqlmock.Sqlmock, *gorm.DB) {
	var (
		sqlMock sqlmock.Sqlmock
		db      *gorm.DB
		err     error
	)
	cfg := &config.Config{}
	_ = config.Load(cfg)
	i := interactor_mock.NewTodo(t)
	inputMock := presenter_mock.NewTodoInput(t)
	outputMock := presenter_mock.NewTodoOutput(t)
	if sqlMock, db, err = test.NewSqlMock(&gorm.Session{}); err != nil {
		require.NoError(t, err)
	}
	return cfg, inputMock, outputMock, i, sqlMock, db
}

func helperNewTodoAndContext(t *testing.T, method, path string, headers http.Header, body any, logger *slog.Logger) (echo.Context, *config.Config, *presenter_mock.TodoInput, *presenter_mock.TodoOutput, *interactor_mock.Todo, sqlmock.Sqlmock, *gorm.DB) {
	e := echo.New()
	require.NotNil(t, e)
	ctx := echo_helper.NewContext(e, method, path, headers, body, logger)
	cfg, inputMock, outputMock, i, sqlMock, db := helperNewTodo(t)
	return ctx, cfg, inputMock, outputMock, i, sqlMock, db
}

func helperAssertTodoExpectations(t *testing.T, inputMock *presenter_mock.TodoInput, outputMock *presenter_mock.TodoOutput, i *interactor_mock.Todo, sqlMock sqlmock.Sqlmock) {
	inputMock.AssertExpectations(t)
	outputMock.AssertExpectations(t)
	i.AssertExpectations(t)
	require.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestNewTodo(t *testing.T) {
	assert.PanicsWithError(t, common_err.ErrNil("cfg").Error(), func() {
		NewTodo(nil, nil, nil)
	})

	cfg := config.Get()
	assert.PanicsWithError(t, common_err.ErrNil("interactor").Error(), func() {
		NewTodo(cfg, nil, nil)
	})

	i := interactor_mock.NewTodo(t)
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
	ctx, cfg, inputMock, outputMock, i, sqlMock, db := helperNewTodoAndContext(t, http.MethodGet, path, http.Header{}, nil, slog.Default())
	expectedErr := "inbound data error"
	inputMock.On("GetAll", mock.Anything).Return(errors.New(expectedErr))

	p := newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)
	err = p.GetAllTodos(ctx)
	assert.EqualError(t, err, expectedErr)
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)

	// interactor error
	ctx, cfg, inputMock, outputMock, i, sqlMock, db = helperNewTodoAndContext(t, http.MethodGet, path, http.Header{}, nil, slog.Default())
	expectedErr = "interactor error"
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()
	inputMock.On("GetAll", mock.Anything).Return(nil)
	i.On("GetAll", mock.Anything).Return(nil, errors.New(expectedErr))

	p = newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)
	err = p.GetAllTodos(ctx)
	assert.EqualError(t, err, expectedErr)
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)

	// outbound error
	ctx, cfg, inputMock, outputMock, i, sqlMock, db = helperNewTodoAndContext(t, http.MethodGet, path, http.Header{}, nil, slog.Default())
	expectedErr = "outbound error"
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	inputMock.On("GetAll", mock.Anything).Return(nil)
	i.On("GetAll", mock.Anything).Return([]model.Todo{}, nil)
	outputMock.On("GetAll", mock.Anything, []model.Todo{}).Return(nil, errors.New(expectedErr))

	p = newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)
	err = p.GetAllTodos(ctx)
	assert.EqualError(t, err, expectedErr)
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)

	// Success case
	ctx, cfg, inputMock, outputMock, i, sqlMock, db = helperNewTodoAndContext(t, http.MethodGet, path, http.Header{}, nil, slog.Default())
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	inputMock.On("GetAll", mock.Anything).Return(nil)
	i.On("GetAll", mock.Anything).Return([]model.Todo{}, nil)
	outputMock.On("GetAll", mock.Anything, []model.Todo{}).Return([]public.ToDo{}, nil)

	p = newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)
	err = p.GetAllTodos(ctx)
	assert.NoError(t, err)
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)
}

func TestCreateTodo(t *testing.T) {
	const path = "/todo/v1/todos"
	var (
		db            *gorm.DB
		sqlMock       sqlmock.Sqlmock
		err           error
		dataInput     *public.ToDo
		data          *model.Todo
		processedData *model.Todo
		dataOutput    *public.ToDo
	)

	// inbound error
	dataInput = http_builder.NewToDo().Build()
	ctx, cfg, inputMock, outputMock, i, sqlMock, db := helperNewTodoAndContext(t, http.MethodPost, path, http.Header{}, dataInput, slog.Default())
	expectedErr := common_err.ErrNil("ctx")
	expectedErrStr := expectedErr.Error()
	expectedHttpErr := echo.NewHTTPError(http.StatusBadRequest, expectedErrStr)
	inputMock.On("Create", mock.Anything).Return(nil, expectedErr)
	e := echo.New()
	require.NotNil(t, e)
	p := newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)

	err = p.CreateTodo(ctx)
	require.Error(t, err)
	assert.EqualError(t, err, expectedHttpErr.Error())
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)

	// interactor error
	dataInput = http_builder.NewToDo().
		WithTitle("title").
		WithDescription("description").
		Build()
	ctx, cfg, inputMock, outputMock, i, sqlMock, db = helperNewTodoAndContext(t, http.MethodPost, path, http.Header{}, dataInput, slog.Default())
	expectedErrStr = "creating todo failed"
	expectedErr = fmt.Errorf("%s", expectedErrStr)
	expectedHttpErr = echo.NewHTTPError(http.StatusInternalServerError, expectedErrStr)
	data = model_builder.NewTodo().WithTitle(dataInput.Title).WithDescription(dataInput.Description).Build()
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()
	inputMock.On("Create", mock.Anything).Return(data, nil)
	i.On("Create", mock.Anything, data).Return(nil, expectedErr)
	p = newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)

	err = p.CreateTodo(ctx)
	require.Error(t, err)
	assert.EqualError(t, err, expectedHttpErr.Error())
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)

	// outbound error
	dataInput = http_builder.NewToDo().
		WithTitle("title").
		WithDescription("description").
		Build()
	ctx, cfg, inputMock, outputMock, i, sqlMock, db = helperNewTodoAndContext(t, http.MethodPost, path, http.Header{}, dataInput, slog.Default())
	expectedErrStr = "failure on formatting output response"
	expectedErr = fmt.Errorf("%s", expectedErrStr)
	expectedHttpErr = echo.NewHTTPError(http.StatusInternalServerError, expectedErrStr)
	data = model_builder.NewTodo().WithTitle(dataInput.Title).WithDescription(dataInput.Description).Build()
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	inputMock.On("Create", mock.Anything).Return(data, nil)
	i.On("Create", mock.Anything, data).Return(data, nil)
	outputMock.On("Create", mock.Anything, data).Return(nil, expectedErr)
	p = newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)

	err = p.CreateTodo(ctx)
	require.Error(t, err)
	assert.EqualError(t, err, expectedHttpErr.Error())
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)

	// Success scenario
	dataInput = http_builder.NewToDo().
		WithTitle("title").
		WithDescription("description").
		Build()
	ctx, cfg, inputMock, outputMock, i, sqlMock, db = helperNewTodoAndContext(t, http.MethodPost, path, http.Header{}, dataInput, slog.Default())
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
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	inputMock.On("Create", mock.Anything).Return(data, nil)
	i.On("Create", mock.Anything, data).Return(processedData, nil)
	outputMock.On("Create", mock.Anything, processedData).Return(dataOutput, nil)
	p = newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)

	err = p.CreateTodo(ctx)
	require.NoError(t, err)
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)
}

func TestDeleteTodo(t *testing.T) {
	var (
		id   = uuid.New()
		path = "/todo/v1/todos/" + id.String()
	)
	ctx, cfg, inputMock, outputMock, i, sqlMock, db := helperNewTodoAndContext(t, http.MethodDelete, path, http.Header{}, nil, slog.Default())
	p := newTodo(cfg, inputMock, outputMock, i, db)
	err := p.DeleteTodo(ctx, uuid.UUID{})
	require.Error(t, err)
	assert.EqualError(t, err, echo.ErrNotImplemented.Error())
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)
}

func TestGetTodo(t *testing.T) {
	var (
		id            openapi_types.UUID = uuid.New()
		path                             = "/todo/v1/todos/" + id.String()
		processedData                    = model_builder.NewTodo().WithID(id).Build()
		outputData                       = http_builder.NewToDo().
				WithID((*openapi_types.UUID)(&processedData.UUID)).
				WithTitle(processedData.Title).
				WithDescription(processedData.Description).
				WithDueDate(processedData.DueDate).
				Build()
		expectedErr error
	)

	// Manage interactor error
	ctx, cfg, inputMock, outputMock, i, sqlMock, db := helperNewTodoAndContext(t, http.MethodGet, path, http.Header{}, nil, slog.Default())
	expectedErrStr := "interactor error"
	expectedErr = errors.New(expectedErrStr)
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()
	i.On("GetByUUID", mock.Anything, id).Return(nil, expectedErr)
	p := newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)

	err := p.GetTodo(ctx, id)
	require.Error(t, err)
	assert.EqualError(t, err, expectedErrStr)
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)

	// Manage outbound error
	ctx, cfg, inputMock, outputMock, i, sqlMock, db = helperNewTodoAndContext(t, http.MethodGet, path, http.Header{}, nil, slog.Default())
	expectedErrStr = "outbound error"
	expectedErr = errors.New(expectedErrStr)
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	i.On("GetByUUID", mock.Anything, id).Return(processedData, nil)
	outputMock.On("Get", mock.Anything, processedData).Return(nil, expectedErr)
	p = newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)

	err = p.GetTodo(ctx, id)
	require.Error(t, err)
	assert.EqualError(t, err, expectedErrStr)
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)

	// Success use case
	ctx, cfg, inputMock, outputMock, i, sqlMock, db = helperNewTodoAndContext(t, http.MethodGet, path, http.Header{}, nil, slog.Default())
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	i.On("GetByUUID", mock.Anything, id).Return(processedData, nil)
	outputMock.On("Get", mock.Anything, processedData).Return(outputData, nil)
	p = newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)

	err = p.GetTodo(ctx, id)
	require.NoError(t, err)
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)
}

func TestPatchTodo(t *testing.T) {
	var (
		id   = uuid.New()
		path = "/todo/v1/todos/" + id.String()
	)

	ctx, cfg, inputMock, outputMock, i, sqlMock, db := helperNewTodoAndContext(t, http.MethodPatch, path, http.Header{}, nil, slog.Default())
	// sqlMock.ExpectBegin()
	// sqlMock.ExpectCommit()
	p := newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)

	err := p.PatchTodo(ctx, uuid.UUID{})
	require.Error(t, err)
	assert.EqualError(t, err, echo.ErrNotImplemented.Error())
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)
}

func TestUpdateTodo(t *testing.T) {
	var (
		id   = uuid.New()
		path = "/todo/v1/todos/" + id.String()
	)
	ctx, cfg, inputMock, outputMock, i, sqlMock, db := helperNewTodoAndContext(t, http.MethodPut, path, http.Header{}, nil, slog.Default())
	// sqlMock.ExpectBegin()
	// sqlMock.ExpectCommit()
	p := newTodo(cfg, inputMock, outputMock, i, db)
	require.NotNil(t, p)

	err := p.UpdateTodo(ctx, uuid.UUID{})
	require.Error(t, err)
	assert.EqualError(t, err, echo.ErrNotImplemented.Error())
	helperAssertTodoExpectations(t, inputMock, outputMock, i, sqlMock)
}
