package test

import "github.com/avisiedo/go-microservice-1/internal/config"

// Config for testing
func GetTestConfig() (cfg *config.Config) {
	cfg = &config.Config{}
	config.Load(cfg)
	// override some default settings
	cfg.Application.PaginationDefaultLimit = 10
	cfg.Application.PaginationMaxLimit = 100
	cfg.Application.PathPrefix = "/api/todo"
	return cfg
}
