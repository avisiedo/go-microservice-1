package echo

import (
	infra_metrics "github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	presenter "github.com/avisiedo/go-microservice-1/internal/interface/presenter/sync/echo"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type instrumentation struct {
	metrics *infra_metrics.Metrics
}

func NewInstrumentation(metrics *infra_metrics.Metrics) presenter.Instrumentation {
	if metrics == nil {
		panic("'metrics' is nil")
	}
	return &instrumentation{
		metrics: metrics,
	}
}

func (p *instrumentation) GetMetrics(ctx echo.Context) error {
	return echo.WrapHandler(promhttp.HandlerFor(
		p.metrics.Registry(),
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
			// Pass custom registry
			Registry: p.metrics.Registry(),
		},
	))(ctx)
}
