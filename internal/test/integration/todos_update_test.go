package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/api/header"
	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	builder_api "github.com/avisiedo/go-microservice-1/internal/test/builder/api/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type BodyFuncTodoUpdate func(t *testing.T, body *public.ToDo)

func WrapBodyFuncTodoResponse(t *testing.T, expected BodyFuncTodoUpdate) BodyFunc {
	if t == nil {
		panic(common_err.ErrNil("t"))
	}
	if expected == nil {
		return func(t *testing.T, body []byte) {
			require.Equal(t, 0, len(body))
		}
	}
	return func(t *testing.T, body []byte) {
		var data *public.ToDo = &public.ToDo{}
		if err := json.Unmarshal(body, data); err != nil {
			require.FailNow(t, fmt.Sprintf("Error unmarshalling body:\n"+
				"error: %q",
				err.Error(),
			))
		}

		// Run body expectetion on the unserialized data
		expected(t, data)
	}
}

// SuiteTodosUpdate is the suite to validate the smoke test when update todos endpoint at PUT /api/todo/v1/todos/:todoId
type SuiteTodosUpdate struct {
	SuiteBaseOneTodo
}

func (s *SuiteTodosUpdate) SetupTest() {
	s.SuiteBaseTodos.SetupTest()
}

func (s *SuiteTodosUpdate) TearDownTest() {
	s.SuiteBaseTodos.TearDownTest()
}

func (s *SuiteTodosUpdate) TestUpdateDomain() {
	t := s.T()
	url := fmt.Sprintf("%s/%s/%s", s.DefaultPublicBaseURL(), "domains", s.Todos.TodoId)
	updateTodo := builder_api.NewToDo().
		WithID(s.Todos.TodoId).
		WithTitle("updated title").
		WithDescription("updated description").
		Build()

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestTodosUpdate",
			Given: TestCaseGiven{
				Method: http.MethodPut,
				URL:    url,
				Header: http.Header{
					header.HdrRequestID: {"test_todos_update"},
				},
				Body: updateTodo,
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					header.HdrRequestID: {"test_todos_update"},
				},
				BodyFunc: WrapBodyFuncTodoResponse(t, func(t *testing.T, body *public.ToDo) {
					assert.Equal(t, updateTodo.TodoId, body.TodoId)
				}),
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
