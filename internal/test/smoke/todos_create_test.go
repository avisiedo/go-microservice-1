package smoke

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	builder_api "github.com/avisiedo/go-microservice-1/internal/test/builder/api/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SuiteDeleteDomain is the suite to validate the smoke test when read domain endpoint at GET /api/idmsvc/v1/domains/:domain_id
type SuiteTodosCreate struct {
	SuiteBase
}

func (s *SuiteTodosCreate) SetupTest() {
	s.SuiteBase.SetupTest()
}

func (s *SuiteTodosCreate) TearDownTest() {
	s.SuiteBase.TearDownTest()
}

func (s *SuiteTodosCreate) TestTodosCreate() {
	t := s.T()
	url := fmt.Sprintf("%s/%s", s.DefaultPublicBaseURL(), "todos")
	resource := builder_api.NewToDo().Build()

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestTodosCreate",
			Given: TestCaseGiven{
				Method: http.MethodPost,
				URL:    url,
				Header: http.Header{
					"Request-Id": {"test_todos_create"},
				},
				Body: resource,
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusCreated,
				Header: http.Header{
					"Request-Id": {"test_todos_create"},
				},
				BodyFunc: WrapBodyFuncTodoResponse(t, func(t *testing.T, body *public.ToDo) {
					require.NotNil(t, body)
					assert.Equal(t, resource.Title, body.Title)
					assert.Equal(t, resource.Description, body.Description)
					assert.Equal(t, resource.DueDate, body.DueDate)
				}),
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
