package impl

import (
	"context"
	"sync"

	"github.com/avisiedo/go-microservice-1/internal/config"
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	handler_impl "github.com/avisiedo/go-microservice-1/internal/handler/http/impl"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/metrics"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/service"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

type svcApplication struct {
	Context   context.Context
	Cancel    context.CancelFunc
	Config    *config.Config
	WaitGroup *sync.WaitGroup
	Api       service.ApplicationService
	Kafka     service.ApplicationService
	Metrics   service.ApplicationService
	// AdditionalService service.ApplicationService
}

func NewApplication(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, db *gorm.DB) service.ApplicationService {
	if ctx == nil {
		panic(common_err.ErrNil("ctx"))
	}
	if wg == nil {
		panic(common_err.ErrNil("wg"))
	}
	if cfg == nil {
		panic(common_err.ErrNil("cfg"))
	}
	if db == nil {
		panic(common_err.ErrNil("db"))
	}

	s := &svcApplication{}
	s.Config = cfg
	s.Context, s.Cancel = context.WithCancel(ctx)
	s.WaitGroup = wg

	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg)

	// Create application handlers
	handler := handler_impl.NewHandler(s.Config, db, metrics)

	// Create Metrics service
	s.Metrics = NewMetrics(s.Context, s.WaitGroup, s.Config, handler)

	// Create Api service
	s.Api = NewApi(s.Context, s.WaitGroup, s.Config, handler, metrics)

	// Create kafka consumer service
	// TODO Uncomment or clean-up when we know if we use kafka
	// s.Kafka = NewKafkaConsumer(s.Context, s.WaitGroup, s.Config, db)

	return s
}

func (svc *svcApplication) Start() error {
	// svc.WaitGroup.Add(3)
	svc.WaitGroup.Add(2)
	go func() {
		defer svc.WaitGroup.Done()
		defer svc.Cancel()
		if err := svc.Api.Start(); err != nil {
			panic(err)
		}
		<-svc.Context.Done()
	}()

	// go func() {
	// 	defer svc.WaitGroup.Done()
	// 	defer svc.Cancel()
	// 	if err := svc.Kafka.Start(); err != nil {
	// 		panic(err)
	// 	}
	// 	<-svc.Context.Done()
	// }()

	go func() {
		defer svc.WaitGroup.Done()
		defer svc.Cancel()
		if err := svc.Metrics.Start(); err != nil {
			panic(err)
		}
		<-svc.Context.Done()
	}()
	return nil
}

func (svc *svcApplication) Stop() error {
	svc.Cancel()
	return nil
}
