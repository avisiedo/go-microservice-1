package impl

import (
	"github.com/avisiedo/go-microservice-1/internal/config"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	handler_http "github.com/avisiedo/go-microservice-1/internal/handler/http"
	metrics "github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	presenter_interface "github.com/avisiedo/go-microservice-1/internal/interface/presenter/sync/echo"
	"github.com/avisiedo/go-microservice-1/internal/usecase/interactor"
	presenter "github.com/avisiedo/go-microservice-1/internal/usecase/presenter/sync/echo"
	repository "github.com/avisiedo/go-microservice-1/internal/usecase/repository/db"
	"gorm.io/gorm"
)

type application struct {
	config  *config.Config
	db      *gorm.DB
	metrics *metrics.Metrics
	presenter_interface.Private
	presenter_interface.Todo
	presenter_interface.OpenAPI
	presenter_interface.Healthcheck
	presenter_interface.Instrumentation
}

func NewHandler(cfg *config.Config, db *gorm.DB, m *metrics.Metrics) handler_http.Application {
	return newHandler(cfg, db, m)
}

func newHandler(cfg *config.Config, db *gorm.DB, m *metrics.Metrics) *application {
	if cfg == nil {
		panic(common_err.ErrNil("config"))
	}
	if db == nil {
		panic(common_err.ErrNil("db"))
	}
	if m == nil {
		panic(common_err.ErrNil("m"))
	}

	// Initialize the presenters
	todoPresenter := presenter.NewTodo(cfg, interactor.NewTodo(repository.NewTodo(cfg)), db)
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
		OpenAPI:         openAPIPresenter,
		Instrumentation: instrumentationPresenter,
		Healthcheck:     healthcheckPresenter,
	}
}
