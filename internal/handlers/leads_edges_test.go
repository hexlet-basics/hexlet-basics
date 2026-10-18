package handlers_test

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/lead"
	"hexletbasics/ent/tag"
	"hexletbasics/internal/api"
	"hexletbasics/internal/events"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// The contact method names the column the value lands in: WhatsApp has its own.
func TestCreateLeadWritesTheWhatsAppColumnForAWhatsAppContact(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	input := leadFromTheSite()
	input.ContactMethod = api.LeadInputContactMethodWhatsapp
	input.ContactValue = "+79990001122"

	res, err := h.Client.CreateLead(t.Context(), input)
	require.NoError(t, err)

	created := res

	stored := h.DB.Lead.Query().Where(lead.ID(int(created.ID))).OnlyX(t.Context())
	assert.Equal(t, "+79990001122", lo.FromPtr(stored.Whatsapp))
	assert.Nil(t, stored.Telegram)
	assert.Nil(t, stored.Phone)
}

// The sales team reads the lead by a name: the user's own, else their email,
// else a number — never an empty string.
func TestCreateLeadNamesTheUserAsLegacyDid(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		user fixture.Label
		want func(*testsupport.Harness) string
	}{
		{"by email when they have no name", fixture.UserWithoutAName,
			func(*testsupport.Harness) string { return "without-a-name@example.com" }},
		{"by number when they have neither", fixture.UserWithoutANameOrEmail,
			func(h *testsupport.Harness) string { return fmt.Sprintf("User #%d", h.UserID) }},
	}

	for _, record := range cases {
		t.Run(record.name, func(t *testing.T) {
			h := testsupport.NewHarnessAs(t, record.user)

			_, err := h.Client.CreateLead(t.Context(), leadFromTheSite())
			require.NoError(t, err)

			facts := testsupport.PublishedEvents[events.LeadCreated](t, h.Services)
			require.Len(t, facts, 1)
			assert.Equal(t, record.want(h), facts[0].UserName)
		})
	}
}

// A user who was never flagged is no change: the lead is stored and the tag's
// counter is untouched. The default harness user carries no should_be_lead
// tagging, so this needs no delete and runs in parallel.
func TestCreateLeadLeavesTheTagAloneForAUserWhoNeverHadIt(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	before := h.DB.Tag.Query().Where(tag.Name("should_be_lead")).OnlyX(ctx).TaggingsCount

	_, err := h.Client.CreateLead(ctx, leadFromTheSite())
	require.NoError(t, err)

	after := h.DB.Tag.Query().Where(tag.Name("should_be_lead")).OnlyX(ctx).TaggingsCount
	assert.Equal(t, before, after)
}

// A database without the tag at all is no change either: the lead is stored.
//
// Not parallel: the lead code finds the tag by name and expects one, so there
// is one should_be_lead row for every test, and this one deletes it and every
// tagging, which would block or deadlock with any test that touches them.
func TestCreateLeadIsStoredWhenThereIsNoShouldBeLeadTag(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	h.DB.Tagging.Delete().ExecX(ctx)
	h.DB.Tag.Delete().Where(tag.Name("should_be_lead")).ExecX(ctx)

	_, err := h.Client.CreateLead(ctx, leadFromTheSite())
	require.NoError(t, err)

	assert.True(t, h.DB.Lead.Query().Where(lead.UserID(h.UserID)).ExistX(ctx))
}
