package events

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	wmsql "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/getsentry/sentry-go"
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
	sentryClient *sentry.Client,
	handlers ...cqrs.EventHandler,
) (*Runtime, error) {
	watermillLogger := watermill.NewSlogLogger(logger)
	router, err := message.NewRouter(message.RouterConfig{CloseTimeout: routerCloseTimeout}, watermillLogger)
	if err != nil {
		return nil, errs.Wrapf(err, "create Watermill router")
	}
	poisonPublisher, err := wmsql.NewPublisher(
		wmsql.BeginnerFromStdSQL(db),
		wmsql.PublisherConfig{SchemaAdapter: wmsql.DefaultPostgreSQLSchema{}},
		watermillLogger,
	)
	if err != nil {
		return nil, errs.Wrapf(err, "create poison queue publisher")
	}
	delivery, err := deliveryMiddleware(watermillLogger, poisonPublisher, sentryClient)
	if err != nil {
		return nil, err
	}
	router.AddMiddleware(delivery...)

	observerSubscriber, err := newSubscriber(db, observerConsumerGroup, watermillLogger)
	if err != nil {
		return nil, errs.Wrapf(err, "create domain event observer subscriber")
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
			return nil, errs.Wrapf(err, "create Watermill event processor")
		}
		if err := processor.AddHandlers(handlers...); err != nil {
			return nil, errs.Wrapf(err, "register domain event handlers")
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
//
// reportPoisoned sits between PoisonQueue and Retry, so it sees only the error
// that survived every retry: a message a retry saved is routine and never
// reaches Sentry, a poisoned one is exactly one event.
func deliveryMiddleware(
	logger watermill.LoggerAdapter,
	poison message.Publisher,
	sentryClient *sentry.Client,
) ([]message.HandlerMiddleware, error) {
	poisonQueue, err := middleware.PoisonQueue(poison, PoisonTopic)
	if err != nil {
		return nil, errs.Wrapf(err, "create poison queue middleware")
	}
	return []message.HandlerMiddleware{
		poisonQueue,
		reportPoisoned(sentryClient),
		middleware.Retry{
			MaxRetries:          handlerRetries,
			InitialInterval:     handlerRetryInitial,
			MaxInterval:         handlerRetryMax,
			Multiplier:          handlerRetryMultiplier,
			ResetContextOnRetry: true,
			Logger:              logger,
		}.Middleware,
		middleware.Recoverer,
	}, nil
}

// reportPoisoned captures the error of a message whose retries are exhausted,
// so growth of the poison table is noticed. It returns the error unchanged for
// PoisonQueue to act on. A nil client reports nothing.
func reportPoisoned(sentryClient *sentry.Client) message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			produced, err := next(msg)
			if err != nil && sentryClient != nil {
				scope := sentry.NewScope()
				scope.SetTags(map[string]string{
					"message.id":      msg.UUID,
					"message.name":    msg.Metadata.Get("name"),
					"message.handler": message.HandlerNameFromCtx(msg.Context()),
					"message.poison":  "true",
				})
				sentryClient.CaptureException(err, &sentry.EventHint{Context: msg.Context(), OriginalException: err}, scope)
			}
			return produced, err
		}
	}
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
		return errs.Wrapf(err, "run Watermill router")
	}
	return nil
}

// Running is closed once every subscriber is ready; the worker process waits on
// it before reporting itself started. Publishing never depends on it: events go
// to the outbox table whether or not a router is running.
func (r *Runtime) Running() <-chan struct{} { return r.router.Running() }

// Close drains subscribers and in-flight handlers.
func (r *Runtime) Close() error { return r.router.Close() }

// How the router closes and how a handler is retried before its message goes to
// the poison queue.
const (
	routerCloseTimeout     = 15 * time.Second
	handlerRetries         = 3
	handlerRetryInitial    = 100 * time.Millisecond
	handlerRetryMax        = 2 * time.Second
	handlerRetryMultiplier = 2
)
