package accounts

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/user"
	"hexletbasics/internal/events"
	"hexletbasics/internal/ids"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := store.NewDB(testdb.DatabaseURL())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestRegistrarCommitsUserAndEventTogether(t *testing.T) {
	db := openTestDB(t)
	txStore := store.New(db)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	registrar := NewRegistrar(txStore, events.NewPublisher(txStore, logger))
	occurredAt := time.Date(2026, time.July, 29, 10, 11, 12, 0, time.UTC)
	registrar.now = func() time.Time { return occurredAt }
	email := fmt.Sprintf("watermill-%s@example.com", ids.New())
	firstName := "Ada"

	u, err := registrar.Register(t.Context(), Registration{
		Email: email, Password: "secret", FirstName: &firstName, Locale: "ru",
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(),
			`DELETE FROM watermill_domain_events WHERE payload->>'user_id' = $1`,
			fmt.Sprint(u.ID),
		)
		_, _ = db.ExecContext(context.Background(), `DELETE FROM users WHERE id = $1`, u.ID)
	})

	entClient := store.NewClient(db)
	saved, err := entClient.User.Query().Where(user.ID(u.ID)).Only(t.Context())
	require.NoError(t, err)
	assert.Equal(t, email, *saved.Email)

	var payloadJSON, metadataJSON []byte
	err = db.QueryRowContext(t.Context(), `
		SELECT payload, metadata
		FROM watermill_domain_events
		WHERE payload->>'user_id' = $1
	`, fmt.Sprint(u.ID)).Scan(&payloadJSON, &metadataJSON)
	require.NoError(t, err)

	var payload events.UserSignedUp
	require.NoError(t, json.Unmarshal(payloadJSON, &payload))
	assert.Equal(t, u.ID, payload.UserID)
	assert.Equal(t, email, *payload.Email)
	assert.Equal(t, firstName, *payload.FirstName)
	assert.Equal(t, "ru", payload.Locale)
	assert.Equal(t, occurredAt, payload.OccurredAt)

	var metadata map[string]string
	require.NoError(t, json.Unmarshal(metadataJSON, &metadata))
	assert.Equal(t, "user_signed_up", metadata["name"])
	assert.Equal(t, "1", metadata["schema_version"])
}

type failingPublisher struct{}

func (failingPublisher) Publish(context.Context, *sql.Tx, events.Event) error {
	return errors.New("outbox unavailable")
}

func TestRegistrarRollsBackUserWhenPublishingFails(t *testing.T) {
	db := openTestDB(t)
	registrar := NewRegistrar(store.New(db), failingPublisher{})
	email := fmt.Sprintf("watermill-rollback-%s@example.com", ids.New())

	_, err := registrar.Register(t.Context(), Registration{
		Email: email, Password: "secret", Locale: "en",
	})
	require.ErrorContains(t, err, "outbox unavailable")

	var users int
	require.NoError(t, db.QueryRowContext(t.Context(),
		`SELECT count(*) FROM users WHERE email = $1`, email,
	).Scan(&users))
	assert.Zero(t, users)
}

func TestRegistrarRejectsPasswordBeforeOpeningTransaction(t *testing.T) {
	registrar := NewRegistrar(nil, nil)
	_, err := registrar.Register(t.Context(), Registration{
		Email:    "long-password@example.com",
		Password: string(make([]byte, 73)),
	})
	require.ErrorIs(t, err, ErrPasswordProcessing)
}
