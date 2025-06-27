package integration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	if value, exist := os.LookupEnv("TEST"); !exist || value != "integration" {
		t.Skip("This TestSuite require to start infrastructure: TEST=integration make compose-up test-integration")
	}
	// TODO Add here your test suites
	suite.Run(t, new(SuiteTodosCreate))
	suite.Run(t, new(SuiteTodosRead))
	// suite.Run(t, new(SuiteTodosDelete))
	// suite.Run(t, new(SuiteTodosUpdate))
	// suite.Run(t, new(SuiteTodosList))
}
