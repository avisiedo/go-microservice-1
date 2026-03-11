package echo

import "github.com/labstack/echo/v5"

type Healthcheck interface {
	GetLivez(ctx echo.Context) error
	GetReadyz(ctx echo.Context) error
}
