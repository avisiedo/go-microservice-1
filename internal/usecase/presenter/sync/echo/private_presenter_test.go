package echo

import (
	"testing"

	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	"github.com/avisiedo/go-microservice-1/internal/test/mock/interface/interactor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrivate(t *testing.T) {
	assert.PanicsWithError(t, common_err.ErrNil("i").Error(), func() {
		_ = NewPrivate(nil)
	})

	i := interactor.NewPrivate(t)
	p := NewPrivate(i)
	require.NotNil(t, p)
}
