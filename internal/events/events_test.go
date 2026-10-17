package events

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// deliver runs a handler through the real delivery policy and returns what
// ended up on the poison topic.
func deliver(t *testing.T, handler message.HandlerFunc) (<-chan *message.Message, error) {
	t.Helper()
	logger := watermill.NopLogger{}
	poison := gochannel.NewGoChannel(gochannel.Config{}, logger)
	t.Cleanup(func() { _ = poison.Close() })
	poisoned, err := poison.Subscribe(t.Context(), PoisonTopic)
	require.NoError(t, err)
	chain, err := deliveryMiddleware(logger, poison)
	require.NoError(t, err)
	for i := len(chain) - 1; i >= 0; i-- {
		handler = chain[i](handler)
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
		return nil, errors.New("always fails")
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
