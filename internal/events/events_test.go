package events

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/ids"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

func TestObserverLogsIdentityWithoutPII(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	msg := message.NewMessage("event-42", []byte(`{"email":"private@example.com","first_name":"Private"}`))
	msg.Metadata.Set("name", "user_signed_up")
	msg.Metadata.Set("schema_version", "1")
	email := "private@example.com"
	firstName := "Private"

	observeMessage(logger, msg)
	logged := output.String()
	assert.Contains(t, logged, `"event_id":"event-42"`)
	assert.Contains(t, logged, `"event_name":"user_signed_up"`)
	assert.NotContains(t, logged, email)
	assert.NotContains(t, logged, firstName)
}

func TestLegacyEventNames(t *testing.T) {
	tests := []struct {
		event Event
		name  string
	}{
		{UserSignedUp{}, "user_signed_up"},
		{UserSignedIn{}, "user_signed_in"},
		{BookRequested{}, "book_requested"},
		{CourseStarted{}, "course_started"},
		{CourseFinished{}, "course_finished"},
		{LessonStarted{}, "lesson_started"},
		{LessonFinished{}, "lesson_finished"},
		{SolutionChecked{}, "solution_checked"},
		{EmailConfirmed{}, "email_confirmed"},
		{LeadCreated{}, "lead_created"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.name, nameOf(tt.event))
		})
	}
}

func TestSQLRouterFansOutToIndependentConsumerGroups(t *testing.T) {
	db, err := store.NewDB(testdb.DatabaseURL())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	logger := watermill.NewSlogLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	require.NoError(t, err)

	processor, err := cqrs.NewEventProcessorWithConfig(router, cqrs.EventProcessorConfig{
		GenerateSubscribeTopic: func(cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
			return Topic, nil
		},
		SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return newSubscriber(db, params.HandlerName, logger)
		},
		AckOnUnknownEvent: true,
		Marshaler:         cqrs.JSONMarshaler{GenerateName: nameOf},
		Logger:            logger,
	})
	require.NoError(t, err)

	suffix := ids.New()
	groupA := "fanout_a_" + suffix
	groupB := "fanout_b_" + suffix
	type delivery struct {
		eventID string
		userID  int
	}
	deliveriesA := make(chan delivery, 8)
	deliveriesB := make(chan delivery, 8)
	handler := func(target chan<- delivery) func(context.Context, *UserSignedUp) error {
		return func(ctx context.Context, event *UserSignedUp) error {
			target <- delivery{
				eventID: cqrs.OriginalMessageFromCtx(ctx).UUID,
				userID:  event.UserID,
			}
			return nil
		}
	}
	require.NoError(t, processor.AddHandlers(
		cqrs.NewEventHandler(groupA, handler(deliveriesA)),
		cqrs.NewEventHandler(groupB, handler(deliveriesB)),
	))

	routerCtx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	routerErrors := make(chan error, 1)
	go func() { routerErrors <- router.Run(routerCtx) }()
	select {
	case <-router.Running():
	case err := <-routerErrors:
		t.Fatalf("router failed to start: %v", err)
	}
	t.Cleanup(func() { _ = router.Close() })

	targetUserID := int(time.Now().UnixNano() & 0x3fffffff)
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	publisher := NewPublisher(store.New(db), slog.New(slog.NewTextHandler(io.Discard, nil)))
	require.NoError(t, publisher.Publish(t.Context(), tx, UserSignedUp{
		UserID: targetUserID, Locale: "en", OccurredAt: time.Now().UTC(),
	}))
	require.NoError(t, tx.Commit())

	waitForTarget := func(ch <-chan delivery) delivery {
		t.Helper()
		timer := time.NewTimer(10 * time.Second)
		defer timer.Stop()
		for {
			select {
			case got := <-ch:
				if got.userID == targetUserID {
					return got
				}
			case <-timer.C:
				t.Fatalf("consumer did not receive user %d", targetUserID)
			}
		}
	}
	first := waitForTarget(deliveriesA)
	second := waitForTarget(deliveriesB)
	assert.Equal(t, first.eventID, second.eventID)

	require.NoError(t, router.Close())
	_, _ = db.ExecContext(context.Background(),
		`DELETE FROM watermill_offsets_domain_events WHERE consumer_group IN ($1, $2)`,
		groupA, groupB,
	)
	_, _ = db.ExecContext(context.Background(),
		`DELETE FROM watermill_domain_events WHERE uuid = $1`,
		first.eventID,
	)

	var groups int
	require.NoError(t, db.QueryRowContext(t.Context(), `
		SELECT count(*) FROM watermill_offsets_domain_events
		WHERE consumer_group IN ($1, $2)
	`, groupA, groupB).Scan(&groups))
	assert.Zero(t, groups, fmt.Sprintf("test consumer offsets must be cleaned: %s, %s", groupA, groupB))
}
