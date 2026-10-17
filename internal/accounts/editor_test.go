package accounts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/user"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

func TestEditorStoresTheEmailTheWayRegistrationDoes(t *testing.T) {
	editor := accounts.NewEditor(testsupport.NewClient(t))

	created, err := editor.Create(t.Context(), accounts.UserFields{Email: "  Mixed.Case@Example.TEST "})
	require.NoError(t, err)

	require.NotNil(t, created.Email)
	assert.Equal(t, "mixed.case@example.test", *created.Email)
}

func TestEditorAnEmailInUseIsAConflict(t *testing.T) {
	editor := accounts.NewEditor(testsupport.NewClient(t))

	_, err := editor.Create(t.Context(), accounts.UserFields{Email: "admin-learner@example.com"})

	require.ErrorIs(t, err, domain.ErrConflict)
}

func TestEditorUpdateClearsWhatIsLeftOut(t *testing.T) {
	db := testsupport.NewClient(t)
	editor := accounts.NewEditor(db)
	id := fixtures.ID(fixtures.UserBlank)
	first := "Robert"
	admin := true

	updated, err := editor.Update(t.Context(), id, accounts.UserFields{
		Email: "blank-user@example.com", FirstName: &first, Admin: &admin,
	})
	require.NoError(t, err)

	require.NotNil(t, updated.FirstName)
	assert.Equal(t, "Robert", *updated.FirstName)
	assert.Nil(t, updated.LastName, "a nil name clears the column")
	require.NotNil(t, updated.Admin)
	assert.True(t, *updated.Admin)
}

func TestEditorUpdateOfAMissingUserIsNotFound(t *testing.T) {
	editor := accounts.NewEditor(testsupport.NewClient(t))

	_, err := editor.Update(t.Context(), 0, accounts.UserFields{Email: "nobody@example.com"})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestEditorKeepsTheLastAdministrator(t *testing.T) {
	db := testsupport.NewClient(t)
	editor := accounts.NewEditor(db)
	lastAdmin := fixtures.ID(fixtures.UserAdminLearner)
	notAdmin := false
	// The fixtures hold two administrators; with one of them demoted, the other is the last.
	db.User.UpdateOneID(fixtures.ID(fixtures.UserSecondAdmin)).SetAdmin(false).ExecX(t.Context())

	_, err := editor.Update(t.Context(), lastAdmin, accounts.UserFields{Email: "admin-learner@example.com", Admin: &notAdmin})
	require.ErrorIs(t, err, accounts.ErrLastAdmin)
	require.ErrorIs(t, err, domain.ErrConflict)

	require.ErrorIs(t, editor.Delete(t.Context(), lastAdmin), accounts.ErrLastAdmin)

	// With a second administrator the same changes go through.
	db.User.UpdateOneID(fixtures.ID(fixtures.UserSecondAdmin)).SetAdmin(true).ExecX(t.Context())
	_, err = editor.Update(t.Context(), lastAdmin, accounts.UserFields{Email: "admin-learner@example.com", Admin: &notAdmin})
	require.NoError(t, err)
}

func TestEditorDeleteRemovesTheUser(t *testing.T) {
	db := testsupport.NewClient(t)
	editor := accounts.NewEditor(db)
	id := fixtures.ID(fixtures.UserBlank)

	require.NoError(t, editor.Delete(t.Context(), id))

	exists, err := db.User.Query().Where(user.ID(id)).Exist(t.Context())
	require.NoError(t, err)
	assert.False(t, exists)
	require.ErrorIs(t, editor.Delete(t.Context(), id), domain.ErrNotFound)
}
