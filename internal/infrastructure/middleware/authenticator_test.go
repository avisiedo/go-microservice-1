package middleware

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock_middleware "github.com/podengo-project/idmsvc-backend/internal/test/mock/infrastructure/middleware"
)

func TestCheckGuardsAuthenticate(t *testing.T) {
	// v is nil
	err := checkGuardsAuthenticate(nil, nil, nil)
	assert.EqualError(t, err, "code=500, message='v' cannot be nil")

	// ctx is nil
	v := mock_middleware.NewXRhIValidator(t)
	err = checkGuardsAuthenticate(v, nil, nil)
	assert.EqualError(t, err, "code=500, message='ctx' cannot be nil")

	// input is nil
	ctx := context.Background()
	err = checkGuardsAuthenticate(v, ctx, nil)
	assert.EqualError(t, err, "code=500, message='input' cannot be nil")

	// Success case
	input := openapi3filter.AuthenticationInput{}
	err = checkGuardsAuthenticate(v, ctx, &input)
	assert.NoError(t, err)
}

func TestNewAuthenticator(t *testing.T) {
	var a openapi3filter.AuthenticationFunc

	// Not panics
	assert.NotPanics(t, func() {
		a = NewAuthenticator(nil)
	})
	require.NotNil(t, a)

	// Not Panics
	v := mock_middleware.NewXRhIValidator(t)
	assert.NotPanics(t, func() {
		a = NewAuthenticator(v)
	})

	// Wrong security schema name
	err := a(nil, nil)
	assert.EqualError(t, err, "code=500, message='ctx' cannot be nil")
}
