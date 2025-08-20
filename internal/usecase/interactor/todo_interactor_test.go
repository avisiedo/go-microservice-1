package interactor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	db "github.com/avisiedo/go-microservice-1/internal/test/mock/interface/repository/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func helperNewTodo(t *testing.T) (*db.TodoRepository, *todoInteractor) {
	mockRepo := db.NewTodoRepository(t)
	sut := newTodo(mockRepo)
	return mockRepo, sut
}

func TestNewTodo(t *testing.T) {
	// Check panic when repo is nil
	assert.PanicsWithError(t, common_err.ErrNil("repo").Error(), func() {
		newTodo(nil)
	})

	// Success scenario for newTodo
	mockRepo, result1 := helperNewTodo(t)
	assert.NotNil(t, result1)
	assert.IsType(t, &todoInteractor{}, result1)

	// Success scenario for NewTodo
	result2 := NewTodo(mockRepo)
	assert.NotNil(t, result2)
}

func TestCheckCtx(t *testing.T) {
	_, sut := helperNewTodo(t)
	assert.EqualError(t, sut.checkCtx(nil), common_err.ErrNil("ctx").Error())

	_, sut = helperNewTodo(t)
	assert.NoError(t, sut.checkCtx(context.Background()))
}

func TestCheckCtxAndTodo(t *testing.T) {
	_, sut := helperNewTodo(t)
	assert.EqualError(t, sut.checkCtxAndTodo(nil, nil), common_err.ErrNil("ctx").Error())

	_, sut = helperNewTodo(t)
	assert.EqualError(t, sut.checkCtxAndTodo(context.Background(), nil), common_err.ErrNil("todo").Error())

	_, sut = helperNewTodo(t)
	assert.NoError(t, sut.checkCtxAndTodo(context.Background(), &model.Todo{}))
}

func TestCheckCtxAndId(t *testing.T) {
	_, sut := helperNewTodo(t)
	assert.EqualError(t, sut.checkCtxAndId(nil, uuid.UUID{}), common_err.ErrNil("ctx").Error())

	_, sut = helperNewTodo(t)
	assert.EqualError(t, sut.checkCtxAndId(context.Background(), uuid.UUID{}), common_err.ErrEmpty("id").Error())

	_, sut = helperNewTodo(t)
	assert.NoError(t, sut.checkCtxAndId(context.Background(), uuid.New()))
}

func TestCheckCreateTodo(t *testing.T) {
	_, sut := helperNewTodo(t)
	assert.EqualError(t, sut.checkCreateTodo(nil, nil), common_err.ErrNil("ctx").Error())

	_, sut = helperNewTodo(t)
	assert.EqualError(t, sut.checkCreateTodo(context.Background(), nil), common_err.ErrNil("todo").Error())

	_, sut = helperNewTodo(t)
	assert.EqualError(t, sut.checkCreateTodo(context.Background(), &model.Todo{}), common_err.ErrEmpty("todo.Description").Error())

	_, sut = helperNewTodo(t)
	assert.NoError(t, sut.checkCreateTodo(context.Background(), &model.Todo{Description: "Some description"}))
}

func TestCreateTodo(t *testing.T) {
	var (
		mockRepo *db.TodoRepository
		sut      *todoInteractor
	)

	// When check fails
	mockRepo, sut = helperNewTodo(t)
	result, err := sut.Create(nil, nil)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, common_err.ErrNil("ctx").Error())

	// When mockRepo fails
	ctx := context.Background()
	data := &model.Todo{
		Title:       "First TODO",
		Description: "A test first TODO",
	}
	expectedErrStr := "db lost connection"
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("Create", ctx, data).Return(nil, errors.New(expectedErrStr))
	result, err = sut.Create(ctx, data)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, expectedErrStr)
	mockRepo.AssertExpectations(t)

	// Success case
	dataResult := &model.Todo{
		Model: gorm.Model{
			ID:        3,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UUID:        uuid.New(),
		Title:       data.Title,
		Description: data.Description,
	}
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("Create", ctx, data).Return(dataResult, nil)
	result, err = sut.Create(ctx, data)
	require.NotNil(t, result)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	assert.Equal(t, result, dataResult)
}

func TestUpdateTodo(t *testing.T) {
	var (
		mockRepo *db.TodoRepository
		sut      *todoInteractor
	)

	// When check fails
	mockRepo, sut = helperNewTodo(t)
	result, err := sut.Update(nil, nil)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, common_err.ErrNil("ctx").Error())

	// When mockRepo fails
	ctx := context.Background()
	data := &model.Todo{
		Title:       "First TODO",
		Description: "A test first TODO",
	}
	expectedErrStr := "db lost connection"
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("Update", ctx, data).Return(nil, errors.New(expectedErrStr))
	result, err = sut.Update(ctx, data)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, expectedErrStr)
	mockRepo.AssertExpectations(t)

	// Success case
	dataResult := &model.Todo{
		Model: gorm.Model{
			ID:        3,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UUID:        uuid.New(),
		Title:       data.Title,
		Description: data.Description,
	}
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("Update", ctx, data).Return(dataResult, nil)
	result, err = sut.Update(ctx, data)
	require.NotNil(t, result)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	assert.Equal(t, result, dataResult)
}

func TestGetByUUID(t *testing.T) {
	var (
		mockRepo *db.TodoRepository
		sut      *todoInteractor
	)

	// When check fails
	_, sut = helperNewTodo(t)
	result, err := sut.GetByUUID(nil, uuid.UUID{})
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, common_err.ErrNil("ctx").Error())

	// When mockRepo fails
	ctx := context.Background()
	defer ctx.Done()
	id := uuid.New()
	expectedErrStr := "db lost connection"
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("GetByUUID", ctx, id).Return(nil, errors.New(expectedErrStr))
	result, err = sut.GetByUUID(ctx, id)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, expectedErrStr)
	mockRepo.AssertExpectations(t)

	// Sucess use case
	ctx = context.Background()
	defer ctx.Done()
	dataResult := &model.Todo{
		Model: gorm.Model{
			ID:        3,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UUID:        uuid.New(),
		Title:       "First TODO",
		Description: "A test first TODO",
	}
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("GetByUUID", ctx, id).Return(dataResult, nil)
	result, err = sut.GetByUUID(ctx, id)
	require.NotNil(t, result)
	assert.NoError(t, err)
	assert.Equal(t, result, dataResult)
	mockRepo.AssertExpectations(t)
}

func TestGetAllTodo(t *testing.T) {
	var (
		mockRepo *db.TodoRepository
		sut      *todoInteractor
	)

	// When check fails
	_, sut = helperNewTodo(t)
	result, err := sut.GetAll(nil)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, common_err.ErrNil("ctx").Error())

	// When mockRepo fails
	ctx := context.Background()
	defer ctx.Done()
	dataResult := []model.Todo{}
	expectedErrStr := "db lost connection"
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("GetAll", ctx).Return(nil, errors.New(expectedErrStr))
	result, err = sut.GetAll(ctx)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, expectedErrStr)
	mockRepo.AssertExpectations(t)

	// Sucess use case
	ctx = context.Background()
	defer ctx.Done()
	dataResult = []model.Todo{
		{
			Model: gorm.Model{
				ID:        3,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UUID:        uuid.New(),
			Title:       "First TODO",
			Description: "A test first TODO",
		},
	}
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("GetAll", ctx).Return(dataResult, nil)
	result, err = sut.GetAll(ctx)
	require.NotNil(t, result)
	assert.NoError(t, err)
	assert.Equal(t, result, dataResult)
	mockRepo.AssertExpectations(t)
}

func TestDeleteByUUIDTodo(t *testing.T) {
	var (
		mockRepo *db.TodoRepository
		sut      *todoInteractor
	)

	// When check fails
	_, sut = helperNewTodo(t)
	err := sut.DeleteByUUID(nil, uuid.UUID{})
	require.Error(t, err)
	assert.EqualError(t, err, common_err.ErrNil("ctx").Error())

	// When mockRepo fails
	ctx := context.Background()
	defer ctx.Done()
	id := uuid.New()
	expectedErrStr := "db lost connection"
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("DeleteByUUID", ctx, id).Return(errors.New(expectedErrStr))
	err = sut.DeleteByUUID(ctx, id)
	require.Error(t, err)
	assert.EqualError(t, err, expectedErrStr)
	mockRepo.AssertExpectations(t)

	// Sucess use case
	ctx = context.Background()
	defer ctx.Done()
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("DeleteByUUID", ctx, id).Return(nil)
	err = sut.DeleteByUUID(ctx, id)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPatchTodo(t *testing.T) {
	var (
		mockRepo *db.TodoRepository
		sut      *todoInteractor
	)

	// When check fails
	mockRepo, sut = helperNewTodo(t)
	result, err := sut.Patch(nil, nil)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, common_err.ErrNil("ctx").Error())

	// When mockRepo.GetByUUID fails
	ctx := context.Background()
	defer ctx.Done()
	id := uuid.New()
	dueDate := time.Now().Add(24 * time.Hour)
	data := &model.Todo{
		Model: gorm.Model{
			ID:        3,
			CreatedAt: time.Now(),
		},
		UUID:        id,
		Title:       "First TODO",
		Description: "A test first TODO",
		DueDate:     &dueDate,
	}
	expectedErrStr := "db lost connection"
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("GetByUUID", ctx, id).Return(nil, errors.New(expectedErrStr))
	result, err = sut.Patch(ctx, data)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, expectedErrStr)
	mockRepo.AssertExpectations(t)

	// When mockRepo.Update fails
	ctx = context.Background()
	defer ctx.Done()
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("GetByUUID", ctx, id).Return(data, nil)
	mockRepo.On("Update", ctx, data).Return(nil, errors.New(expectedErrStr))
	result, err = sut.Patch(ctx, data)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.EqualError(t, err, expectedErrStr)
	mockRepo.AssertExpectations(t)

	// Success case
	dataResult := &model.Todo{
		Model: gorm.Model{
			ID:        3,
			CreatedAt: data.Model.CreatedAt,
			UpdatedAt: time.Now(),
		},
		UUID:        id,
		Title:       data.Title + "Updated",
		Description: data.Description + "Updated",
		DueDate:     &dueDate,
	}
	mockRepo, sut = helperNewTodo(t)
	mockRepo.On("GetByUUID", ctx, id).Return(data, nil)
	mockRepo.On("Update", ctx, data).Return(dataResult, nil)
	result, err = sut.Patch(ctx, data)
	require.NotNil(t, result)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	assert.Equal(t, result, dataResult)
}
