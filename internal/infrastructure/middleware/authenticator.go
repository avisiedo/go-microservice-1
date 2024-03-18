package middleware

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3filter"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type XRhIValidator interface {
	ValidateXRhIdentity(xrhi *identity.XRHID) error
}

func NewAuthenticator(v XRhIValidator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		// TODO Call here to your Authenticate function
		// return Authenticate(v, ctx, input)
		return nil
	}
}

func checkGuardsAuthenticate(v XRhIValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if v == nil {
		return internal_errors.NilArgError("v")
	}
	if ctx == nil {
		return internal_errors.NilArgError("ctx")
	}
	if input == nil {
		return internal_errors.NilArgError("input")
	}
	return nil
}
