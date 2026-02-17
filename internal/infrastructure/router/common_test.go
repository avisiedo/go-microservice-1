package router

import (
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/config"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func helperNewConfig() *config.Config {
	cfg := &config.Config{}
	_ = config.Load(cfg)
	return cfg
}

func helperNewEchoRouteConfig(t *testing.T) (*echo.Echo, *config.Config) {
	e := echo.New()
	cfg := helperNewConfig()
	return e, cfg
}

func TestNewRouterWithConfigCommonGuards(t *testing.T) {
	e, cfg := helperNewEchoRouteConfig(t)

	assert.PanicsWithError(t, common_err.ErrNil("e").Error(), func() {
		newRouterWithConfigCommonGuards(nil, nil)
	})
	assert.PanicsWithError(t, common_err.ErrNil("cfg").Error(), func() {
		newRouterWithConfigCommonGuards(e, nil)
	})
	assert.NotPanics(t, func() {
		newRouterWithConfigCommonGuards(e, cfg)
	})
}
