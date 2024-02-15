package impl

import (
	"context"
	"sync"

	api_event "github.com/avisiedo/go-microservice-1/internal/api/event"
	"github.com/avisiedo/go-microservice-1/internal/config"
	handler_event "github.com/avisiedo/go-microservice-1/internal/handler/event"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/event"
	event_handler "github.com/avisiedo/go-microservice-1/internal/infrastructure/event/handler"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/service"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

type kafkaConsumer struct {
	context   context.Context
	cancel    context.CancelFunc
	waitGroup *sync.WaitGroup
	config    *config.Config

	db *gorm.DB
}

func NewKafkaConsumer(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, db *gorm.DB) service.ApplicationService {
	ctx, cancel := context.WithCancel(ctx)
	return &kafkaConsumer{
		context:   ctx,
		cancel:    cancel,
		waitGroup: wg,
		config:    cfg,

		db: db,
	}
}

func (s *kafkaConsumer) Start() error {
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()

		// Create event router
		eventRouter := event_handler.NewRouter()
		eventRouter.Add(api_event.TopicTodoCreated, handler_event.NewTodoCreatedEventHandler(s.db))

		// Start service
		event.Start(s.context, &s.config.Kafka, eventRouter)
		slog.Info("kafkaConsumer stopped")
	}()
	return nil
}

func (s *kafkaConsumer) Stop() error {
	s.cancel()
	return nil
}
