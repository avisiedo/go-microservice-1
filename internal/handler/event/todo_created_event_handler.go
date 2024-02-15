package impl

import (
	"fmt"

	"github.com/avisiedo/go-microservice-1/internal/infrastructure/event"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gorm.io/gorm"
)

type todoCreatedEventHandler struct {
	db *gorm.DB
}

func NewTodoCreatedEventHandler(db *gorm.DB) event.Eventable {
	if db == nil {
		return nil
	}
	return &todoCreatedEventHandler{
		db: db,
	}
}

func (h *todoCreatedEventHandler) OnMessage(msg *kafka.Message) error {
	return fmt.Errorf("Not implemented")
}
