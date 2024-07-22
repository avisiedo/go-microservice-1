package echo

import (
	"github.com/labstack/echo/v4"
)

// Events represents the presenter for the /events resources
// The combination of all the presenters will match the
// API interface generated
type Events interface {
	// List All events
	// (GET /events)
	GetEvents(ctx echo.Context) error
	// Create a Events
	// (POST /events)
	CreateEvent(ctx echo.Context) error
	// Delete a Events
	// (DELETE /events/{eventsId})
	DeleteEvent(ctx echo.Context, eventsId string) error
	// Get a Events
	// (GET /events/{eventsId})
	GetEventByID(ctx echo.Context, eventsId string) error
	// Update a Events
	// (PUT /events/{eventsId})
	UpdateEvent(ctx echo.Context, eventsId string) error
}
