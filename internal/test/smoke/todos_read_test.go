package smoke

import (
	"fmt"
	"net/http"

	builder_api "github.com/avisiedo/go-microservice-1/internal/test/builder/api/http"
)

// SuiteDeleteDomain is the suite to validate the smoke test when read domain endpoint at GET /api/idmsvc/v1/domains/:domain_id
type SuiteTodosRead struct {
	SuiteBaseOneTodo
}

func (s *SuiteTodosRead) SetupTest() {
	s.SuiteBase.SetupTest()
}

func (s *SuiteTodosRead) TearDownTest() {
	s.SuiteBase.TearDownTest()
}

func (s *SuiteTodosRead) TestReadDomain() {
	url := fmt.Sprintf("%s/%s/%s", s.DefaultPublicBaseURL(), "todos", s.Todos[0].TodoId)

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestTodosRead",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url,
				Header: http.Header{
					// FIXME Remove hardcoded header key
					"Request-Id": {"test_todos_read"},
				},
				Body: builder_api.NewToDo().Build(),
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusNoContent,
				Header: http.Header{
					// FIXME Remove hardcoded header key
					"Request-Id": {"test_todos_read"},
				},
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
