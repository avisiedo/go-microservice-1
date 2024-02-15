package interactor

import (
	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/interface/interactor"
)

type healtchcheckInteractor struct {
	cfg *config.Config
}

func NewHealthcheck(cfg *config.Config) interactor.HealthcheckInteractor {
	return &healtchcheckInteractor{
		cfg: cfg,
	}
}

func (i *healtchcheckInteractor) IsLive() error {
	// TODO implement checks here
	// IsLive means the process is up and running
	return nil
}

func (i *healtchcheckInteractor) IsReady() error {
	// TODO implement checks here
	// IsReady means IsLive and the 3rd parties are too,
	// so the input requests can be attended
	// - database
	// - s3
	// - redis
	return nil
}
