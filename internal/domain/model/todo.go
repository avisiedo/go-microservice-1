package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Todo is the busines struct container for the todos
type Todo struct {
	gorm.Model
	UUID        uuid.UUID
	Title       *string
	Description *string
	DueDate     *time.Time
}
