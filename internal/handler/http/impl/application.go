package impl

import (
	"github.com/podengo-project/idmsvc-backend/internal/config"
	handler "github.com/podengo-project/idmsvc-backend/internal/handler/http"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository/client"
	metrics "github.com/podengo-project/idmsvc-backend/internal/metrics"
	usecase_interactor "github.com/podengo-project/idmsvc-backend/internal/usecase/interactor"
	"gorm.io/gorm"
)

type application struct {
	config        *config.Config
	metrics       *metrics.Metrics
	db            *gorm.DB
	todoPresenter presenter.Todo
}

func NewHandler(config *config.Config, db *gorm.DB, m *metrics.Metrics, inventory client.HostInventory) handler.Application {
	if config == nil {
		panic("config is nil")
	}
	if db == nil {
		panic("db is nil")
	}
	tr := NewTodoHandlertodoResource{
		interactor: usecase_interactor.NewTodo(),
	}

	// Instantiate application
	return &application{
		config:  config,
		db:      db,
		metrics: m,
		todo:    tr,
	}
}
