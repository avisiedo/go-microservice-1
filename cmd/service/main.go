package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/datastore"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/logger"
	impl_service "github.com/avisiedo/go-microservice-1/internal/infrastructure/service/impl"
)

func startSignalHandler(c context.Context) (context.Context, context.CancelFunc) {
	if c == nil {
		c = context.Background()
	}
	ctx, cancel := context.WithCancel(c)
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-exit
		cancel()
	}()
	return ctx, cancel
}

func main() {
	wg := &sync.WaitGroup{}
	logger.LogBuildInfo("idmscv-backend")
	cfg := config.Get()
	logger.InitLogger(cfg)
	db := datastore.NewDB(cfg)
	defer datastore.Close(db)

	ctx, cancel := startSignalHandler(context.Background())
	s := impl_service.NewApplication(ctx, wg, cfg, db)
	if e := s.Start(); e != nil {
		panic(e)
	}
	<-ctx.Done()
	defer cancel()
	s.Stop()
	wg.Wait()
}
