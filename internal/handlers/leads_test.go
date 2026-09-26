package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/lead"
	"hexletbasics/internal/api"
	"hexletbasics/internal/events"
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

func leadFromTheSite() *api.LeadInput {
	return &api.LeadInput{
		ContactMethod: api.LeadInputContactMethodTelegram,
		ContactValue:  "@alice",
		YmClientId:    api.NewNilString("1700000000123456789"),
		FirstVisit: api.NewOptNilFirstVisit(api.FirstVisit{
			UtmSource:   api.NewNilString("vk"),
			UtmMedium:   api.NewNilString("cpc"),
			UtmCampaign: api.NewNilString("autumn"),
			UtmContent:  api.NewNilString("banner"),
			UtmTerm:     api.NewNilString("python"),
			LandingPage: api.NewNilString("https://code-basics.com/ru?utm_source=vk&gclid=g-1"),
			Referrer:    api.NewNilString("https://vk.com/"),
		}),
	}
}

func TestCreateLeadStoresTheLeadAndPublishesItsAttribution(t *testing.T) {
	h := testsupport.NewHarness(t)
	testsupport.ArriveFrom(h, "203.0.113.7")

	res, err := h.Client.CreateLead(t.Context(), leadFromTheSite())
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())
	created, ok := res.(*api.Lead)
	require.True(t, ok, "expected the created lead, got %T", res)

	// The contact method names the column the value lands in; the others stay
	// empty, as legacy's write_attribute left them.
	stored := h.DB.Lead.Query().Where(lead.ID(int(created.ID))).OnlyX(t.Context())
	assert.Equal(t, h.UserID, stored.UserID)
	assert.Equal(t, "@alice", lo.FromPtr(stored.Telegram))
	assert.Nil(t, stored.Phone)
	assert.Nil(t, stored.Whatsapp)
	assert.Equal(t, "1700000000123456789", lo.FromPtr(stored.YmClientID))
	// The harness user is enrolled in JavaScript with one finished lesson.
	assert.JSONEq(t, `[{"slug":"javascript","lessons_finished_count":1}]`, lo.FromPtr(stored.CoursesData))
	assert.JSONEq(t, `[]`, lo.FromPtr(stored.SurveyAnswersData))

	require.Len(t, h.Events.Published, 1)
	published, ok := h.Events.Published[0].(events.LeadCreated)
	require.True(t, ok, "expected LeadCreated, got %T", h.Events.Published[0])
	assert.Equal(t, stored.ID, published.LeadID)
	assert.Equal(t, h.UserID, published.UserID)
	assert.Equal(t, "Alice Anderson", published.UserName)
	assert.Equal(t, "alice@example.com", lo.FromPtr(published.Email))
	assert.Equal(t, "@alice", lo.FromPtr(published.Telegram))
	assert.Equal(t, "1700000000123456789", lo.FromPtr(published.YMClientID))
	assert.Equal(t, "vk", lo.FromPtr(published.UTMSource))
	assert.Equal(t, "cpc", lo.FromPtr(published.UTMMedium))
	assert.Equal(t, "autumn", lo.FromPtr(published.UTMCampaign))
	assert.Equal(t, "banner", lo.FromPtr(published.UTMContent))
	assert.Equal(t, "python", lo.FromPtr(published.UTMTerm))
	assert.Equal(t, "https://code-basics.com/ru?utm_source=vk&gclid=g-1", lo.FromPtr(published.LandingPage))
	assert.Equal(t, "https://vk.com/", lo.FromPtr(published.Referrer))
	assert.Equal(t, "203.0.113.7", lo.FromPtr(published.IP))
}

func TestCreateLeadWritesThePhoneColumnForAPhoneContact(t *testing.T) {
	h := testsupport.NewHarness(t)
	input := leadFromTheSite()
	input.ContactMethod = api.LeadInputContactMethodPhone
	input.ContactValue = "+79990000003"
	// A browser that kept no first-visit cookie still submits.
	input.FirstVisit = api.OptNilFirstVisit{}

	res, err := h.Client.CreateLead(t.Context(), input)
	require.NoError(t, err)
	created, ok := res.(*api.Lead)
	require.True(t, ok, "expected the created lead, got %T", res)

	stored := h.DB.Lead.Query().Where(lead.ID(int(created.ID))).OnlyX(t.Context())
	assert.Equal(t, "+79990000003", lo.FromPtr(stored.Phone))
	assert.Nil(t, stored.Telegram)

	require.Len(t, h.Events.Published, 1)
	published, ok := h.Events.Published[0].(events.LeadCreated)
	require.True(t, ok)
	assert.Equal(t, "+79990000003", lo.FromPtr(published.Phone))
	assert.Nil(t, published.UTMSource)
	assert.Nil(t, published.LandingPage)
}

func TestCreateLeadRequiresASignedInUser(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	res, err := h.Client.CreateLead(t.Context(), leadFromTheSite())
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus())
	assert.IsType(t, &api.ProblemDetails{}, res)
	assert.Empty(t, h.Events.Published)
	assert.Equal(t, totalLeads, h.DB.Lead.Query().CountX(t.Context()))
}
