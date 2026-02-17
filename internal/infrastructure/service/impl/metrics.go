package impl

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"log/slog"

	"github.com/avisiedo/go-microservice-1/internal/config"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	handler "github.com/avisiedo/go-microservice-1/internal/handler/http"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/router"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/service"
	"github.com/labstack/echo/v4"
)

type metricsService struct {
	context   context.Context
	cancel    context.CancelFunc
	waitGroup *sync.WaitGroup
	config    *config.Config

	echo *echo.Echo
}

func NewMetrics(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, app handler.Application) service.ApplicationService {
	if cfg == nil {
		panic(common_err.ErrNil("config"))
	}
	if wg == nil {
		panic(common_err.ErrNil("wg"))
	}

	result := &metricsService{}
	result.context, result.cancel = context.WithCancel(ctx)
	result.waitGroup = wg
	result.config = cfg

	result.echo = router.NewMetricsRouter(
		echo.New(),
		cfg,
		app,
	)
	result.echo.HideBanner = true

	if result.config.Logging.Level == "debug" {
		routes := result.echo.Routes()
		slog.Debug("Printing metrics routes")
		for idx, route := range routes {
			slog.Debug("routing",
				slog.Int("index", idx),
				slog.Any("route", route),
			)
		}
	}

	return result
}

func (srv *metricsService) Start() error {
	srv.waitGroup.Add(2)
	go func() {
		defer srv.waitGroup.Done()
		srvAddress := fmt.Sprintf(":%d", srv.config.Metrics.Port)
		slog.Debug("metrics", slog.String("srvAddress", srvAddress))
		if err := srv.echo.Start(srvAddress); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start metricsService", slog.Any("error", err))
		}
	}()

	go func() {
		defer srv.waitGroup.Done()
		defer srv.cancel()
		<-srv.context.Done()
		slog.Info("Shutting down metricsService")
		if err := srv.echo.Shutdown(context.Background()); err != nil {
			slog.Error(
				"error shuttingdown metricsService",
				slog.Any("error", err),
			)
		}
	}()

	return nil
}

func (srv *metricsService) Stop() error {
	srv.cancel()
	return nil
}
