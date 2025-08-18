package router

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/config"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	"github.com/avisiedo/go-microservice-1/internal/test"
	mock_openapi "github.com/avisiedo/go-microservice-1/internal/test/mock/api/http/openapi"
	mock_http "github.com/avisiedo/go-microservice-1/internal/test/mock/handler/http"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func helperNewContextForSkipper(method, path string, headers map[string]string) echo.Context {
	// See: https://echo.labstack.com/guide/testing/
	e := echo.New()
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)
	return c
}

func helperNewGroupPublic(t *testing.T) (*echo.Echo, *config.Config, *mock_http.Application, *mock_openapi.ServerInterface, *metrics.Metrics) {
	var (
		err error
		db  *gorm.DB
	)

	e, cfg := helperNewEchoRouteConfig(t)

	reg := prometheus.NewRegistry()
	require.NotNil(t, reg)
	m := metrics.NewMetrics(reg)
	require.NotNil(t, m)

	_, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, db)

	presenterPublic := mock_http.NewApplication(t)
	presenterOpenAPI := mock_openapi.NewServerInterface(t)

	public.RegisterHandlers(e.Group(cfg.Application.PathPrefix), presenterPublic)
	return e, cfg, presenterPublic, presenterOpenAPI, m
}

func helperOpenApiURI(cfg *config.Config) string {
	return cfg.Application.PathPrefix + "/openapi.json"
}

func TestNewGroupPublicPanics(t *testing.T) {
	e, cfg, presenterPublic, presenterOpenAPI, m := helperNewGroupPublic(t)
	assert.PanicsWithError(t, common_err.ErrNil("e").Error(), func() {
		newPublic(nil, nil, nil, nil, nil)
	})
	assert.PanicsWithError(t, common_err.ErrNil("cfg").Error(), func() {
		newPublic(e.Group(cfg.Application.PathPrefix), nil, nil, nil, nil)
	})
	assert.PanicsWithError(t, common_err.ErrNil("publicHandler").Error(), func() {
		newPublic(e.Group(cfg.Application.PathPrefix), cfg, nil, nil, nil)
	})
	assert.PanicsWithError(t, common_err.ErrNil("openapiHandler").Error(), func() {
		newPublic(e.Group(cfg.Application.PathPrefix), cfg, presenterPublic, nil, nil)
	})
	assert.PanicsWithError(t, common_err.ErrNil("metricsHandler").Error(), func() {
		newPublic(e.Group(cfg.Application.PathPrefix), cfg, presenterPublic, presenterOpenAPI, nil)
	})
	assert.NotPanics(t, func() {
		newPublic(e.Group(cfg.Application.PathPrefix), cfg, presenterPublic, presenterOpenAPI, m)
	})
}

func TestNewGroupPublicGroupRegistered(t *testing.T) {
	const (
		appPrefix = "/api/todo"
	)
	type TestCaseGiven struct {
		HandlerName string
		Params      []string
	}
	type TestCase struct {
		Given    TestCaseGiven
		Expected string
	}

	swagger, err := public.GetSwagger()
	require.NoError(t, err)
	require.NotNil(t, swagger)

	majorVersion := strings.Split(swagger.Info.Version, ".")[0]
	// majorMinorVersion := majorVersion + "." + strings.Split(swagger.Info.Version, ".")[1]

	testCases := []TestCase{
		// openapi.json
		{
			Given: TestCaseGiven{
				HandlerName: "github.com/avisiedo/go-microservice-1/internal/api/http/openapi.(*ServerInterfaceWrapper).GetOpenapi-fm",
				Params:      []string{},
			},
			Expected: appPrefix + "/v" + majorVersion + "/openapi.json",
		},

		// public api
		{
			Given: TestCaseGiven{
				HandlerName: "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).GetAllTodos-fm",
				Params:      []string{},
			},
			Expected: appPrefix + "/v" + majorVersion + "/todos",
		},
		{
			Given: TestCaseGiven{
				HandlerName: "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).CreateTodo-fm",
				Params:      []string{},
			},
			Expected: appPrefix + "/v" + majorVersion + "/todos",
		},
	}

	// Check panic
	assert.PanicsWithError(t, common_err.ErrNil("e").Error(), func() {
		newPublic(nil, nil, nil, nil, nil)
	})

	// Check the group generated
	for _, testCase := range testCases {
		t.Logf("HandlerName=%s", testCase.Given.HandlerName)
		e, cfg, publicPresenter, openAPIPresenter, m := helperNewGroupPublic(t)

		prefix := cfg.Application.PathPrefix
		require.NotNil(t, newPublic(e.Group(prefix), cfg, publicPresenter, openAPIPresenter, m))

		result := e.Reverse(testCase.Given.HandlerName, testCase.Given.Params)
		require.Equal(t, testCase.Expected, result)
		publicPresenter.AssertExpectations(t)
		openAPIPresenter.AssertExpectations(t)
	}
}

func TestGetOpenapiPaths(t *testing.T) {
	// Check cfg is nil
	assert.PanicsWithError(t, common_err.ErrNil("cfg").Error(), func() {
		getOpenapiPaths(nil, nil)
	})

	// Check swagger is nil
	cfg := helperNewConfig()
	assert.PanicsWithError(t, common_err.ErrNil("swagger").Error(), func() {
		getOpenapiPaths(cfg, nil)
	})

	// Check swagger.Info.Version is empty
	swagger, err := public.GetSwagger()
	swagger.Info.Version = ""
	require.NoError(t, err)
	require.NotNil(t, swagger)
	assert.PanicsWithError(t, common_err.ErrEmpty("swagger.Info.Version").Error(), func() {
		getOpenapiPaths(cfg, swagger)
	})

	// Success use case
	swagger, err = public.GetSwagger()
	require.NoError(t, err)
	require.NotNil(t, swagger)
	cachedPaths := getOpenapiPaths(cfg, swagger)
	assert.NotNil(t, cachedPaths)
	assert.Equal(t,
		[]string{
			"/api/todo/v1/openapi.json",
			"/api/todo/v1.0/openapi.json",
		},
		cachedPaths(),
	)
}

func TestNewSkipperOpenapi(t *testing.T) {
	cfg := helperNewConfig()

	skipper := newSkipperOpenapi(cfg)
	assert.NotNil(t, skipper)

	path := helperOpenApiURI(cfg)
	ctx := helperNewContextForSkipper(echo.GET, path, map[string]string{})
	assert.True(t, skipper(ctx))

	ctx = helperNewContextForSkipper(echo.GET, "/something/does/not/match", map[string]string{})
	assert.False(t, skipper(ctx))
}
