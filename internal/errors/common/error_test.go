package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckEmptyFieldName(t *testing.T) {
	require.PanicsWithError(t, ErrEmpty("fieldName").Error(), func() {
		checkEmptyFieldName("")
	})
	require.NotPanics(t, func() {
		checkEmptyFieldName("demo")
	})
}

func TestErrNil(t *testing.T) {
	assert.PanicsWithError(t, ErrEmpty("fieldName").Error(), func() {
		_ = ErrNil("")
	})
	assert.EqualError(t, ErrNil("someField"), "'someField' is nil")
}

func TestErrEmpty(t *testing.T) {
	assert.PanicsWithError(t, ErrEmpty("fieldName").Error(), func() {
		_ = ErrEmpty("")
	})
	require.EqualError(t, ErrEmpty("someField"), "'someField' is empty")
}
