// Code generated by mockery v2.38.0. DO NOT EDIT.

package echo

import (
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"
)

// Healthcheck is an autogenerated mock type for the Healthcheck type
type Healthcheck struct {
	mock.Mock
}

// GetLivez provides a mock function with given fields: ctx
func (_m *Healthcheck) GetLivez(ctx echo.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetLivez")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetReadyz provides a mock function with given fields: ctx
func (_m *Healthcheck) GetReadyz(ctx echo.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetReadyz")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewHealthcheck creates a new instance of Healthcheck. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHealthcheck(t interface {
	mock.TestingT
	Cleanup(func())
}) *Healthcheck {
	mock := &Healthcheck{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}