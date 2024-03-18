package handler

import (
	"github.com/podengo-project/idmsvc-backend/internal/api/http/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/api/http/openapi"
	"github.com/podengo-project/idmsvc-backend/internal/api/http/private"
	"github.com/podengo-project/idmsvc-backend/internal/api/http/public"
)

type Application interface {
	public.ServerInterface
	private.ServerInterface
	metrics.ServerInterface
	openapi.ServerInterface
}
