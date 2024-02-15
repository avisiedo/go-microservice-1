package model

import (
	"time"

	"github.com/avisiedo/go-microservice-1/internal/domain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Todo is the builder interface to create model.Todo components
type Todo interface {
	Build() *model.Todo
	WithModel(value gorm.Model) Todo
	WithTitle(value *string) Todo
	WithDescription(value *string) Todo
	WithDueDate(value *time.Time) Todo
	WithID(value uuid.UUID) Todo
}

type todo model.Todo

// NewTodo create builder instance to build Todo model components
func NewTodo() Todo {
	return &todo{}
}

func (b *todo) Build() *model.Todo {
	return (*model.Todo)(b)
}

func (b *todo) WithModel(value gorm.Model) Todo {
	b.Model = value
	return b
}

func (b *todo) WithDescription(value *string) Todo {
	b.Description = value
	return b
}

func (b *todo) WithDueDate(value *time.Time) Todo {
	b.DueDate = value
	return b
}

func (b *todo) WithTitle(value *string) Todo {
	b.Title = value
	return b
}

func (b *todo) WithID(value uuid.UUID) Todo {
	b.UUID = value
	return b
}
