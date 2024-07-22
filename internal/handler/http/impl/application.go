package impl

import (
	"github.com/avisiedo/go-microservice-1/internal/config"
	handler_http "github.com/avisiedo/go-microservice-1/internal/handler/http"
	metrics "github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	presenter_interface "github.com/avisiedo/go-microservice-1/internal/interface/presenter/echo"
	"github.com/avisiedo/go-microservice-1/internal/usecase/interactor"
	presenter "github.com/avisiedo/go-microservice-1/internal/usecase/presenter/echo"
	repository "github.com/avisiedo/go-microservice-1/internal/usecase/repository/db"
	"gorm.io/gorm"
)

type application struct {
	config  *config.Config
	db      *gorm.DB
	metrics *metrics.Metrics
	presenter_interface.Private
	presenter_interface.Todo
	presenter_interface.Events
	presenter_interface.OpenAPI
	presenter_interface.Healthcheck
	presenter_interface.Instrumentation
}

func NewHandler(cfg *config.Config, db *gorm.DB, m *metrics.Metrics) handler_http.Application {
	return newHandler(cfg, db, m)
}

func newHandler(cfg *config.Config, db *gorm.DB, m *metrics.Metrics) *application {
	if cfg == nil {
		panic("config is nil")
	}
	if db == nil {
		panic("db is nil")
	}
	if m == nil {
		panic("m is nil")
	}

	// Initialize the presenters
	todoPresenter := presenter.NewTodo(cfg, interactor.NewTodo(repository.NewTodo(cfg)), db)
	eventsPresenter := presenter.NewEvents(cfg)
	privatePresenter := presenter.NewPrivate(interactor.NewPrivate())
	openAPIPresenter := presenter.NewOpenAPI()
	instrumentationPresenter := presenter.NewInstrumentation(m)
	healthcheckPresenter := presenter.NewHealthcheck(interactor.NewHealthcheck(cfg))

	// Instantiate application
	return &application{
		config:          cfg,
		db:              db,
		metrics:         m,
		Private:         privatePresenter,
		Todo:            todoPresenter,
		Events:          eventsPresenter,
		OpenAPI:         openAPIPresenter,
		Instrumentation: instrumentationPresenter,
		Healthcheck:     healthcheckPresenter,
	}
}
