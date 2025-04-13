package echo

import (
	"errors"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	. "github.com/avisiedo/go-microservice-1/internal/interface/presenter/sync/echo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type todoOutput struct{}

func NewTodoOutput() TodoOutput {
	return &todoOutput{}
}

func (o *todoOutput) createGuards(ctx echo.Context, data *model.Todo) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if data == nil {
		return errors.New("dataOutput is nil")
	}
	return nil
}

func (o *todoOutput) Create(ctx echo.Context, data *model.Todo) (*public.ToDo, error) {
	if err := o.createGuards(ctx, data); err != nil {
		return nil, err
	}
	output := &public.ToDo{}
	uuidTemp := &uuid.UUID{}
	*uuidTemp = data.UUID
	output.TodoId = uuidTemp
	output.Title = data.Title
	output.Description = data.Description
	output.DueDate = data.DueDate
	return output, nil
}

func (o *todoOutput) getAllGuards(ctx echo.Context, data []model.Todo) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if data == nil {
		return errors.New("dataOutput is nil")
	}
	return nil
}

func (o *todoOutput) GetAll(ctx echo.Context, data []model.Todo) ([]public.ToDo, error) {
	if err := o.getAllGuards(ctx, data); err != nil {
		return nil, err
	}
	output := make([]public.ToDo, len(data))
	for i := range data {
		output[i].TodoId = &data[i].UUID
		output[i].Title = data[i].Title
		output[i].Description = data[i].Description
		output[i].DueDate = data[i].DueDate
	}
	return output, nil
}

func (o *todoOutput) getGuards(ctx echo.Context, data *model.Todo) error {
	if ctx == nil {
		return errors.New("'ctx' is nil")
	}
	if data == nil {
		return errors.New("'data' is nil")
	}
	if (data.UUID == uuid.UUID{}) {
		return errors.New("'data.UUID' is empty")
	}
	return nil
}

func (o *todoOutput) Get(ctx echo.Context, data *model.Todo) (*public.ToDo, error) {
	if err := o.getGuards(ctx, data); err != nil {
		return nil, err
	}
	out := &public.ToDo{}
	outUUID := &uuid.UUID{}
	*outUUID = data.UUID
	out.TodoId = outUUID
	out.Title = data.Title
	out.Description = data.Description
	out.DueDate = data.DueDate
	return out, nil
}
