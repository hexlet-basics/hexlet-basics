package account_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/user"
	"hexletbasics/internal/account"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

const newUserEmail = "new-sign-up@example.com"

func TestRegistrarCreatesTheUser(t *testing.T) {
	db, _, txStore := testsupport.NewClientWithTxAndTransactor(t)
	registrar := account.NewRegistrar(txStore)
	firstName := "Ada"

	created, err := registrar.Register(t.Context(), account.Registration{
		Email: newUserEmail, Password: fixture.Password, FirstName: &firstName,
	})
	require.NoError(t, err)

	saved, err := db.User.Query().Where(user.ID(created.ID)).Only(t.Context())
	require.NoError(t, err)
	assert.Equal(t, newUserEmail, *saved.Email)
	assert.Equal(t, firstName, *saved.FirstName)
}

func TestRegistrarRejectsPasswordBeforeOpeningTransaction(t *testing.T) {
	registrar := account.NewRegistrar(nil)
	_, err := registrar.Register(t.Context(), account.Registration{
		Email:     "long-password@example.com",
		Password:  string(make([]byte, 73)),
		FirstName: nil,
	})
	require.ErrorIs(t, err, account.ErrPasswordProcessing)
}

// A unique violation on a user field is a refusal of the input, not a conflict
// of state: it is domain.ErrInvalid, and the error names the field so the
// adapter can put the message where the person can fix it.
func TestRegistrarNamesTheFieldOfATakenEmail(t *testing.T) {
	db, _, txStore := testsupport.NewClientWithTxAndTransactor(t)
	registrar := account.NewRegistrar(txStore)
	taken := db.User.Query().Where(user.EmailNotNil()).FirstX(t.Context())

	_, err := registrar.Register(t.Context(), account.Registration{
		Email: *taken.Email, Password: fixture.Password, FirstName: nil,
	})

	require.ErrorIs(t, err, domain.ErrInvalid)
	field, ok := domain.FieldOf(err)
	require.True(t, ok)
	assert.Equal(t, "email", field)
}
