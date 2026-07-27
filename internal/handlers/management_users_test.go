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

// Reuses users.yml (totalUsers) — the management surface reads the same table.

func TestAdminListManagementUsers(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListManagementUsers(ctx, api.AdminListManagementUsersParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, int32(totalUsers), page.Total)
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}
}

func TestAdminGetManagementUser(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	want := userByEmail(t, h, "alice@example.com")

	got, err := h.Client.AdminGetManagementUser(ctx, api.AdminGetManagementUserParams{ID: want.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "Alice", got.FirstName.Value)
}

func TestAdminGetManagementUserNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminGetManagementUser(ctx, api.AdminGetManagementUserParams{ID: 999999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminUpdateManagementUser(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	target := userByEmail(t, h, "carol@example.com")

	updated, err := h.Client.AdminUpdateManagementUser(ctx, &api.UserInput{
		Email:     "carol@example.com",
		FirstName: api.NewNilString("Caroline"),
		LastName:  api.NilString{Null: true}, // clears the column
		Admin:     api.NewNilBool(true),
	}, api.AdminUpdateManagementUserParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "Caroline", updated.FirstName.Value)
	assert.True(t, updated.LastName.Null)
	assert.True(t, updated.Admin.Value)
}
