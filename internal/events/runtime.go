package events

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	wmsql "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

const (
	observerConsumerGroup = "observe_domain_events_v1"

	// PoisonTopic receives a message once its handler has failed every retry.
	PoisonTopic = "domain_events_poison"
)

// Runtime exposes the Watermill Router as a blocking process actor. The process
// coordinator owns startup ordering and graceful shutdown.
type Runtime struct {
	router *message.Router
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
	poisonPublisher, err := wmsql.NewPublisher(
		wmsql.BeginnerFromStdSQL(db),
		wmsql.PublisherConfig{SchemaAdapter: wmsql.DefaultPostgreSQLSchema{}},
		watermillLogger,
	)
	if err != nil {
		return nil, fmt.Errorf("create poison queue publisher: %w", err)
	}
	delivery, err := deliveryMiddleware(watermillLogger, poisonPublisher)
	if err != nil {
		return nil, err
	}
	router.AddMiddleware(delivery...)

	observerSubscriber, err := newSubscriber(db, observerConsumerGroup, watermillLogger)
	if err != nil {
		return nil, fmt.Errorf("create domain event observer subscriber: %w", err)
	}
	// The observer takes every event type, so it is a raw handler rather than a
	// typed cqrs handler.
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
			Marshaler:         Marshaler,
			Logger:            watermillLogger,
		})
		if err != nil {
			return nil, fmt.Errorf("create Watermill event processor: %w", err)
		}
		if err := processor.AddHandlers(handlers...); err != nil {
			return nil, fmt.Errorf("register domain event handlers: %w", err)
		}
	}

	return &Runtime{router: router}, nil
}

// deliveryMiddleware is the per-message delivery policy, outermost first.
// PoisonQueue catches what is left after the retries and moves it to
// PoisonTopic with the failure reason in its metadata, so the message is acked
// and a handler that always fails cannot hold back its consumer group. Retry
// wraps Recoverer so that a panic, which Recoverer turns into an error, is
// retried like any other failure; with the order reversed a panic would skip
// Retry. Messages in the poison table are replayed by hand once the cause is
// fixed.
func deliveryMiddleware(
	logger watermill.LoggerAdapter,
	poison message.Publisher,
) ([]message.HandlerMiddleware, error) {
	poisonQueue, err := middleware.PoisonQueue(poison, PoisonTopic)
	if err != nil {
		return nil, fmt.Errorf("create poison queue middleware: %w", err)
	}
	return []message.HandlerMiddleware{
		poisonQueue,
		middleware.Retry{
			MaxRetries:          3,
			InitialInterval:     100 * time.Millisecond,
			MaxInterval:         2 * time.Second,
			Multiplier:          2,
			ResetContextOnRetry: true,
			Logger:              logger,
		}.Middleware,
		middleware.Recoverer,
	}, nil
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

// Run blocks while the Watermill router is alive, allowing the process
// supervisor to observe router termination directly.
func (r *Runtime) Run(ctx context.Context) error {
	if err := r.router.Run(ctx); err != nil {
		return fmt.Errorf("run Watermill router: %w", err)
	}
	return nil
}

// Running is closed once every subscriber is ready; the worker process waits on
// it before reporting itself started. Publishing never depends on it: events go
// to the outbox table whether or not a router is running.
func (r *Runtime) Running() <-chan struct{} { return r.router.Running() }

// Close drains subscribers and in-flight handlers.
func (r *Runtime) Close() error { return r.router.Close() }
