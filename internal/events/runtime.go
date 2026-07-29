package events

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	wmsql "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

const observerConsumerGroup = "observe_domain_events_v1"

// Runtime owns the Watermill Router lifecycle and reports asynchronous router
// failures to main so they trigger the same graceful shutdown as a signal.
type Runtime struct {
	router    *message.Router
	errors    chan error
	startOnce sync.Once
}

// NewRuntime builds the durable SQL subscriber and the initial PII-safe
// structured-log observer. Future integrations register as independent event
// handlers with their own stable consumer groups.
func NewRuntime(
	db *sql.DB,
	logger *slog.Logger,
	handlers ...cqrs.EventHandler,
) (*Runtime, error) {
	watermillLogger := watermill.NewSlogLogger(logger)
	router, err := message.NewRouter(message.RouterConfig{CloseTimeout: 15 * time.Second}, watermillLogger)
	if err != nil {
		return nil, fmt.Errorf("create Watermill router: %w", err)
	}
	router.AddMiddleware(
		middleware.Recoverer,
		middleware.Retry{
			MaxRetries:          3,
			InitialInterval:     100 * time.Millisecond,
			MaxInterval:         2 * time.Second,
			Multiplier:          2,
			ResetContextOnRetry: true,
			Logger:              watermillLogger,
		}.Middleware,
	)

	observerSubscriber, err := newSubscriber(db, observerConsumerGroup, watermillLogger)
	if err != nil {
		return nil, fmt.Errorf("create domain event observer subscriber: %w", err)
	}
	router.AddConsumerHandler(
		observerConsumerGroup,
		Topic,
		observerSubscriber,
		func(msg *message.Message) error {
			observeMessage(logger, msg)
			return nil
		},
	)

	if len(handlers) > 0 {
		processor, err := cqrs.NewEventProcessorWithConfig(router, cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return Topic, nil
			},
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return newSubscriber(db, params.HandlerName, watermillLogger)
			},
			AckOnUnknownEvent: true,
			Marshaler:         cqrs.JSONMarshaler{GenerateName: nameOf},
			Logger:            watermillLogger,
		})
		if err != nil {
			return nil, fmt.Errorf("create Watermill event processor: %w", err)
		}
		if err := processor.AddHandlers(handlers...); err != nil {
			return nil, fmt.Errorf("register domain event handlers: %w", err)
		}
	}

	return &Runtime{router: router, errors: make(chan error, 1)}, nil
}

func observeMessage(logger *slog.Logger, msg *message.Message) {
	logger.Info("domain event observed",
		"event_name", msg.Metadata.Get("name"),
		"event_id", msg.UUID,
		"schema_version", msg.Metadata.Get("schema_version"),
	)
}

func newSubscriber(
	db *sql.DB,
	consumerGroup string,
	logger watermill.LoggerAdapter,
) (*wmsql.Subscriber, error) {
	return wmsql.NewSubscriber(
		wmsql.BeginnerFromStdSQL(db),
		wmsql.SubscriberConfig{
			ConsumerGroup:  consumerGroup,
			SchemaAdapter:  wmsql.DefaultPostgreSQLSchema{},
			OffsetsAdapter: wmsql.DefaultPostgreSQLOffsetsAdapter{},
		},
		logger,
	)
}

// Start runs the router in the background and waits until all subscribers are
// ready, so HTTP traffic cannot publish before the consumer runtime is alive.
func (r *Runtime) Start(ctx context.Context) error {
	var startErr error
	r.startOnce.Do(func() {
		go func() {
			r.errors <- r.router.Run(ctx)
			close(r.errors)
		}()
		select {
		case <-r.router.Running():
		case err := <-r.errors:
			startErr = fmt.Errorf("start Watermill router: %w", err)
		}
	})
	return startErr
}

// Errors reports a terminal router failure after a successful start.
func (r *Runtime) Errors() <-chan error { return r.errors }

// Close drains subscribers and in-flight handlers.
func (r *Runtime) Close() error { return r.router.Close() }
