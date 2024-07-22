package echo

import (
	"net/http"

	"github.com/avisiedo/go-microservice-1/internal/config"
	presenter "github.com/avisiedo/go-microservice-1/internal/interface/presenter/echo"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type eventsPresenter struct {
	db *gorm.DB
	// input      input.EventsInput
	// output     output.EventsOutput
}

func NewEvents(cfg *config.Config) presenter.Events {
	return &eventsPresenter{
		db: nil,
	}
}

// List All events
// (GET /events)
func (p *eventsPresenter) GetEvents(ctx echo.Context) error {
	code := http.StatusNotImplemented
	return echo.NewHTTPError(code, http.StatusText(code))
}

// Create a Events
// (POST /events)
func (p *eventsPresenter) CreateEvent(ctx echo.Context) error {
	code := http.StatusNotImplemented
	return echo.NewHTTPError(code, http.StatusText(code))
}

// Delete a Events
// (DELETE /events/{eventsId})
func (p *eventsPresenter) DeleteEvent(ctx echo.Context, eventsId string) error {
	code := http.StatusNotImplemented
	return echo.NewHTTPError(code, http.StatusText(code))
}

// Get a Events
// (GET /events/{eventsId})
func (p *eventsPresenter) GetEventByID(ctx echo.Context, eventsId string) error {
	code := http.StatusNotImplemented
	return echo.NewHTTPError(code, http.StatusText(code))
}

// Update a Events
// (PUT /events/{eventsId})
func (p *eventsPresenter) UpdateEvent(ctx echo.Context, eventsId string) error {
	code := http.StatusNotImplemented
	return echo.NewHTTPError(code, http.StatusText(code))
}
