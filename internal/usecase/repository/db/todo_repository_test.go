package db

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"context"
	"fmt"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/avisiedo/go-microservice-1/internal/config"
	model "github.com/avisiedo/go-microservice-1/internal/domain/model"
	app_context "github.com/avisiedo/go-microservice-1/internal/infrastructure/context"
	"github.com/avisiedo/go-microservice-1/internal/test"
	"github.com/avisiedo/go-microservice-1/internal/test/builder/helper"
	builder "github.com/avisiedo/go-microservice-1/internal/test/builder/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SuiteTodo struct {
	suite.Suite
	ctx        context.Context
	cancel     context.CancelFunc
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository *todoRepository
}

// https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
func (s *SuiteTodo) SetupTest() {
	var err error
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.mock, s.DB, err = test.NewSqlMock(&gorm.Session{
		SkipHooks: true,
	})
	if err != nil {
		s.Suite.FailNow("Error calling gorm.Open: %s", err.Error())
		return
	}
	s.ctx = app_context.WithDB(context.Background(), s.DB)
	s.repository = &todoRepository{}
}

func (s *SuiteTodo) TestNewTodo() {
	t := s.Suite.T()
	assert.NotPanics(t, func() {
		cfg := config.Get()
		_ = NewTodo(cfg)
	})
}

//
// TestCreateTodo
//

func (s *SuiteTodo) prepareInsertIntoTodo(data *model.Todo, mock sqlmock.Sqlmock, expectedErr error) {
	var dataUUID any
	if (data.UUID == uuid.UUID{}) {
		dataUUID = sqlmock.AnyArg()
	} else {
		dataUUID = data.UUID
	}
	expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "todos" ("created_at","updated_at","deleted_at","uuid","title","description","due_date") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,

			dataUUID,
			data.Title,
			data.Description,
			data.DueDate,
		)
	if expectedErr != nil {
		expectQuery.WillReturnError(expectedErr)
	} else {
		expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(helper.GenRandNum(1, 100000000)))
	}
}

func (s *SuiteTodo) helperTestCreateTodo(stage int, data *model.Todo, mock sqlmock.Sqlmock, expectedErr error) {
	var err error
	for i := 1; i <= stage; i++ {
		if i == stage && expectedErr != nil {
			err = expectedErr
		} else {
			err = nil
		}
		switch i {
		case 1:
			s.prepareInsertIntoTodo(data, mock, err)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *SuiteTodo) TestCreateTodoNil() {
	t := s.Suite.T()
	s.helperTestCreateTodo(0, nil, s.mock, nil)
	_, err := s.repository.Create(s.ctx, nil)
	require.EqualError(t, err, "'todo' is nil")
}

func (s *SuiteTodo) TestCreateTodo() {
	t := s.Suite.T()
	// testUuid := uuid.New()
	var (
		err           error
		gormModel     gorm.Model  = builder.NewModel().WithID(0).Build()
		data          *model.Todo = builder.NewTodo().WithModel(gormModel).Build()
		result        *model.Todo
		expectedError error
	)

	// Check ctx is nil
	result, err = s.repository.Create(nil, data)
	require.Errorf(t, err, "")
	assert.Nil(t, result)

	// Check uuid == {}
	data.UUID = uuid.UUID{}
	result, err = s.repository.Create(s.ctx, data)
	require.Errorf(t, err, "")
	assert.Nil(t, result)

	// Check nil
	data.UUID = uuid.New()
	expectedError = fmt.Errorf("code=500, message='data' of type '*model.Todo' cannot be nil")
	s.helperTestCreateTodo(1, data, s.mock, expectedError)
	result, err = s.repository.Create(s.ctx, data)
	assert.EqualError(t, err, "code=500, message='data' of type '*model.Todo' cannot be nil")
	assert.Nil(t, result)

	// Error on INSERT INTO "todo"
	expectedError = fmt.Errorf(`error at INSERT INTO "todo"`)
	s.helperTestCreateTodo(1, data, s.mock, expectedError)
	result, err = s.repository.Create(s.ctx, data)
	assert.EqualError(t, err, expectedError.Error())
	assert.Nil(t, result)

	// Success scenario
	expectedError = nil
	s.helperTestCreateTodo(1, data, s.mock, nil)
	result, err = s.repository.Create(s.ctx, data)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

//
// TestDeleteByUUID
//

func (s *SuiteTodo) prepareDeleteByUUID(data *model.Todo, mock sqlmock.Sqlmock, expectedErr error) {
	var dataUUID any
	if (data.UUID == uuid.UUID{}) {
		dataUUID = sqlmock.AnyArg()
	} else {
		dataUUID = data.UUID
	}
	expectQuery := s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "todos" WHERE uuid = $1`)).
		WithArgs(
			dataUUID,
		)
	if expectedErr != nil {
		expectQuery.WillReturnError(expectedErr)
	} else {
		expectQuery.WillReturnResult(sqlmock.NewResult(0, 1))
	}
}

func (s *SuiteTodo) helperTestDeleteByUUID(stage int, data *model.Todo, mock sqlmock.Sqlmock, expectedErr error) {
	var err error
	for i := 1; i <= stage; i++ {
		if i == stage && expectedErr != nil {
			err = expectedErr
		} else {
			err = nil
		}
		switch i {
		case 1:
			s.prepareDeleteByUUID(data, mock, err)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *SuiteTodo) TestDeleteByUUID() {
	t := s.Suite.T()
	var (
		err           error
		gormModel     gorm.Model  = builder.NewModel().Build()
		data          *model.Todo = builder.NewTodo().WithModel(gormModel).WithID(uuid.New()).Build()
		expectedError error
	)

	// Check nil
	expectedError = fmt.Errorf("'ctx' is nil")
	s.helperTestDeleteByUUID(0, data, s.mock, expectedError)
	err = s.repository.DeleteByUUID(nil, uuid.UUID{})
	require.EqualError(t, err, expectedError.Error())

	// Check empty uuid
	expectedError = fmt.Errorf("'todo_uuid' is empty")
	s.helperTestDeleteByUUID(0, data, s.mock, expectedError)
	err = s.repository.DeleteByUUID(s.ctx, uuid.UUID{})
	require.EqualError(t, err, expectedError.Error())

	// Success scenario
	expectedError = nil
	s.helperTestDeleteByUUID(1, data, s.mock, expectedError)
	err = s.repository.DeleteByUUID(s.ctx, data.UUID)
	assert.NoError(t, err)
}

//
// TestGetByUUID
//

func (s *SuiteTodo) prepareGetByUUID(data *model.Todo, mock sqlmock.Sqlmock, expectedErr error) {
	var dataUUID any
	if (data.UUID == uuid.UUID{}) {
		dataUUID = sqlmock.AnyArg()
	} else {
		dataUUID = data.UUID
	}
	expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE uuid = $1 AND "todos"."deleted_at" IS NULL ORDER BY "todos"."id" LIMIT $2`)).
		WithArgs(
			dataUUID,
			int64(1),
		)
	if expectedErr != nil {
		expectQuery.WillReturnError(expectedErr)
	} else {
		expectQuery.WillReturnRows(sqlmock.NewRows([]string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt", "UUID", "Title"}).
			AddRow(
				data.ID,
				data.CreatedAt,
				data.UpdatedAt,
				data.DeletedAt,

				data.UUID,
				data.Title,
			),
		)
	}
}

func (s *SuiteTodo) helperTestGetByUUID(stage int, data *model.Todo, mock sqlmock.Sqlmock, expectedErr error) {
	var err error
	for i := 1; i <= stage; i++ {
		if i == stage && expectedErr != nil {
			err = expectedErr
		} else {
			err = nil
		}
		switch i {
		case 1:
			s.prepareGetByUUID(data, mock, err)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *SuiteTodo) TestGetByUUID() {
	t := s.Suite.T()
	var (
		err           error
		gormModel     gorm.Model  = builder.NewModel().Build()
		data          *model.Todo = builder.NewTodo().WithModel(gormModel).WithID(uuid.New()).Build()
		expectedError error
		result        *model.Todo = nil
	)

	// Check nil
	expectedError = fmt.Errorf("'ctx' is nil")
	s.helperTestGetByUUID(0, data, s.mock, expectedError)
	result, err = s.repository.GetByUUID(nil, uuid.UUID{})
	require.EqualError(t, err, expectedError.Error())
	assert.Nil(t, result)

	// Check empty ID
	expectedError = fmt.Errorf("'id' is empty")
	s.helperTestGetByUUID(0, data, s.mock, expectedError)
	result, err = s.repository.GetByUUID(s.ctx, uuid.UUID{})
	require.EqualError(t, err, expectedError.Error())
	assert.Nil(t, result)

	// Check error requesting by ID
	expectedError = fmt.Errorf("some error")
	s.helperTestGetByUUID(1, data, s.mock, expectedError)
	result, err = s.repository.GetByUUID(s.ctx, data.UUID)
	require.EqualError(t, err, expectedError.Error())
	assert.Nil(t, result)

	// Check empty result
	expectedError = nil
	s.helperTestGetByUUID(1, data, s.mock, expectedError)
	result, err = s.repository.GetByUUID(s.ctx, data.UUID)
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, data.CreatedAt, result.CreatedAt)
}

func (s *SuiteTodo) prepareCountGetAll(data []model.Todo, mock sqlmock.Sqlmock, expectedErr error) {
	expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "todos" WHERE "todos"."deleted_at" IS NULL`))
	if expectedErr != nil {
		expectQuery.WillReturnError(expectedErr)
	} else {
		rows := sqlmock.NewRows([]string{"count(*)"})
		rows.AddRow(int64(len(data)))
		expectQuery.WillReturnRows(rows)
	}
}

func (s *SuiteTodo) prepareGetAll(data []model.Todo, mock sqlmock.Sqlmock, expectedErr error) {
	expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "todos"."deleted_at" IS NULL`))
	if expectedErr != nil {
		expectQuery.WillReturnError(expectedErr)
	} else {
		rows := sqlmock.NewRows([]string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt", "UUID", "Title"})
		for idx := range data {
			rows.AddRow(
				data[idx].ID,
				data[idx].CreatedAt,
				data[idx].UpdatedAt,
				data[idx].DeletedAt,

				data[idx].UUID,
				data[idx].Title,
			)
		}
		expectQuery.WillReturnRows(rows)
	}
}

func (s *SuiteTodo) helperTestGetAll(stage int, data []model.Todo, mock sqlmock.Sqlmock, expectedErr error) {
	var err error
	for i := 1; i <= stage; i++ {
		if i == stage && expectedErr != nil {
			err = expectedErr
		} else {
			err = nil
		}
		switch i {
		case 1:
			s.prepareCountGetAll(data, mock, err)
		case 2:
			s.prepareGetAll(data, mock, err)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func (s *SuiteTodo) TestGetAll() {
	t := s.Suite.T()
	var (
		err           error
		gormModel     gorm.Model   = builder.NewModel().Build()
		data          []model.Todo = []model.Todo{*builder.NewTodo().WithModel(gormModel).WithID(uuid.New()).Build()}
		expectedError error
		result        []model.Todo = nil
	)

	// Check nil
	expectedError = fmt.Errorf("'ctx' is nil")
	s.helperTestGetAll(0, data, s.mock, expectedError)
	result, err = s.repository.GetAll(nil)
	require.EqualError(t, err, expectedError.Error())
	assert.Nil(t, result)

	// Check error requesting all
	expectedError = fmt.Errorf("some error")
	s.helperTestGetAll(1, data, s.mock, expectedError)
	result, err = s.repository.GetAll(s.ctx)
	require.EqualError(t, err, expectedError.Error())
	require.NotNil(t, result)
	assert.Empty(t, result)

	// Check error counting
	expectedError = fmt.Errorf("some error")
	s.helperTestGetAll(1, data, s.mock, expectedError)
	result, err = s.repository.GetAll(s.ctx)
	require.EqualError(t, err, expectedError.Error())
	require.NotNil(t, result)
	assert.Empty(t, result)

	// Check empty result
	expectedError = nil
	// Only run the COUNT(*) not the select
	s.helperTestGetAll(1, []model.Todo{}, s.mock, expectedError)
	result, err = s.repository.GetAll(s.ctx)
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result)

	// Success scenario
	expectedError = nil
	s.helperTestGetAll(2, data, s.mock, expectedError)
	result, err = s.repository.GetAll(s.ctx)
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result)
	for idx := range result {
		assert.Equal(t, data[idx].CreatedAt, result[idx].CreatedAt)
		assert.Equal(t, data[idx].UpdatedAt, result[idx].UpdatedAt)
		assert.Equal(t, gorm.DeletedAt{}, result[idx].DeletedAt)

		assert.Equal(t, data[idx].Title, result[idx].Title)
	}
}
