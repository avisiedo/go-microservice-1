package echo

import (
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Todo represents the presenter for the /todo resources
// The combination of all the presenters will match the
// API interface generated
type Todo interface {
	// Retrieve all ToDo items
	// (GET /todos)
	GetAllTodos(ctx echo.Context) error
	// Create a new ToDo item
	// (POST /todos)
	CreateTodo(ctx echo.Context) error
	// Remove item by ID
	// (DELETE /todos/{todoId})
	DeleteTodo(ctx echo.Context, todoId openapi_types.UUID) error
	// Get item by ID
	// (GET /todos/{todoId})
	GetTodo(ctx echo.Context, todoId openapi_types.UUID) error
	// Patch an existing ToDo item
	// (PATCH /todos/{todoId})
	PatchTodo(ctx echo.Context, todoId openapi_types.UUID) error
	// Substitute an existing ToDo item
	// (PUT /todos/{todoId})
	UpdateTodo(ctx echo.Context, todoId openapi_types.UUID) error
}
