package smoke

import (
	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	builder_api "github.com/avisiedo/go-microservice-1/internal/test/builder/api/http"
)

// SuiteTodosDelete is the suite to validate the smoke test when read domain endpoint at GET /api/idmsvc/v1/domains/:domain_id
type SuiteBaseOneTodo struct {
	SuiteBase
	Todos []public.ToDo
}

func (s *SuiteBaseOneTodo) SetupTest() {
	s.SuiteBase.SetupTest()

	s.Todos = []public.ToDo{}
	oneTodo, err := s.TodosCreate(builder_api.NewToDo().Build())
	if err != nil {
		s.FailNow("error creating todo", err.Error())
	}
	s.Todos = append(s.Todos, *oneTodo)
}

func (s *SuiteBaseOneTodo) TearDownTest() {
	s.Todos = nil

	s.SuiteBase.TearDownTest()
}

func (s *SuiteBaseOneTodo) TestReadDomain() {
	t := s.T()
	t.Skip("Skipping wrapped test to avoid duplication")
}
