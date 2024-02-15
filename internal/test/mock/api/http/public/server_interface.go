// Code generated by mockery v2.38.0. DO NOT EDIT.

package public

import (
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// ServerInterface is an autogenerated mock type for the ServerInterface type
type ServerInterface struct {
	mock.Mock
}

// CreateTodo provides a mock function with given fields: ctx
func (_m *ServerInterface) CreateTodo(ctx echo.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for CreateTodo")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTodo provides a mock function with given fields: ctx, todoId
func (_m *ServerInterface) DeleteTodo(ctx echo.Context, todoId uuid.UUID) error {
	ret := _m.Called(ctx, todoId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteTodo")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, todoId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllTodos provides a mock function with given fields: ctx
func (_m *ServerInterface) GetAllTodos(ctx echo.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllTodos")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetTodo provides a mock function with given fields: ctx, todoId
func (_m *ServerInterface) GetTodo(ctx echo.Context, todoId uuid.UUID) error {
	ret := _m.Called(ctx, todoId)

	if len(ret) == 0 {
		panic("no return value specified for GetTodo")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, todoId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PatchTodo provides a mock function with given fields: ctx, todoId
func (_m *ServerInterface) PatchTodo(ctx echo.Context, todoId uuid.UUID) error {
	ret := _m.Called(ctx, todoId)

	if len(ret) == 0 {
		panic("no return value specified for PatchTodo")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, todoId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTodo provides a mock function with given fields: ctx, todoId
func (_m *ServerInterface) UpdateTodo(ctx echo.Context, todoId uuid.UUID) error {
	ret := _m.Called(ctx, todoId)

	if len(ret) == 0 {
		panic("no return value specified for UpdateTodo")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, todoId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewServerInterface creates a new instance of ServerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ServerInterface {
	mock := &ServerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
