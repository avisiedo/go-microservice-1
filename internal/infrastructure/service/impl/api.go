package impl

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/avisiedo/go-microservice-1/internal/api/http/public"
	"github.com/avisiedo/go-microservice-1/internal/config"
	handler "github.com/avisiedo/go-microservice-1/internal/handler/http"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/router"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/service"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

type apiService struct {
	context   context.Context
	cancel    context.CancelFunc
	waitGroup *sync.WaitGroup
	config    *config.Config

	echo *echo.Echo
}

func NewApi(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, app handler.Application, metrics *metrics.Metrics) service.ApplicationService {
	if cfg == nil {
		panic("config is nil")
	}
	if wg == nil {
		panic("wg is nil")
	}

	result := &apiService{}
	result.context, result.cancel = context.WithCancel(ctx)
	result.waitGroup = wg
	result.config = cfg
	swagger, err := public.GetSwagger()
	if err != nil {
		panic(err.Error())
	}
	result.echo = router.NewRouterWithConfig(
		echo.New(),
		cfg,
		swagger,
		app,
		metrics,
	)
	result.echo.HideBanner = true
	if result.config.Logging.Level == "debug" || result.config.Logging.Level == "trace" {
		result.echo.Debug = true
	}
	if result.config.Logging.Level == "debug" || result.config.Logging.Level == "trace" {
		routes := result.echo.Routes()
		slog.Debug("Printing routes")
		for idx, route := range routes {
			slog.Debug("routing",
				slog.Int("index", idx),
				slog.Any("route", route),
			)
		}
	}

	return result
}

func (srv *apiService) Start() error {
	srv.waitGroup.Add(2)
	go func() {
		defer srv.waitGroup.Done()
		srvAddress := fmt.Sprintf(":%d", srv.config.Web.Port)
		slog.Debug("staring echo server", slog.String("srvAddress", srvAddress))
		if err := srv.echo.Start(srvAddress); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", slog.Any("error", err))
		}
	}()

	go func() {
		defer srv.waitGroup.Done()
		defer srv.cancel()
		<-srv.context.Done()
		slog.Info("Shutting down apiService")
		if err := srv.echo.Shutdown(context.Background()); err != nil {
			slog.Error("Failed to shut down apiService", slog.Any("error", err))
		}
	}()

	return nil
}

func (srv *apiService) Stop() error {
	srv.cancel()
	return nil
}
