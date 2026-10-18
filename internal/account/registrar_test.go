package account_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/user"
	"hexletbasics/internal/account"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/events"
	"hexletbasics/internal/testsupport"
)

// Faults a test injects; each is told apart from the others by identity.
var (
	errOutboxUnavailable = errors.New("outbox unavailable")
)

const newUserEmail = "new-sign-up@example.com"

func TestRegistrarCommitsUserAndEventTogether(t *testing.T) {
	db, tx, txStore := testsupport.NewClientWithTxAndTransactor(t)
	occurredAt := time.Date(2026, time.July, 29, 10, 11, 12, 0, time.UTC)
	registrar := account.NewRegistrar(txStore, events.NewPublisher(txStore, slog.New(slog.DiscardHandler)),
		account.WithClock(func() time.Time { return occurredAt }))
	firstName := "Ada"

	created, err := registrar.Register(t.Context(), account.Registration{
		Email: newUserEmail, Password: "secret", FirstName: &firstName, Locale: "ru",
	})
	require.NoError(t, err)

	saved, err := db.User.Query().Where(user.ID(created.ID)).Only(t.Context())
	require.NoError(t, err)
	assert.Equal(t, newUserEmail, *saved.Email)

	var payloadJSON, metadataJSON []byte

	err = tx.QueryRowContext(t.Context(), `
		SELECT payload, metadata
		FROM watermill_domain_events
		WHERE payload->>'user_id' = $1
	`, strconv.Itoa(created.ID)).Scan(&payloadJSON, &metadataJSON)
	require.NoError(t, err)

	var payload events.UserSignedUp
	require.NoError(t, json.Unmarshal(payloadJSON, &payload))
	assert.Equal(t, created.ID, payload.UserID)
	assert.Equal(t, newUserEmail, *payload.Email)
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
	return errOutboxUnavailable
}

func TestRegistrarRollsBackUserWhenPublishingFails(t *testing.T) {
	db, _, txStore := testsupport.NewClientWithTxAndTransactor(t)
	registrar := account.NewRegistrar(txStore, failingPublisher{})

	_, err := registrar.Register(t.Context(), account.Registration{
		Email: newUserEmail, Password: "secret", Locale: "en",
		FirstName: nil,
	})
	require.ErrorContains(t, err, "outbox unavailable")

	assert.False(t, db.User.Query().Where(user.Email(newUserEmail)).ExistX(t.Context()), "no user is left behind")
}

func TestRegistrarRejectsPasswordBeforeOpeningTransaction(t *testing.T) {
	registrar := account.NewRegistrar(nil, nil)
	_, err := registrar.Register(t.Context(), account.Registration{
		Email:     "long-password@example.com",
		Password:  string(make([]byte, 73)),
		FirstName: nil,
		Locale:    "",
	})
	require.ErrorIs(t, err, account.ErrPasswordProcessing)
}

// A unique violation on a user field is a refusal of the input, not a conflict
// of state: it is domain.ErrInvalid, and the error names the field so the
// adapter can put the message where the person can fix it.
func TestRegistrarNamesTheFieldOfATakenEmail(t *testing.T) {
	db, _, txStore := testsupport.NewClientWithTxAndTransactor(t)
	registrar := account.NewRegistrar(txStore, events.NewPublisher(txStore, slog.New(slog.DiscardHandler)))
	taken := db.User.Query().Where(user.EmailNotNil()).FirstX(t.Context())

	_, err := registrar.Register(t.Context(), account.Registration{
		Email: *taken.Email, Password: "secret", Locale: "en",
		FirstName: nil,
	})

	require.ErrorIs(t, err, domain.ErrInvalid)
	field, ok := domain.FieldOf(err)
	require.True(t, ok)
	assert.Equal(t, "email", field)
}
