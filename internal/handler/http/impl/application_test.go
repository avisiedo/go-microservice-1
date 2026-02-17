package impl

import (
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/config"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	"github.com/avisiedo/go-microservice-1/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewHandler(t *testing.T) {
	sqlMock, gormDB, err := test.NewSqlMock(&gorm.Session{SkipHooks: true})
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, gormDB)
	assert.Panics(t, func() {
		NewHandler(nil, nil, nil)
	})
	assert.PanicsWithError(t, common_err.ErrNil("db").Error(), func() {
		NewHandler(&config.Config{}, nil, nil)
	})
	cfg := test.GetTestConfig()
	assert.PanicsWithError(t, common_err.ErrNil("m").Error(), func() {
		NewHandler(cfg, gormDB, nil)
	})
	assert.NotPanics(t, func() {
		NewHandler(cfg, gormDB, &metrics.Metrics{})
	})
}
