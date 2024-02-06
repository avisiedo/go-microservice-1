package adapter

import (
	"fmt"

	"github.com/avisiedo/go-microservice-1/internal/api/event"
	"github.com/avisiedo/go-microservice-1/internal/api/header"
	"github.com/avisiedo/go-microservice-1/internal/infrastructure/event/message"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

// KafkaHeaders is the adapter interface to translate to kafka.Header slice
// which is used to compose a kafka message.
type KafkaHeaders interface {
	FromEchoContext(ctx echo.Context, event string) (headers []kafka.Header, err error)
}

// KafkaAdapter represent a specific implementation from the KafkaHeaders adapter interface.
type KafkaAdapter struct{}

// NewKafkaHeaders create KafkaAdapter and return the KafkaHeaders interface.
// Return KafkaHeaders interface.
func NewKafkaHeaders() KafkaHeaders {
	return KafkaAdapter{}
}

// FIXME Code duplicated from pkg/handler but if it is included a cycle dependency happens
// Find a better solution than duplicate it
func getEchoHeader(ctx echo.Context, key string, defvalues []string) []string {
	val := ctx.Request().Header.Get(key)
	if val == "" {
		return defvalues
	}
	return []string{val}
}

// FromEchoContext translate from an echo.Context to []kafka.Header
// ctx is the echo.Context from an http handler.
// event is an additional type to identify exactly the schema which match
// with the kafka message.
// Return headers a slice of kafka.Header and nil error when success, else
// an error reference filled and an empty slice of kafka.Header.
func (a KafkaAdapter) FromEchoContext(ctx echo.Context, e string) (headers []kafka.Header, err error) {
	if ctx == nil {
		return []kafka.Header{}, fmt.Errorf("ctx cannot be nil")
	}
	if e == "" {
		return []kafka.Header{}, fmt.Errorf("event cannot be an empty string")
	}
	var (
		headerKey string
	)

	headerKey = string(header.HdrRequestID)
	xrhInsightsRequestId := getEchoHeader(ctx, headerKey, []string{random.String(32)})

	// Fill headers
	headers = []kafka.Header{
		{
			Key:   string(message.HdrTopic),
			Value: []byte(event.TopicTodoCreated),
		},
		{
			Key:   string(header.HdrRequestID),
			Value: []byte(xrhInsightsRequestId[0]),
		},
	}

	return headers, nil
}
