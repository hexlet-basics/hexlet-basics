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

// leads.yml seeds these two.
const totalLeads = 2

func TestAdminListLeads(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListLeads(ctx, api.AdminListLeadsParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(totalLeads), page.Total)
	assert.Len(t, page.Items, totalLeads)
	// Newest first: ids strictly descending.
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}

	// Business facts from the fixtures: emails present, full_name null (derived
	// from the user, not yet implemented), a nullable column left blank stays null.
	byEmail := map[string]api.Lead{}
	for _, it := range page.Items {
		byEmail[it.Email.Value] = it
	}
	anna, ok := byEmail["anna@example.com"]
	require.True(t, ok, "lead anna@example.com not found")
	assert.True(t, anna.FullName.Null, "fullName should be null until User schema lands")
	assert.Equal(t, "@anna", anna.Telegram.Value)
	assert.True(t, anna.Whatsapp.Null, "blank whatsapp column should be null")
	assert.False(t, anna.CreatedAt.IsZero())
}

func TestAdminListLeadsPaginated(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListLeads(ctx, api.AdminListLeadsParams{
		Page:    api.NewOptInt32(2),
		PerPage: api.NewOptInt32(1),
	})
	require.NoError(t, err)

	assert.Equal(t, int32(totalLeads), page.Total)
	assert.Equal(t, int32(2), page.Page)
	assert.Len(t, page.Items, 1) // 2 rows, page 2 of size 1 -> the 2nd row
}
