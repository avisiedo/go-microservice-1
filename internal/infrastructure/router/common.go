package router

import (
	"github.com/avisiedo/go-microservice-1/internal/config"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	"github.com/labstack/echo/v4"
)

func newRouterWithConfigCommonGuards(e *echo.Echo, cfg *config.Config) {
	if e == nil {
		panic(common_err.ErrNil("e"))
	}
	if cfg == nil {
		panic(common_err.ErrNil("cfg"))
	}
}
