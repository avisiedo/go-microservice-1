package output

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/http/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
)

type TodoOutput struct{}

func (o TodoOutput) createGuards(ctx echo.Context, dataOutput *model.Todo, apiOutput *public.ToDo) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if dataOutput == nil {
		return errors.New("dataOutput is nil")
	}
	if apiOutput == nil {
		return errors.New("apiOutput is nil")
	}
	return nil
}

func (o TodoOutput) Create(ctx echo.Context, dataOutput *model.Todo, apiOutput *public.ToDo) error {
	var err error
	if err = o.createGuards(ctx, dataOutput, apiOutput); err != nil {
		return err
	}

	return nil
}
