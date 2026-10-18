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

func TestEditorStoresTheEmailTheWayRegistrationDoes(t *testing.T) {
	editor := account.NewEditor(testsupport.NewClient(t))

	created, err := editor.Create(t.Context(), account.UserFields{Email: "  Mixed.Case@Example.TEST "})
	require.NoError(t, err)

	require.NotNil(t, created.Email)
	assert.Equal(t, "mixed.case@example.test", *created.Email)
}

func TestEditorAnEmailInUseIsAConflict(t *testing.T) {
	editor := account.NewEditor(testsupport.NewClient(t))

	_, err := editor.Create(t.Context(), account.UserFields{Email: "admin-learner@example.com"})

	require.ErrorIs(t, err, domain.ErrConflict)
}

func TestEditorUpdateClearsWhatIsLeftOut(t *testing.T) {
	db := testsupport.NewClient(t)
	editor := account.NewEditor(db)
	id := fixture.ID(fixture.UserBlank)
	first := "Robert"
	admin := true

	updated, err := editor.Update(t.Context(), id, account.UserFields{
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
	editor := account.NewEditor(testsupport.NewClient(t))

	_, err := editor.Update(t.Context(), 0, account.UserFields{Email: "nobody@example.com"})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestEditorKeepsTheLastAdministrator(t *testing.T) {
	db := testsupport.NewClient(t)
	editor := account.NewEditor(db)
	lastAdmin := fixture.ID(fixture.UserAdminLearner)
	notAdmin := false
	// The fixtures hold three administrators (the learner, the second one and the
	// one who is also staff); with the other two demoted, the learner is the last.
	db.User.UpdateOneID(fixture.ID(fixture.UserSecondAdmin)).SetAdmin(false).ExecX(t.Context())
	db.User.UpdateOneID(fixture.ID(fixture.UserAdminStaff)).SetAdmin(false).ExecX(t.Context())

	_, err := editor.Update(t.Context(), lastAdmin, account.UserFields{Email: "admin-learner@example.com", Admin: &notAdmin})
	require.ErrorIs(t, err, account.ErrLastAdmin)
	require.ErrorIs(t, err, domain.ErrConflict)

	require.ErrorIs(t, editor.Delete(t.Context(), lastAdmin), account.ErrLastAdmin)

	// With a second administrator the same changes go through.
	db.User.UpdateOneID(fixture.ID(fixture.UserSecondAdmin)).SetAdmin(true).ExecX(t.Context())
	_, err = editor.Update(t.Context(), lastAdmin, account.UserFields{Email: "admin-learner@example.com", Admin: &notAdmin})
	require.NoError(t, err)
}

func TestEditorDeleteRemovesTheUser(t *testing.T) {
	db := testsupport.NewClient(t)
	editor := account.NewEditor(db)
	id := fixture.ID(fixture.UserBlank)

	require.NoError(t, editor.Delete(t.Context(), id))

	exists, err := db.User.Query().Where(user.ID(id)).Exist(t.Context())
	require.NoError(t, err)
	assert.False(t, exists)
	require.ErrorIs(t, editor.Delete(t.Context(), id), domain.ErrNotFound)
}
