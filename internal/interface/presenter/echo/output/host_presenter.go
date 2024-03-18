package presenter

import (
	"github.com/podengo-project/idmsvc-backend/internal/api/http/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
)

type HostPresenter interface {
	HostConf(domain *model.Domain, token public.HostToken) (*public.HostConfResponse, error)
}
