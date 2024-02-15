package echo

import (
	"net/http"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	presenter "github.com/avisiedo/go-microservice-1/internal/interface/presenter/echo"
	"github.com/labstack/echo/v4"
)

type openAPIPresenter struct{}

func NewOpenAPI() presenter.OpenAPI {
	return &openAPIPresenter{}
}

var openapiSpec = public.PathToRawSpec("/openapi.json")

// GetOpenapi return the openapi specification as a json content
// from the boilerplate generated.
// ctx is the echo context.
// Return nil for success execution or an error object.
func (p openAPIPresenter) GetOpenapi(ctx echo.Context) error {
	resp, err := openapiSpec["/openapi.json"]()
	if err != nil {
		return err
	}
	return ctx.Blob(http.StatusOK, echo.MIMEApplicationJSON, resp)
}
