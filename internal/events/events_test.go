package events

import (
	"bytes"
	"errors"
	"log/slog"
	"slices"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Faults a test injects; each is told apart from the others by identity.
var (
	errFlaky       = errors.New("flaky")
	errAlwaysFails = errors.New("always fails")
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
		{StaffReadPersonalData{}, "staff_read_personal_data"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.name, nameOf(tt.event))
		})
	}
}

// deliver runs a handler through the real delivery policy and returns what
// ended up on the poison topic.
func deliver(t *testing.T, handler message.HandlerFunc) (<-chan *message.Message, error) {
	t.Helper()

	return deliverReporting(t, nil, handler)
}

// deliverReporting is deliver with a Sentry client for the poisoned messages.
func deliverReporting(t *testing.T, sentryClient *sentry.Client, handler message.HandlerFunc) (<-chan *message.Message, error) {
	t.Helper()

	logger := watermill.NopLogger{}
	poison := gochannel.NewGoChannel(gochannel.Config{}, logger)

	t.Cleanup(func() { _ = poison.Close() })
	poisoned, err := poison.Subscribe(t.Context(), PoisonTopic)
	require.NoError(t, err)
	chain, err := deliveryMiddleware(logger, poison, sentryClient)
	require.NoError(t, err)

	for _, c := range slices.Backward(chain) {
		handler = c(handler)
	}

	_, err = handler(message.NewMessage("event-1", nil))

	return poisoned, err
}

func TestDeliveryRetriesAHandlerThatPanics(t *testing.T) {
	attempts := 0
	poisoned, err := deliver(t, func(*message.Message) ([]*message.Message, error) {
		attempts++
		if attempts == 1 {
			panic("boom")
		}

		return nil, nil
	})

	require.NoError(t, err)
	assert.Equal(t, 2, attempts)
	assert.Empty(t, poisoned)
}

func TestDeliveryMovesAnAlwaysFailingMessageToThePoisonQueue(t *testing.T) {
	poisoned, err := deliver(t, func(*message.Message) ([]*message.Message, error) {
		return nil, errAlwaysFails
	})

	require.NoError(t, err, "the message is acked so the consumer group moves on")

	select {
	case msg := <-poisoned:
		assert.Equal(t, "event-1", msg.UUID)
		assert.Equal(t, "always fails", msg.Metadata.Get(middleware.ReasonForPoisonedKey))
	case <-time.After(time.Second):
		t.Fatal("message did not reach the poison queue")
	}
}

// The poison table grows silently otherwise: a message that exhausted its
// retries is one Sentry event, a message that a retry saved is none.
func TestPoisonedMessageIsReportedToSentry(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{Dsn: "https://public@example.com/1", Transport: transport})
	require.NoError(t, err)
	t.Cleanup(client.Close) // its batch processor goroutine would trip goleak

	attempts := 0
	_, err = deliverReporting(t, client, func(*message.Message) ([]*message.Message, error) {
		attempts++
		if attempts == 1 {
			return nil, errFlaky
		}

		return nil, nil
	})
	require.NoError(t, err)
	assert.Empty(t, transport.Events(), "a message a retry saved is not an incident")

	_, err = deliverReporting(t, client, func(*message.Message) ([]*message.Message, error) {
		return nil, errAlwaysFails
	})
	require.NoError(t, err)

	captured := transport.Events()
	require.Len(t, captured, 1)
	assert.Equal(t, "always fails", captured[0].Exception[0].Value)
	assert.Equal(t, "event-1", captured[0].Tags["message.id"])
}
