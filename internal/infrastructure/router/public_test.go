package router

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	"github.com/avisiedo/go-microservice-1/internal/test"
	openapi "github.com/avisiedo/go-microservice-1/internal/test/mock/api/http/openapi"
	presenter "github.com/avisiedo/go-microservice-1/internal/test/mock/interface/presenter/echo"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func helperNewContextForSkipper(route string, method string, path string, headers map[string]string) echo.Context {
	// See: https://echo.labstack.com/guide/testing/
	e := echo.New()
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(route)
	return c
}

func helperNewGroupPublic(t *testing.T) (*echo.Echo, *config.Config, public.ServerInterface, *openapi.ServerInterface, *metrics.Metrics) {
	var (
		err error
		db  *gorm.DB
	)

	reg := prometheus.NewRegistry()
	require.NotNil(t, reg)
	m := metrics.NewMetrics(reg)
	require.NotNil(t, m)

	cfg := test.GetTestConfig()
	_, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, db)

	presenterPublic := presenter.NewTodo(t)
	presenterOpenAPI := openapi.NewServerInterface(t)

	e := echo.New()
	require.NotNil(t, e)
	return e, cfg, presenterPublic, presenterOpenAPI, m
}

func TestNewGroupPublicPanics(t *testing.T) {
	assert.Panics(t, func() {
		_, cfg, presenterPublic, presenterOpenAPI, m := helperNewGroupPublic(t)
		newGroupPublic(nil, cfg, presenterPublic, presenterOpenAPI, m)
	})
	assert.Panics(t, func() {
		e, cfg, presenterPublic, _, _ := helperNewGroupPublic(t)
		newGroupPublic(e.Group(cfg.Application.PathPrefix), cfg, presenterPublic, nil, nil)
	})
	assert.Panics(t, func() {
		e, cfg, presenterPublic, presenterOpenAPI, _ := helperNewGroupPublic(t)
		newGroupPublic(e.Group(cfg.Application.PathPrefix), cfg, presenterPublic, presenterOpenAPI, nil)
	})
	assert.NotPanics(t, func() {
		e, cfg, presenterPublic, presenterOpenAPI, m := helperNewGroupPublic(t)
		newGroupPublic(e.Group(cfg.Application.PathPrefix), cfg, presenterPublic, presenterOpenAPI, m)
	})
}

func TestNewGroupPublic(t *testing.T) {
	const (
		appPrefix = "/api"
		appName   = "/todo"
	)
	type TestCaseExpected map[string]map[string]string
	swagger, err := public.GetSwagger()
	require.NoError(t, err)
	require.NotNil(t, swagger)

	majorVersion := strings.Split(swagger.Info.Version, ".")[0]
	majorMinorVersion := majorVersion + "." + strings.Split(swagger.Info.Version, ".")[1]

	testCases := TestCaseExpected{
		"/private/readyz": {
			"GET": "github.com/avisiedo/go-microservice-1/internal/api/http/handler.ping",
		},
		"/private/livez": {
			"GET": "github.com/avisiedo/go-microservice-1/internal/api/http/public.ping",
		},

		appPrefix + appName + "/v" + majorMinorVersion + "/openapi.json": {
			"GET": "github.com/avisiedo/go-microservice-1/internal/api/http/openapi.(*ServerInterfaceWrapper).GetOpenapi-fm",
		},
		appPrefix + appName + "/v" + majorVersion + "/openapi.json": {
			"GET": "github.com/avisiedo/go-microservice-1/internal/api/http/openapi.(*ServerInterfaceWrapper).GetOpenapi-fm",
		},

		appPrefix + appName + "/v" + majorMinorVersion + "/todos": {
			"GET":  "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).GetAllTodos-fm",
			"POST": "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).CreateTodo-fm",
		},
		appPrefix + appName + "/v" + majorVersion + "/todos": {
			"GET":  "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).GetAllTodos-fm",
			"POST": "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).CreateTodo-fm",
		},

		appPrefix + appName + "/v" + majorMinorVersion + "/todos/:todoId": {
			"GET":    "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).GetTodo-fm",
			"PUT":    "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).UpdateTodo-fm",
			"PATCH":  "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).PatchTodo-fm",
			"DELETE": "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).DeleteTodo-fm",
		},
		appPrefix + appName + "/v" + majorVersion + "/todos/:todoId": {
			"GET":    "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).GetTodo-fm",
			"PUT":    "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).UpdateTodo-fm",
			"PATCH":  "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).PatchTodo-fm",
			"DELETE": "github.com/avisiedo/go-microservice-1/internal/api/http/public.(*ServerInterfaceWrapper).DeleteTodo-fm",
		},

		// This routes are added when the group is created
		appPrefix + appName + "/v" + majorMinorVersion + "/*": {
			"echo_route_not_found": "github.com/labstack/echo/v4.init.func1",
		},
		appPrefix + appName + "/v" + majorVersion + "/*": {
			"echo_route_not_found": "github.com/labstack/echo/v4.init.func1",
		},
		appPrefix + appName + "/v" + majorMinorVersion: {
			"echo_route_not_found": "github.com/labstack/echo/v4.init.func1",
		},
		appPrefix + appName + "/v" + majorVersion: {
			"echo_route_not_found": "github.com/labstack/echo/v4.init.func1",
		},
	}

	e, _, _, _, _ := helperNewGroupPublic(t)

	// Match Routes in expected
	for _, route := range e.Routes() {
		t.Logf("Method=%s Path=%s Name=%s", route.Method, route.Path, route.Name)

		methods, okPath := testCases[route.Path]
		if !okPath {
			t.Logf("Path=%s not found", route.Path)
		}
		assert.Truef(t, okPath, "path=%s not found into the expected ones", route.Path)

		name, okMethod := methods[route.Method]
		if !okMethod {
			t.Logf("Method=%s not found for path=%s", route.Method, route.Path)
		}
		assert.Truef(t, okMethod, "method=%s not found into the expected ones for the path=%s", route.Method, route.Path)
		assert.Equalf(t, name, route.Name, "handler for path=%s method=%s does not match", route.Path, route.Method)
	}

	e, cfg, publicPresenter, openAPIPresenter, m := helperNewGroupPublic(t)

	prefix := cfg.Application.PathPrefix + "/v1"
	group := newGroupPublic(e.Group(prefix), cfg, publicPresenter, openAPIPresenter, m)
	require.NotNil(t, group)

	for _, route := range e.Routes() {
		t.Logf("Method=%s Path=%s Name=%s", route.Method, route.Path, route.Name)

		methods, okPath := testCases[route.Path]
		assert.Truef(t, okPath, "path=%s not found into the expected ones", route.Path)

		name, okMethod := methods[route.Method]
		assert.Truef(t,
			okMethod,
			"method=%s not found into the expected ones for the path=%s",
			route.Method,
			route.Path)
		assert.Equalf(t,
			name,
			route.Name,
			"handler for path=%s method=%s does not match",
			route.Path,
			route.Method)
	}
}

func TestGetOpenapiPaths(t *testing.T) {
	assert.PanicsWithValue(t, "'cfg' is nil", func() {
		getOpenapiPaths(nil, nil)
	})

	cfg := &config.Config{}
	_ = config.Load(cfg)
	assert.PanicsWithValue(t, "'swagger' is nil", func() {
		getOpenapiPaths(cfg, nil)
	})

	swagger, err := public.GetSwagger()
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
	cfg := &config.Config{}
	_ = config.Load(cfg)

	skipper := newSkipperOpenapi(cfg)
	assert.NotNil(t, skipper)

	path := fmt.Sprintf("%s/openapi.json", cfg.Application.PathPrefix+"/v1")
	ctx := helperNewContextForSkipper(path, echo.GET, path, map[string]string{})
	assert.True(t, skipper(ctx))
}
