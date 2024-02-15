package smoke

import (
	"fmt"
	"net/http"

	builder_api "github.com/avisiedo/go-microservice-1/internal/test/builder/api/http"
)

// SuiteTodosDelete is the suite to validate the smoke test when read domain endpoint at GET /api/idmsvc/v1/domains/:domain_id
type SuiteTodosDelete struct {
	SuiteBaseOneTodo
}

func (s *SuiteTodosDelete) SetupTest() {
	s.SuiteBaseOneTodo.SetupTest()
}

func (s *SuiteTodosDelete) TearDownTest() {
	s.SuiteBaseOneTodo.TearDownTest()
}

func (s *SuiteTodosDelete) TestDeleteDomain() {
	url := fmt.Sprintf("%s/%s/%s", s.DefaultPublicBaseURL(), "todos", s.Todos[0].TodoId)

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestTodosDelete",
			Given: TestCaseGiven{
				Method: http.MethodDelete,
				URL:    url,
				Header: http.Header{
					"Request-Id": {"test_todos_delete"},
				},
				Body: builder_api.NewToDo().Build(),
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusNoContent,
				Header: http.Header{
					"Request-Id": {"test_todos_delete"},
				},
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
