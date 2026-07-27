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
const totalUsers = 3

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
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListUsers(ctx, api.AdminListUsersParams{})
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
	h := testsupport.NewHarness(t)
	ctx := context.Background()

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
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// Matches on last name, case-insensitive.
	got, err := h.Client.AdminSearchUsers(ctx, api.AdminSearchUsersParams{Q: "brown"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	require.Len(t, got, 1)
	assert.Equal(t, "bob@example.com", got[0].Email.Value)

	// Matches on email substring.
	got, err = h.Client.AdminSearchUsers(ctx, api.AdminSearchUsersParams{Q: "example.com"})
	require.NoError(t, err)
	assert.Len(t, got, totalUsers)
}

func TestAdminSearchUsersBlank(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// A blank term returns an empty list, not every user.
	got, err := h.Client.AdminSearchUsers(ctx, api.AdminSearchUsersParams{Q: ""})
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestAdminGetUser(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	want := userByEmail(t, h, "alice@example.com")

	got, err := h.Client.AdminGetUser(ctx, api.AdminGetUserParams{ID: want.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "Alice", got.FirstName.Value)
	assert.True(t, got.Admin.Value)
}

func TestAdminGetUserNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminGetUser(ctx, api.AdminGetUserParams{ID: 999999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminCreateUser(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

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
// email index: a duplicate surfaces as 409 via the central ErrorHandler, not a
// 500, and no pre-check race exists.
func TestAdminCreateUserDuplicateEmail(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminCreateUser(ctx, &api.UserInput{
		Email:     "alice@example.com", // already seeded
		FirstName: api.NilString{Null: true},
		LastName:  api.NilString{Null: true},
		Admin:     api.NilBool{Null: true},
	})
	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())
}

func TestAdminUpdateUser(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	target := userByEmail(t, h, "bob@example.com")

	updated, err := h.Client.AdminUpdateUser(ctx, &api.UserInput{
		Email:     "bob@example.com",
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
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	target := userByEmail(t, h, "carol@example.com")

	err := h.Client.AdminDeleteUser(ctx, api.AdminDeleteUserParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	total, err := h.DB.User.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalUsers-1, total)
}
