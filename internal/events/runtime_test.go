package events

import (
	"bytes"
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/ids"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

// syncBuffer is a log sink the router's goroutines and the test can share.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// TestRuntimeDeliversAPublishedEventToItsHandlerAndTheObserver runs the real
// runtime — the SQL subscribers, the delivery policy, the observer — from start
// to graceful stop, with an event published through the outbox the way the HTTP
// process does it.
func TestRuntimeDeliversAPublishedEventToItsHandlerAndTheObserver(t *testing.T) {
	db, err := store.NewDB(testdb.DatabaseURL())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	logs := &syncBuffer{}
	logger := slog.New(slog.NewTextHandler(logs, nil))

	group := "runtime_" + ids.New()
	userID := int(time.Now().UnixNano() & 0x3fffffff)
	received := make(chan int, 8)
	runtime, err := NewRuntime(db, logger,
		cqrs.NewEventHandler(group, func(_ context.Context, event *UserSignedUp) error {
			received <- event.UserID
			return nil
		}),
	)
	require.NoError(t, err)

	runCtx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	stopped := make(chan error, 1)
	go func() { stopped <- runtime.Run(runCtx) }()

	select {
	case <-runtime.Running():
	case err := <-stopped:
		t.Fatalf("runtime failed to start: %v", err)
	case <-time.After(15 * time.Second):
		t.Fatal("runtime did not report running")
	}

	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, NewPublisher(store.New(db), logger).Publish(t.Context(), tx, UserSignedUp{
		UserID: userID, Locale: "en", OccurredAt: time.Now().UTC(),
	}))
	require.NoError(t, tx.Commit())

	deadline := time.After(15 * time.Second)
	for delivered := false; !delivered; {
		select {
		case got := <-received:
			delivered = got == userID
		case <-deadline:
			t.Fatalf("handler did not receive user %d", userID)
		}
	}
	assert.Eventually(t, func() bool {
		return bytes.Contains([]byte(logs.String()), []byte("domain event observed"))
	}, 15*time.Second, 50*time.Millisecond, "the observer logs every event")

	cancel()
	select {
	case err := <-stopped:
		require.NoError(t, err)
	case <-time.After(20 * time.Second):
		t.Fatal("runtime did not stop after its context was cancelled")
	}
	require.NoError(t, runtime.Close())

	_, _ = db.ExecContext(context.Background(),
		`DELETE FROM watermill_offsets_domain_events WHERE consumer_group = $1`, group)
}

func TestRuntimeRunReportsARouterThatIsAlreadyRunning(t *testing.T) {
	db, err := store.NewDB(testdb.DatabaseURL())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	runtime, err := NewRuntime(db, slog.New(slog.NewTextHandler(&syncBuffer{}, nil)))
	require.NoError(t, err)
	runCtx, cancel := context.WithCancel(t.Context())
	first := make(chan error, 1)
	go func() { first <- runtime.Run(runCtx) }()
	t.Cleanup(func() {
		cancel()
		<-first
		_ = runtime.Close()
	})
	select {
	case <-runtime.Running():
	case <-time.After(15 * time.Second):
		t.Fatal("runtime did not report running")
	}

	err = runtime.Run(t.Context())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "run Watermill router")
}
