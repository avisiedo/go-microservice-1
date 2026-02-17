// The scope of this file is:
// - Define the configuration struct.
// - Set default configuration values.
// - Map the data so viper can load the configuration there.
// See: https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
// See: https://consoledot.pages.redhat.com/docs/dev/getting-started/migration/config.html
package config

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	validator "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"k8s.io/utils/env"
)

const (
	// DefaultAppName is used to compose the route paths
	DefaultAppName = "todo"
	// API URL path prefix
	DefaultPathPrefix = "/api/" + DefaultAppName + "/v1"
	// DefaultExpirationTime is used for the default token expiration period
	// expressed in seconds. The default value is set to 7200 (2 hours)
	DefaultTokenExpirationTimeSeconds = 7200
	// DefaultWebPort is the default port where the public API is listening
	DefaultWebPort = 8000

	// https://github.com/project-koku/koku/blob/main/koku/api/common/pagination.py

	// PaginationDefaultLimit is the default limit for the pagination
	PaginationDefaultLimit = 10
	// PaginationMaxLimit is the default max limit for the pagination
	PaginationMaxLimit = 1000

	// DefaultValidateAPI is true
	DefaultValidateAPI = true
)

type Config struct {
	Loaded      bool
	Web         Web
	Database    Database
	Logging     Logging
	Cloudwatch  Cloudwatch
	Metrics     Metrics
	Clients     Clients
	Application Application `mapstructure:"app"`
}

type Web struct {
	Port int16
}

type Database struct {
	Host     string
	Port     int
	User     string `json:"-"`
	Password string `json:"-"`
	Name     string
	// https://stackoverflow.com/questions/54844546/how-to-unmarshal-golang-viper-snake-case-values
	CACertPath string `mapstructure:"ca_cert_path"`
}

type Logging struct {
	Level    string
	Console  bool
	Location bool
}

type Cloudwatch struct {
	Region  string
	Key     string `json:"-"`
	Secret  string `json:"-"`
	Session string
	Group   string
	Stream  string
}

type Metrics struct {
	// Defines the path to the metrics server that the app should be configured to
	// listen on for metric traffic.
	Path string `mapstructure:"path"`

	// Defines the metrics port that the app should be configured to listen on for
	// metric traffic.
	Port int `mapstructure:"port"`
}

// Clients gather all the client settings for the
type Clients struct {
	Inventory InventoryClient
}

type InventoryClient struct {
	// Define the base url for the host inventory service
	BaseUrl string `mapstructure:"base_url"`
}

// Application hold specific application settings
type Application struct {
	// PathPrefix is the API URL's path prefix, e.g. /api/todo
	PathPrefix string `mapstructure:"url_path_prefix" validate:"required"`
	// Indicate the default pagination limit when it is 0 or not filled
	PaginationDefaultLimit int `mapstructure:"pagination_default_limit"`
	// Indicate the max pagination limit when it is grather
	PaginationMaxLimit int `mapstructure:"pagination_max_limit"`
	// ValidateAPI enable the API validation for every request
	ValidateAPI bool `mapstructure:"validate_api"`
}

var (
	config *Config = nil
	once   sync.Once
)

func setDefaults(v *viper.Viper) {
	if v == nil {
		panic("viper instance cannot be nil")
	}
	// Web
	v.SetDefault("web.port", DefaultWebPort)

	// Database
	v.SetDefault("database.host", "")
	v.SetDefault("database.port", "")
	v.SetDefault("database.user", "")
	v.SetDefault("database.password", "")
	v.SetDefault("database.name", "")
	v.SetDefault("database.ca_cert_path", "")

	// Cloudwatch

	// Miscelanea
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.console", true)
	v.SetDefault("logging.location", false)

	// Clients
	v.SetDefault("clients.host_inventory_base_url", "http://localhost:8010/api/inventory/v1")

	// Application specific

	// Set default value for application expiration time for
	// the token created by the RHEL IDM domains
	v.SetDefault("app.token_expiration_seconds", DefaultTokenExpirationTimeSeconds)
	v.SetDefault("app.pagination_default_limit", PaginationDefaultLimit)
	v.SetDefault("app.pagination_max_limit", PaginationMaxLimit)
	v.SetDefault("app.validate_api", DefaultValidateAPI)
	v.SetDefault("app.url_path_prefix", DefaultPathPrefix)
	v.SetDefault("app.secret", "")
	v.SetDefault("app.debug", false)
}

func Load(cfg *Config) *viper.Viper {
	var err error

	if cfg == nil {
		panic(common_err.ErrNil("cfg"))
	}

	v := viper.New()
	v.AddConfigPath(env.GetString("CONFIG_PATH", "./configs"))
	v.SetConfigName("config.yaml")
	v.SetConfigType("yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	if err = v.ReadInConfig(); err != nil {
		slog.Warn("Not using config.yaml", slog.Any("error", err))
	}
	if err = v.Unmarshal(cfg); err != nil {
		slog.Warn("Mapping to configuration", slog.Any("error", err))
	}

	return v
}

func reportError(err error) {
	for _, err := range err.(validator.ValidationErrors) {
		slog.Error(
			"Configuration validation error",
			slog.String("namespace", err.Namespace()),
			slog.Group("rule",
				slog.String("tag", err.Tag()),
				slog.Any("value", err.Value),
			),
			slog.String("got", err.Param()),
			slog.String("type", err.Kind().String()),
		)
	}
}

func Validate(cfg *Config) (err error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(cfg)
}

// Get is a singleton to get the global loaded configuration.
func Get() *Config {
	once.Do(func() {
		config = &Config{}
		_ = Load(config)
		if err := Validate(config); err != nil {
			reportError(err)
			panic("Invalid configuration")
		}
	})

	return config
}

func readEnv(key string, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		value = def
	}
	return value
}
