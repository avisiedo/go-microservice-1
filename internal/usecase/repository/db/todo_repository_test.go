package db

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"context"
	"fmt"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/avisiedo/go-microservice-1/internal/config"
	model "github.com/avisiedo/go-microservice-1/internal/domain/model"
	"github.com/avisiedo/go-microservice-1/internal/test"
	builder "github.com/avisiedo/go-microservice-1/internal/test/builder/model"
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
	s.ctx = ContextWithDb(context.Background(), s.DB)
	s.repository = &todoRepository{}
}

func (s *SuiteTodo) TestNewTodo() {
	t := s.Suite.T()
	assert.NotPanics(t, func() {
		cfg := config.Get()
		_ = NewTodo(cfg)
	})
}

func (s *SuiteTodo) helperTestCreateTodo(stage int, data *model.Todo, mock sqlmock.Sqlmock, expectedErr error) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			expectQuery := s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "todos" ("created_at","updated_at","deleted_at","uuid","title","description","due_date") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
				WithArgs(
					data.Model.CreatedAt,
					data.Model.UpdatedAt,
					nil,

					data.UUID,
					data.Title,
					data.Description,
					data.DueDate,
				)
			if i == stage && expectedErr != nil {
				expectQuery.WillReturnError(expectedErr)
			} else {
				expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(data.ID))
			}
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
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

	// Check nil
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
