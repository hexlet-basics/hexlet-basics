package events

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	wmsql "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"

	"hexletbasics/ent"
	"hexletbasics/internal/store"
)

// TxPublisher is the transactional seam used by business modules. It does not
// implement a bus: the production adapter delegates persistence and delivery
// semantics to Watermill.
type TxPublisher interface {
	Publish(ctx context.Context, tx *sql.Tx, event Event) error
}

// StandalonePublisher records facts with no accompanying state mutation.
type StandalonePublisher interface {
	PublishStandalone(ctx context.Context, event Event) error
}

// Publisher writes Watermill messages through a caller-owned SQL transaction.
type Publisher struct {
	store  store.Transactor
	logger watermill.LoggerAdapter
}

// NewPublisher builds the PostgreSQL domain-event publisher.
func NewPublisher(txStore store.Transactor, logger *slog.Logger) *Publisher {
	return &Publisher{store: txStore, logger: watermill.NewSlogLogger(logger)}
}

// PublishStandalone owns a short transaction for facts such as sign-in.
func (p *Publisher) PublishStandalone(ctx context.Context, event Event) error {
	return p.store.WithinTx(ctx, func(tx *sql.Tx, _ *ent.Client) error {
		return p.Publish(ctx, tx, event)
	})
}

// Publish inserts the event into the outbox transaction. Atlas owns the
// Watermill schema, so runtime auto-initialization remains disabled.
func (p *Publisher) Publish(ctx context.Context, tx *sql.Tx, event Event) error {
	sqlPublisher, err := wmsql.NewPublisher(
		wmsql.TxFromStdSQL(tx),
		wmsql.PublisherConfig{SchemaAdapter: wmsql.DefaultPostgreSQLSchema{}},
		p.logger,
	)
	if err != nil {
		return fmt.Errorf("create Watermill transaction publisher: %w", err)
	}
	defer func() { _ = sqlPublisher.Close() }()

	marshaler := cqrs.JSONMarshaler{GenerateName: nameOf}
	bus, err := cqrs.NewEventBusWithConfig(sqlPublisher, cqrs.EventBusConfig{
		GeneratePublishTopic: func(cqrs.GenerateEventPublishTopicParams) (string, error) {
			return Topic, nil
		},
		Marshaler: marshaler,
		Logger:    p.logger,
		OnPublish: func(params cqrs.OnEventSendParams) error {
			if params.EventName == "" {
				return fmt.Errorf("domain event has no stable name")
			}
			params.Message.Metadata.Set("schema_version", schemaVersion)
			return nil
		},
	})
	if err != nil {
		return fmt.Errorf("create Watermill event bus: %w", err)
	}
	if err := bus.Publish(ctx, event); err != nil {
		return fmt.Errorf("publish %s: %w", event.eventName(), err)
	}
	return nil
}
