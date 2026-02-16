package config

import (
	"os"
	"testing"

	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	validator "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetDefaults(t *testing.T) {
	assert.Panics(t, func() {
		setDefaults(nil)
	})

	v := viper.New()
	assert.NotPanics(t, func() {
		setDefaults(v)
	})

	assert.Equal(t, DefaultWebPort, v.Get("web.port"))
	assert.Equal(t, "info", v.Get("logging.level"))
	assert.Equal(t, "http://localhost:8010/api/inventory/v1", v.Get("clients.host_inventory_base_url"))
	assert.Equal(t, DefaultTokenExpirationTimeSeconds, v.Get("app.token_expiration_seconds"))
	assert.Equal(t, PaginationDefaultLimit, v.Get("app.pagination_default_limit"))
	assert.Equal(t, PaginationMaxLimit, v.Get("app.pagination_max_limit"))
}

func TestLoad(t *testing.T) {
	var v *viper.Viper
	// 'cfg' is nil panic
	assert.PanicsWithError(t, common_err.ErrNil("cfg").Error(), func() {
		Load(nil)
	})

	// Success Load
	cfg := Config{}
	assert.NotPanics(t, func() {
		v = Load(&cfg)
	})
	assert.NotNil(t, v)
}

func TestValidateConfig(t *testing.T) {
	cfg := Config{}

	err := Validate(&cfg)
	assert.Error(t, err)

	cfg.Application.PathPrefix = DefaultPathPrefix
	ve, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)
	assert.Equal(t, len(ve), 1)
	assert.Equal(t, ve[0].Namespace(), "Config.Application.PathPrefix")
	assert.Equal(t, "required", ve[0].Tag())
}

func TestReportError(t *testing.T) {
	vSlice := validator.ValidationErrors{}
	assert.NotPanics(t, func() {
		reportError(vSlice)
	})
}

func TestReadEnv(t *testing.T) {
	require.NoError(t, os.Unsetenv("TEST"))
	assert.Equal(t, "default", readEnv("TEST", "default"))
	require.NoError(t, os.Setenv("TEST", "nodefault"))
	assert.Equal(t, "nodefault", readEnv("TEST", "default"))
}

func TestGet(t *testing.T) {
	var c *Config
	require.NotPanics(t, func() {
		c = Get()
	})
	require.NotNil(t, c)
}
