package account_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/account"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

func TestEditorStoresTheEmailTheWayRegistrationDoes(t *testing.T) {
	editor := account.NewEditor(testsupport.NewClient(t))

	created, err := editor.Create(t.Context(), administrator(), account.UserFields{
		Email:     "  Mixed.Case@Example.TEST ",
		FirstName: nil,
		LastName:  nil,
		Admin:     nil,
	})
	require.NoError(t, err)

	require.NotNil(t, created.Email)
	assert.Equal(t, "mixed.case@example.test", *created.Email)
}

func TestEditorAnEmailInUseIsInvalidInputNamingTheEmail(t *testing.T) {
	editor := account.NewEditor(testsupport.NewClient(t))

	_, err := editor.Create(t.Context(), administrator(), account.UserFields{
		Email:     "admin-learner@example.com",
		FirstName: nil,
		LastName:  nil,
		Admin:     nil,
	})

	require.ErrorIs(t, err, domain.ErrInvalid)
	field, ok := domain.FieldOf(err)
	require.True(t, ok)
	assert.Equal(t, "email", field)
}

func TestEditorUpdateClearsWhatIsLeftOut(t *testing.T) {
	db := testsupport.NewClient(t)
	editor := account.NewEditor(db)
	id := fixture.ID(fixture.UserBlank)
	first := "Robert"
	raised := true

	updated, err := editor.Update(t.Context(), administrator(), id, account.UserFields{
		Email: "blank-user@example.com", FirstName: &first, Admin: &raised,
		LastName: nil,
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

	_, err := editor.Update(t.Context(), administrator(), 0, account.UserFields{
		Email:     "nobody@example.com",
		FirstName: nil,
		LastName:  nil,
		Admin:     nil,
	})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestEditorKeepsTheLastAdministrator(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	editor := account.NewEditor(db)
	lastAdmin := fixture.ID(fixture.UserAdminLearner)
	notAdmin := false
	// The fixtures hold three administrators (the learner, the second one and the
	// one who is also staff); with the other two demoted, the learner is the last.
	db.User.UpdateOneID(fixture.ID(fixture.UserSecondAdmin)).SetAdmin(false).ExecX(t.Context())
	db.User.UpdateOneID(fixture.ID(fixture.UserAdminStaff)).SetAdmin(false).ExecX(t.Context())

	_, err := editor.Update(t.Context(), administrator(), lastAdmin, account.UserFields{
		Email:     "admin-learner@example.com",
		Admin:     &notAdmin,
		FirstName: nil,
		LastName:  nil,
	})
	require.ErrorIs(t, err, account.ErrLastAdmin)
	require.ErrorIs(t, err, domain.ErrConflict)

	removeErr := account.NewRemover(transactor).Remove(t.Context(), administrator(), lastAdmin)
	require.ErrorIs(t, removeErr, account.ErrLastAdmin)

	// With a second administrator the same changes go through.
	db.User.UpdateOneID(fixture.ID(fixture.UserSecondAdmin)).SetAdmin(true).ExecX(t.Context())
	_, err = editor.Update(t.Context(), administrator(), lastAdmin, account.UserFields{
		Email:     "admin-learner@example.com",
		Admin:     &notAdmin,
		FirstName: nil,
		LastName:  nil,
	})
	require.NoError(t, err)
}

// admin is an administrator acting through the editor, someone other than the
// users the tests change.
func administrator() domain.Actor { return domain.Actor{UserID: 0, Administrator: true} }

func TestOnlyAnAdministratorChangesTheAdministratorFlag(t *testing.T) {
	db := testsupport.NewClient(t)
	editor := account.NewEditor(db)
	staff := domain.Actor{UserID: fixture.ID(fixture.UserStaffEditor), Administrator: false}
	yes := true

	_, err := editor.Create(t.Context(), staff, account.UserFields{Email: "raised@example.com", Admin: &yes, FirstName: nil, LastName: nil})
	require.ErrorIs(t, err, account.ErrAdministratorFlagRestricted)
	require.ErrorIs(t, err, domain.ErrForbidden)

	id := fixture.ID(fixture.UserBlank)
	_, err = editor.Update(t.Context(), staff, id, account.UserFields{Email: "blank@example.com", Admin: &yes, FirstName: nil, LastName: nil})
	require.ErrorIs(t, err, account.ErrAdministratorFlagRestricted)
	assert.False(t, lo.FromPtr(db.User.GetX(t.Context(), id).Admin), "a refused change writes nothing")

	_, err = editor.Update(t.Context(), administrator(), id, account.UserFields{Email: "blank@example.com", Admin: &yes, FirstName: nil, LastName: nil})
	require.NoError(t, err)
}
