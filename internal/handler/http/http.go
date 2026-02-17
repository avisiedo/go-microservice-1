package http

import (
	"github.com/avisiedo/go-microservice-1/internal/api/http/healthcheck"
	"github.com/avisiedo/go-microservice-1/internal/api/http/metrics"
	"github.com/avisiedo/go-microservice-1/internal/api/http/openapi"
	"github.com/avisiedo/go-microservice-1/internal/api/http/private"
	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
)

type Application interface {
	healthcheck.ServerInterface
	public.ServerInterface
	private.ServerInterface
	metrics.ServerInterface
	openapi.ServerInterface
}
