package http

import (
	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ToDo is the builder interface to create public.ToDo components
type ToDo interface {
	Build() *public.ToDo
	WithDescription(value *string) ToDo
	WithDueDate(value *openapi_types.Date) ToDo
	WithTitle(value string) ToDo
	WithID(value *openapi_types.UUID) ToDo
}

type todo public.ToDo

// NewToDo create builder instance to build ToDo components
func NewToDo() ToDo {
	return &todo{}
}

func (b *todo) Build() *public.ToDo {
	return (*public.ToDo)(b)
}

func (b *todo) WithDescription(value *string) ToDo {
	b.Description = value
	return b
}

func (b *todo) WithDueDate(value *openapi_types.Date) ToDo {
	b.DueDate = value
	return b
}

func (b *todo) WithTitle(value string) ToDo {
	b.Title = value
	return b
}

func (b *todo) WithID(value *openapi_types.UUID) ToDo {
	b.TodoId = value
	return b
}
