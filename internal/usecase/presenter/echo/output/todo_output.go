package output

import (
	"errors"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
)

type TodoOutput struct{}

func (o TodoOutput) createGuards(ctx echo.Context, data *model.Todo) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if data == nil {
		return errors.New("dataOutput is nil")
	}
	return nil
}

func (o TodoOutput) Create(ctx echo.Context, data *model.Todo) (*public.ToDo, error) {
	if err := o.createGuards(ctx, data); err != nil {
		return nil, err
	}
	output := &public.ToDo{}
	return output, nil
}

func (o TodoOutput) getAllGuards(ctx echo.Context, data []model.Todo) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if data == nil {
		return errors.New("dataOutput is nil")
	}
	return nil
}

func (o TodoOutput) GetAll(ctx echo.Context, data []model.Todo) ([]public.ToDo, error) {
	if err := o.getAllGuards(ctx, data); err != nil {
		return []public.ToDo{}, err
	}
	output := make([]public.ToDo, len(data))
	for i := range data {
		output[i].DueDate = &types.Date{Time: *data[i].DueDate}
		output[i].Title = *data[i].Title
		output[i].TodoId = &data[i].UUID
		output[i].Description = data[i].Description
	}
	return output, nil
}
