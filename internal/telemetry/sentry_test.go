package telemetry_test

import (
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/config"
	"hexletbasics/internal/telemetry"
)

func newSentryClient(t *testing.T) (*sentry.Client, *sentry.MockTransport) {
	t.Helper()

	transport := &sentry.MockTransport{}
	client, err := telemetry.NewSentryClient(
		&config.Config{SentryDSN: "https://public@example.com/1"},
		telemetry.WithSentryTransport(transport),
	)
	require.NoError(t, err)

	return client, transport
}

func TestSentryClientCopiesOopsAttributesIntoTheEvent(t *testing.T) {
	client, transport := newSentryClient(t)

	err := oops.In("lessons").With("lesson_id", 42).Wrapf(assert.AnError, "load lesson")
	client.CaptureException(err, &sentry.EventHint{OriginalException: err}, nil)

	events := transport.Events()
	require.Len(t, events, 1)
	assert.Equal(t, "lessons", events[0].Tags["oops.domain"])
	assert.Equal(t, 42, events[0].Contexts["oops"]["lesson_id"])
}

func TestSentryClientLeavesPlainErrorsAlone(t *testing.T) {
	client, transport := newSentryClient(t)

	client.CaptureException(assert.AnError, &sentry.EventHint{OriginalException: assert.AnError}, nil)

	events := transport.Events()
	require.Len(t, events, 1)
	assert.NotContains(t, events[0].Contexts, "oops")
}
