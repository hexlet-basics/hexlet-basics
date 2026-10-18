package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// users.yml seeds these three.

// typeaheadCap is how many users the admin typeahead returns at most (the legacy
// `.limit(20)`).
const typeaheadCap = 20

// userByEmail resolves a fixture user by its email (a stable business key),
// reading its id back rather than hard-coding it.
func userByEmail(t *testing.T, h *testsupport.Harness, email string) api.UserCrud {
	t.Helper()

	page, err := h.Client.AdminListUsers(context.Background(), api.AdminListUsersParams{
		PerPage: api.NewOptInt32(100),
	})
	require.NoError(t, err)

	for _, it := range page.Items {
		if it.Email.Value == email {
			return it
		}
	}

	t.Fatalf("user with email %q not found", email)

	return api.UserCrud{}
}

func TestAdminListUsers(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalUsers := h.DB.User.Query().CountX(ctx)

	page, err := h.Client.AdminListUsers(ctx, api.AdminListUsersParams{PerPage: api.NewOptInt32(100)})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(totalUsers), page.Total)
	assert.Len(t, page.Items, totalUsers)
	// Newest first: ids strictly descending.
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}
}

func TestAdminListUsersPaginated(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalUsers := h.DB.User.Query().CountX(ctx)

	page, err := h.Client.AdminListUsers(ctx, api.AdminListUsersParams{
		Page:    api.NewOptInt32(2),
		PerPage: api.NewOptInt32(1),
	})
	require.NoError(t, err)

	assert.Equal(t, int32(totalUsers), page.Total)
	assert.Equal(t, int32(2), page.Page)
	assert.Len(t, page.Items, 1)
}

func TestAdminSearchUsers(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalUsers := h.DB.User.Query().CountX(ctx)

	// Matches on last name, case-insensitive (user_staff_editor is the only Editor).
	got, err := h.Client.AdminSearchUsers(ctx, api.AdminSearchUsersParams{Q: "editor"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	require.Len(t, got, 1)
	assert.Equal(t, "staff-editor@example.com", got[0].Email.Value)

	// Matches on email substring.
	// Every fixture user matches; the typeahead returns at most twenty of them.
	got, err = h.Client.AdminSearchUsers(ctx, api.AdminSearchUsersParams{Q: "example.com"})
	require.NoError(t, err)
	assert.Len(t, got, min(totalUsers, typeaheadCap))
}

func TestAdminSearchUsersBlank(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// A blank term returns an empty list, not every user.
	got, err := h.Client.AdminSearchUsers(ctx, api.AdminSearchUsersParams{Q: ""})
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestAdminGetUser(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	want := userByEmail(t, h, "admin-learner@example.com")

	got, err := h.Client.AdminGetUser(ctx, api.AdminGetUserParams{ID: want.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "Admin", got.FirstName.Value)
	assert.True(t, got.Admin.Value)
}

func TestAdminGetUserNotFound(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminGetUser(ctx, api.AdminGetUserParams{ID: 999999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminCreateUser(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalUsers := h.DB.User.Query().CountX(ctx)

	created, err := h.Client.AdminCreateUser(ctx, &api.UserInput{
		Email:     "dave@example.com",
		FirstName: api.NewNilString("Dave"),
		LastName:  api.NewNilString("Davis"),
		Admin:     api.NewNilBool(true),
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())
	assert.NotZero(t, created.ID)
	assert.Equal(t, "dave@example.com", created.Email.Value)
	assert.True(t, created.Admin.Value)

	total, err := h.DB.User.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalUsers+1, total)
}

// TestAdminCreateUserDuplicateEmail is the parity proof for the baseline unique
// email index: a duplicate is the person's input to fix, a 422 that names the
// email field through the central ErrorHandler, not a 500, and no pre-check
// race exists.
func TestAdminCreateUserDuplicateEmail(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminCreateUser(ctx, &api.UserInput{
		Email:     "admin-learner@example.com", // already seeded
		FirstName: api.NilString{Null: true},
		LastName:  api.NilString{Null: true},
		Admin:     api.NilBool{Null: true},
	})
	problem := problemWithStatus(t, err, http.StatusUnprocessableEntity)
	assert.Contains(t, problem.Response.Errors.Value, "email")
}

func TestAdminUpdateUser(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	target := userByEmail(t, h, "staff-editor@example.com")

	updated, err := h.Client.AdminUpdateUser(ctx, &api.UserInput{
		Email:     "staff-editor@example.com",
		FirstName: api.NewNilString("Robert"),
		LastName:  api.NilString{Null: true}, // clears the column
		Admin:     api.NewNilBool(true),
	}, api.AdminUpdateUserParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "Robert", updated.FirstName.Value)
	assert.True(t, updated.LastName.Null)
	assert.True(t, updated.Admin.Value)

	row, err := h.DB.User.Get(ctx, int(target.ID))
	require.NoError(t, err)
	assert.Nil(t, row.LastName, "null lastName should clear the column")
}

func TestAdminDeleteUser(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalUsers := h.DB.User.Query().CountX(ctx)

	target := userByEmail(t, h, "blank-user@example.com")

	err := h.Client.AdminDeleteUser(ctx, api.AdminDeleteUserParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	total, err := h.DB.User.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalUsers-1, total)
}

func TestAdminCannotRemoveOrDemoteTheLastAdministrator(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	// The fixtures have three administrators: user_admin_learner (the caller), user_second_admin
	// and the staff member user_admin_staff; the caller is left the only one.
	leaveOnlyAlice(t, h)

	input := &api.UserInput{
		Email:     "admin-learner@example.com",
		FirstName: api.NilString{Null: true},
		LastName:  api.NilString{Null: true},
		Admin:     api.NewNilBool(false),
	}

	_, err := h.Client.AdminUpdateUser(ctx, input, api.AdminUpdateUserParams{ID: int32(h.UserID)})
	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())

	// Deleting oneself is refused on its own account (self-protection_test.go);
	// the last administrator is also protected from someone else.
	err = h.Client.AdminDeleteUser(ctx, api.AdminDeleteUserParams{ID: int32(h.UserID)})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())

	// With a second administrator the same changes go through.
	h.DB.User.UpdateOneID(secondAdminID).SetAdmin(true).ExecX(ctx)
	_, err = h.Client.AdminUpdateUser(ctx, input, api.AdminUpdateUserParams{ID: int32(h.UserID)})
	require.NoError(t, err)
}

func TestAdminStoresTheEmailTheWayRegistrationDoes(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	created, err := h.Client.AdminCreateUser(t.Context(), &api.UserInput{
		Email:     "Mixed.Case@Example.TEST",
		FirstName: api.NilString{Null: true},
		LastName:  api.NilString{Null: true},
		Admin:     api.NilBool{Null: true},
	})
	require.NoError(t, err)

	assert.Equal(t, api.NewNilString("mixed.case@example.test"), created.Email)
}
