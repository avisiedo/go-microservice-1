package echo

import (
	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	. "github.com/avisiedo/go-microservice-1/internal/interface/presenter/sync/echo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type todoOutput struct{}

func todoCommonGuards(ctx echo.Context, data *model.Todo) error {
	if ctx == nil {
		return common_err.ErrNil("ctx")
	}
	if data == nil {
		return common_err.ErrNil("data")
	}
	return nil
}

func NewTodoOutput() TodoOutput {
	return &todoOutput{}
}

func (o *todoOutput) createGuards(ctx echo.Context, data *model.Todo) error {
	return todoCommonGuards(ctx, data)
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
		return common_err.ErrNil("ctx")
	}
	if data == nil {
		return common_err.ErrNil("data")
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
	if err := todoCommonGuards(ctx, data); err != nil {
		return err
	}
	if (data.UUID == uuid.UUID{}) {
		return common_err.ErrEmpty("data.UUID")
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
