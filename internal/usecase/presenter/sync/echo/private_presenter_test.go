package echo

import (
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/test/mock/interface/interactor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrivate(t *testing.T) {
	assert.PanicsWithValue(t, "interactor is nil", func() {
		_ = NewPrivate(nil)
	})

	i := interactor.NewPrivate(t)
	p := NewPrivate(i)
	require.NotNil(t, p)
}
