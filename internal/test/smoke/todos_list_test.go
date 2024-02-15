package smoke

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	builder_api "github.com/avisiedo/go-microservice-1/internal/test/builder/api/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SuiteTokenCreate is the suite token for smoke tests at /api/idmsvc/v1/domains/token
type SuiteTodosList struct {
	SuiteBase
	Resources []public.ToDo
}

// BodyFuncListResponse is the function that wrap
type BodyFuncTodoList func(t *testing.T, body []public.ToDo)

// WrapBodyFuncListResponse allow to implement custom body expectations for the specific type of the response.
// expected is the specific BodyFuncDomain for Domain type
// Returns a BodyFunc that wrap the generic expectation function.
func WrapBodyFuncTodoList(t *testing.T, expected BodyFuncTodoList) BodyFunc {
	if expected == nil {
		return func(t *testing.T, body []byte) {
			require.Equal(t, 0, len(body))
		}
	}
	return func(t *testing.T, body []byte) {
		var data []public.ToDo
		if err := json.Unmarshal(body, &data); err != nil {
			require.FailNow(t, fmt.Sprintf("Error unmarshalling body:\n"+
				"error: %q",
				err.Error(),
			))
		}

		// Run body expectetion on the unserialized data
		expected(t, data)
	}
}

func (s *SuiteTodosList) SetupTest() {
	s.SuiteBase.SetupTest()

	s.Resources = []public.ToDo{}
	for i := 1; i < 50; i++ {
		resource, err := s.TodosCreate(builder_api.NewToDo().Build())
		if err != nil {
			s.FailNow(err.Error())
		}
		s.Resources = append(s.Resources, *resource)
	}
}

func (s *SuiteTodosList) TearDownTest() {
	s.Resources = nil

	s.SuiteBase.TearDownTest()
}

func (s *SuiteTodosList) existTodo(id string) bool {
	for j := range s.Resources {
		if s.Resources[j].TodoId.String() == id {
			return true
		}
	}
	return false
}

func (s *SuiteTodosList) assertInTodos(t *testing.T, data []public.ToDo, msgAndArgs ...any) bool {
	if data == nil {
		return true
	}
	if len(data) == 0 {
		return true
	}

	for i := range data {
		resourceID := data[i].TodoId.String()
		if !s.existTodo(resourceID) {
			return assert.Fail(t, fmt.Sprintf("Not in slice: TodoID=%s\n", resourceID), msgAndArgs...)
		}
	}

	return true
}

func (s *SuiteTodosList) TestTodosList() {
	t := s.T()
	req, err := http.NewRequest(http.MethodGet, s.DefaultPublicBaseURL()+"/todos", nil)
	require.NoError(t, err)
	q := req.URL.Query()
	q.Add("offset", "0")
	q.Add("limit", "10")
	url1 := req.URL.String() + "?" + q.Encode()
	q.Set("offset", "40")
	q.Set("limit", "10")
	url2 := req.URL.String() + "?" + q.Encode()
	q.Set("offset", "20")
	q.Set("limit", "10")
	url3 := req.URL.String() + "?" + q.Encode()
	q.Del("offset")
	q.Del("limit")
	url4 := req.URL.String() + "?" + q.Encode()

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestTodosList: offset=0&limit=10 case",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url1,
				Header: http.Header{
					"Request-Id": {"test_todos_list_page_1"},
				},
				Body: nil,
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					"Request-Id": {"test_todos_list_page_1"},
				},
				BodyFunc: WrapBodyFuncTodoList(t, func(t *testing.T, body []public.ToDo) {
					require.NotNil(t, body)
					require.Equal(t, 10, len(body))

					// Check items
					s.assertInTodos(t, body)
				}),
			},
		},
		{
			Name: "TestTodosList: offset=40&limit=10 case",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url2,
				Header: http.Header{
					"Request-Id": {"test_todos_list_page_2"},
				},
				Body: builder_api.NewToDo().Build(),
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					"Request-Id": {"test_todos_list_page_last"},
				},
				BodyFunc: WrapBodyFuncTodoList(t, func(t *testing.T, body []public.ToDo) {
					require.NotNil(t, body)

					// Check items
					assert.Equal(t, 9, len(body))
					s.assertInTodos(t, body)
				}),
			},
		},
		{
			Name: "TestTodosList: offset=20&limit=10 case",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url3,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					"Request-Id": {"test_todos_list_page_last"},
				},
				Body: builder_api.NewToDo().Build(),
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					"Request-Id": {"test_todos_list_page_last"},
				},
				BodyFunc: WrapBodyFuncTodoList(t, func(t *testing.T, body []public.ToDo) {
					require.NotNil(t, body)

					// Check items
					assert.Equal(t, 10, len(body))
					s.assertInTodos(t, body)
				}),
			},
		},
		{
			Name: "TestTodosList: no params",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url4,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					"Request-Id": {"test_todos_list_page_last"},
				},
				Body: builder_api.NewToDo().Build(),
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					"Request-Id": {"test_todos_list_page_last"},
				},
				BodyFunc: WrapBodyFuncTodoList(t, func(t *testing.T, body []public.ToDo) {
					require.NotNil(t, body)

					// Check items
					assert.Equal(t, 10, len(body))
					s.assertInTodos(t, body)
				}),
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
