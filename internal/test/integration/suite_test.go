package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/datastore"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/logger"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/service"
	service_impl "github.com/avisiedo/go-microservice-1/internal/infrastructure/service/impl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// SuiteBase represents the base Suite to be used for smoke tests, this
// start the services before run the smoke tests.
// TODO the smoke tests cannot be executed in parallel yet, an alternative
// for them would be to use specific http and metrics service in one side,
// and to use a specific OrgID per test by using a generator for it which
// would provide data partition between the tests.
type Suite struct {
	suite.Suite

	cancel context.CancelFunc
	svc    service.ApplicationService
	mocks  map[string]service.ApplicationService
	wg     *sync.WaitGroup
	db     *gorm.DB

	cfg *config.Config
}

func (s *Suite) StartMocks() {
	// TODO Add here the logic to start the mocks
	// s.mock["rbac"] = mock_rbac.New(...)
}

func (s *Suite) StopMocks() {
	for _, mock := range s.mocks {
		mock.Stop()
	}
}

func (s *Suite) WaitMocks() {
	for range /* mock */ s.mocks {
		// TODO Add the logic to await the mock is ready to use
	}
}

// SetupTest start the services and await until they are ready
// for being used.
func (s *Suite) SetupTest() {
	s.cfg = config.Get()
	s.wg = &sync.WaitGroup{}
	logger.InitLogger(s.cfg)
	s.db = datastore.NewDB(s.cfg, slog.Default())
	s.mocks = map[string]service.ApplicationService{}

	ctx, cancel := StartSignalHandler(context.Background())
	s.cancel = cancel
	s.svc = service_impl.NewApplication(ctx, s.wg, s.cfg, s.db)
	s.StartMocks()
	go func() {
		if e := s.svc.Start(); e != nil {
			panic(e)
		}
	}()
	s.WaitReady(s.cfg)
	s.WaitMocks()
}

// TearDownTest Stop the services in an ordered way before every
// smoke test executed.
func (s *Suite) TearDownTest() {
	TearDownSignalHandler()
	defer datastore.Close(s.db)
	defer s.cancel()
	s.svc.Stop()
	s.StopMocks()
	s.wg.Wait()
}

// WaitReady poll the ready healthcheck until the response is http.StatusOK
// cfg is the current configuration to use for the application.
func (s *Suite) WaitReady(cfg *config.Config) {
	t := s.T()
	if cfg == nil {
		panic("cfg is nil")
	}
	header := http.Header{}
	path := s.DefaultHealthcheckBaseURL() + "/readyz"
	for i := 0; i < 10; i++ {
		resp, err := s.DoRequest(
			http.MethodGet,
			path,
			header,
			nil,
		)
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
		t.Logf("Not available '%s'", path)
		time.Sleep(1000 * time.Millisecond)
	}
	panic("WaitReady didn't return after 30 seconds checking for it")
}

// RunTestCase run test for one specific testcase
func (s *Suite) RunTestCase(testCase *TestCase) {
	t := s.T()

	var (
		body []byte
		resp *http.Response
		err  error
	)

	// GIVEN testCase
	bodyCount := 0
	if testCase.Given.Body != nil {
		bodyCount++
	}
	if testCase.Given.BodyBytes != nil {
		bodyCount++
	}
	if bodyCount > 1 {
		t.Errorf("Given Body and BodyBytes are exclusive between them.")
	}
	bodyCount = 0
	if testCase.Expected.Body != nil {
		bodyCount++
	}
	if testCase.Expected.BodyFunc != nil {
		bodyCount++
	}
	if testCase.Expected.BodyBytes != nil {
		bodyCount++
	}
	if bodyCount > 1 {
		t.Errorf("Expected Body, BodyFunc and BodyBytes are exclusive between them.")
	}

	// WHEN
	resp, err = s.DoRequest(testCase.Given.Method, testCase.Given.URL, testCase.Given.Header, testCase.Given.Body)

	// THEN

	// Check no error
	require.NoError(t, err)
	if resp != nil {
		body, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)
	}

	// Check response status code
	t.Log("Checking response status")
	require.Equal(t, testCase.Expected.StatusCode, resp.StatusCode)

	// Check response headers
	t.Log("Checking response headers")
	for key := range testCase.Expected.Header {
		expectedValue := fmt.Sprintf("%s: %s", key, testCase.Expected.Header.Get(key))
		currentValue := fmt.Sprintf("%s: %s", key, resp.Header.Get(key))
		assert.Equal(t, expectedValue, currentValue)
	}

	// Check response body
	t.Log("Checking response body")
	if bodyCount == 0 && len(body) == 0 {
		return
	}
	if testCase.Expected.Body != nil {
		assert.Equal(t, testCase.Expected.Body, body)
	}
	if testCase.Expected.BodyFunc != nil {
		testCase.Expected.BodyFunc(t, body)
	}
	if testCase.Expected.BodyBytes != nil {
		assert.Equal(t, testCase.Expected.BodyBytes, body)
	}
}

// RunTestCases run a slice of test cases.
// testCases is the list of test cases to be executed.
func (s *Suite) RunTestCases(testCases []TestCase) {
	t := s.T()
	for i := range testCases {
		t.Log(testCases[i].Name)
		s.RunTestCase(&testCases[i])
	}
}

func (s *Suite) DefaultHostAPI() string {
	return "http://localhost"
}

// DefaultPublicBaseURL retrieve the public base endpoint URL.
// Return for the URL for the current configuration.
func (s *Suite) DefaultPublicBaseURL() string {
	// TODO Update this base URL
	return fmt.Sprintf("%s:%d%s", s.DefaultHostAPI(), s.cfg.Web.Port, s.cfg.Application.PathPrefix)
}

// DefaultPrivateBaseURL retrieve the private base endpoint URL.
// Return for the URL for the current configuration.
func (s *Suite) DefaultPrivateBaseURL() string {
	return fmt.Sprintf("%s:%d%s", s.DefaultHostAPI(), s.cfg.Web.Port, "/internal")
}

func (s *Suite) DefaultHealthcheckBaseURL() string {
	return fmt.Sprintf("%s:%d", s.DefaultHostAPI(), s.cfg.Web.Port)
}

func (s *Suite) DefaultMetricsURL() string {
	return fmt.Sprintf("%s:%d%s", s.DefaultHostAPI(), s.cfg.Metrics.Port, s.cfg.Metrics.Path)
}

// DoRequest execute a http request against a url using headers and the body specified.
// method is the HTTP method to use for the request.
// url is the to reach out.
// header represents the request headers.
// body is any golang type to be marshalled or a Reader interface (TODO future).
// Return the http.Response object and nil when the endpoint is reached out,
// or nil and a non error when some non API situation happens trying to reach
// out the endpoint.
func (s *Suite) DoRequest(method, url string, header http.Header, body any) (*http.Response, error) {
	var reader io.Reader = nil
	client := &http.Client{}

	if body != nil {
		// TODO add type assert so if a Reader interface
		// is provided, it will be used directly; this will
		// allow to wrong body to check for integration tests
		_body, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		if len(_body) > 0 {
			reader = bytes.NewReader(_body)
		}
	} else {
		reader = bytes.NewBufferString("")
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range header {
		req.Header.Set(key, strings.Join(value, "; "))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type BodyFunc func(t *testing.T, body []byte)

// TestCaseGiven represents the requirements for the smoke test to implement.
type TestCaseGiven struct {
	// Method represents the http method for the request.
	Method string
	// URL represents the url for the route to test.
	URL string
	// Header represents the set of header of the request.
	Header http.Header
	// Body represents a golang type to be marshalled before send the request;
	// this field exclude the BodyBytes field.
	Body any
	// BodyBytes represents a specific buffer for the request body; this
	// field exlude the Body field. This works for bad formed json documents,
	// and other scenarios where Body does not fit.
	BodyBytes []byte
}

// TestCaseExpect represents the expected response for a smoke test.
type TestCaseExpect struct {
	// StatusCode represents the http status code expected.
	StatusCode int
	// Header represents the expected http response headers.
	Header http.Header
	// Body represent an API type struct that after marshall should match the
	// returned response; this could be a situation, because the order of the
	// properties could not match. It is useful only when the property order
	// is deterministic, else use BodyFunc.
	Body any
	// BodyBytes represent a specific bytes returned on the expectations.
	BodyBytes []byte
	// BodyFunc represent a custom function that will return nil or error
	// to check some specifc body unserialized. This option exclude Body and
	// BodyBytes and is useful when we want expectations based on a
	// valid json document, but it is not a perfect fit of the Body.
	BodyFunc BodyFunc
}

// TestCase represents a test case for the smoke test
type TestCase struct {
	// Name represents a string to be printed out which will be displayed
	// in case of a failure happens.
	Name string
	// Given represents the given specification for the test case.
	Given TestCaseGiven
	// Expected represents the expected result for the operations.
	Expected TestCaseExpect
}

// StartSignalHandler set up the signal handler. This method MUST NOT
// be called several times, as that make no deterministic which signal
// handler will receive the call.
// c is the golang context associated, if it is nil a new background
// context is used.
// Return the cancel context generated that will called on exit and
// the cancel function associted to the context.
// See: https://pkg.go.dev/os/signal
func StartSignalHandler(c context.Context) (context.Context, context.CancelFunc) {
	if c == nil {
		c = context.Background()
	}
	ctx, cancel := context.WithCancel(c)
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-exit
		cancel()
	}()
	return ctx, cancel
}

// TearDownSignalHandler reset the signal handlers
func TearDownSignalHandler() {
	signal.Reset(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
}

func TestSuite(t *testing.T) {
	if value, exist := os.LookupEnv("TEST"); !exist || value != "integration" {
		t.Skip("This TestSuite needs to start infrastructure: Enforce by TEST=integration if it is up: make compose-up && make test-integration TEST=integration")
	}
	// TODO Add here your test suites
	suite.Run(t, new(SuiteTodosCreate))
	suite.Run(t, new(SuiteTodosRead))
	// suite.Run(t, new(SuiteTodosDelete))
	// suite.Run(t, new(SuiteTodosUpdate))
	// suite.Run(t, new(SuiteTodosList))
}
