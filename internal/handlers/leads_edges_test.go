package handlers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/lead"
	"hexletbasics/ent/tag"
	"hexletbasics/ent/tagging"
	"hexletbasics/internal/api"
	"hexletbasics/internal/events"
	"hexletbasics/internal/testsupport"
)

// The contact method names the column the value lands in: WhatsApp has its own.
func TestCreateLeadWritesTheWhatsAppColumnForAWhatsAppContact(t *testing.T) {
	h := testsupport.NewHarness(t)
	input := leadFromTheSite()
	input.ContactMethod = api.LeadInputContactMethodWhatsapp
	input.ContactValue = "+79990001122"

	res, err := h.Client.CreateLead(t.Context(), input)
	require.NoError(t, err)
	created, ok := res.(*api.Lead)
	require.True(t, ok, "got %T", res)

	stored := h.DB.Lead.Query().Where(lead.ID(int(created.ID))).OnlyX(t.Context())
	assert.Equal(t, "+79990001122", lo.FromPtr(stored.Whatsapp))
	assert.Nil(t, stored.Telegram)
	assert.Nil(t, stored.Phone)
}

// The sales team reads the lead by a name: the user's own, else their email,
// else a number — never an empty string.
func TestCreateLeadNamesTheUserAsLegacyDid(t *testing.T) {
	cases := []struct {
		name    string
		prepare func(context.Context, *testsupport.Harness)
		want    func(*testsupport.Harness) string
	}{
		{"by email when they have no name",
			func(ctx context.Context, h *testsupport.Harness) {
				h.DB.User.UpdateOneID(h.UserID).ClearFirstName().ClearLastName().ExecX(ctx)
			},
			func(*testsupport.Harness) string { return "alice@example.com" }},
		{"by number when they have neither",
			func(ctx context.Context, h *testsupport.Harness) {
				h.DB.User.UpdateOneID(h.UserID).ClearFirstName().ClearLastName().ClearEmail().ExecX(ctx)
			},
			func(h *testsupport.Harness) string { return fmt.Sprintf("User #%d", h.UserID) }},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			c.prepare(t.Context(), h)

			_, err := h.Client.CreateLead(t.Context(), leadFromTheSite())
			require.NoError(t, err)

			facts := testsupport.PublishedEvents[events.LeadCreated](t, h.Services)
			require.Len(t, facts, 1)
			assert.Equal(t, c.want(h), facts[0].UserName)
		})
	}
}

// A user who was never flagged is no change: the lead is stored and the tag's
// counter is untouched, whether the flag was never raised for them or there is
// no such tag at all.
func TestCreateLeadLeavesTheTagAloneForAUserWhoNeverHadIt(t *testing.T) {
	t.Run("the tag exists but not on the user", func(t *testing.T) {
		h := testsupport.NewHarness(t)
		ctx := t.Context()
		h.DB.Tagging.Delete().Where(tagging.TaggableType("User"), tagging.TaggableID(h.UserID)).ExecX(ctx)
		before := h.DB.Tag.Query().Where(tag.Name("should_be_lead")).OnlyX(ctx).TaggingsCount

		_, err := h.Client.CreateLead(ctx, leadFromTheSite())
		require.NoError(t, err)

		after := h.DB.Tag.Query().Where(tag.Name("should_be_lead")).OnlyX(ctx).TaggingsCount
		assert.Equal(t, before, after)
	})

	t.Run("there is no such tag", func(t *testing.T) {
		h := testsupport.NewHarness(t)
		ctx := t.Context()
		h.DB.Tagging.Delete().ExecX(ctx)
		h.DB.Tag.Delete().Where(tag.Name("should_be_lead")).ExecX(ctx)

		_, err := h.Client.CreateLead(ctx, leadFromTheSite())
		require.NoError(t, err)

		assert.True(t, h.DB.Lead.Query().Where(lead.UserID(h.UserID)).ExistX(ctx))
	})
}
