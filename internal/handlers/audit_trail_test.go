package handlers_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/events"
	"hexletbasics/internal/testsupport"
)

// The audit trail page is closed by a capability Permission: audit.read has to
// be granted to the Staff Role of whoever reads it, an administrator included.

// seedRecord stores one Audit Record directly, with the time and fields a test
// needs; the trail is append-only for the application, but a test may insert.
func seedRecord(t *testing.T, h *testsupport.Harness, at time.Time, operation string, mutate func(*audit.Entry)) int {
	t.Helper()
	entry := audit.Entry{Operation: operation, Outcome: audit.OutcomeSucceeded}
	if mutate != nil {
		mutate(&entry)
	}
	create := h.DB.AuditRecord.Create().
		SetOccurredAt(at).
		SetChannel(string(audit.ChannelWeb)).
		SetOperation(operation).
		SetOutcome(string(entry.Outcome))
	if entry.Actor != nil {
		create = create.SetActorUserID(*entry.Actor)
	}
	if entry.Source.Channel != "" {
		create = create.SetChannel(string(entry.Source.Channel))
	}
	if entry.EntityType != "" {
		create = create.SetEntityType(entry.EntityType)
	}
	if entry.EntityID != "" {
		create = create.SetEntityID(entry.EntityID)
	}
	return create.SaveX(context.Background()).ID
}

// grantAuditRead makes the signed-in administrator a staff member whose role
// holds audit.read, the way an operator would delegate the capability.
func grantAuditRead(t *testing.T, h *testsupport.Harness, withGrant bool) {
	t.Helper()
	ctx := context.Background()
	role := h.DB.StaffRole.Create().SetName("Auditors " + t.Name()).SetDescription("reads the trail").SaveX(ctx)
	h.DB.StaffMember.Create().SetUserID(h.UserID).SetRoleID(role.ID).SetAllowedLocales([]string{"ru", "en"}).ExecX(ctx)
	if withGrant {
		h.DB.StaffRoleGrant.Create().SetRoleID(role.ID).SetPermission("audit.read").ExecX(ctx)
	}
}

func listTrail(h *testsupport.Harness, params api.AdminListAuditRecordsParams) (*api.AuditRecordPage, error) {
	res, err := h.Client.AdminListAuditRecords(context.Background(), params)
	if err != nil {
		return nil, err
	}
	page, ok := res.(*api.AuditRecordPage)
	if !ok {
		return nil, assert.AnError
	}
	return page, nil
}

func TestAnAdministratorWithoutAStaffRecordIsRefusedTheTrail(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := listTrail(h, api.AdminListAuditRecordsParams{})

	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
}

func TestAnAdministratorWhoseRoleLacksAuditReadIsRefusedTheTrail(t *testing.T) {
	h := testsupport.NewHarness(t)
	grantAuditRead(t, h, false)

	_, err := listTrail(h, api.AdminListAuditRecordsParams{})

	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
}

func TestGrantingAuditReadOpensTheTrailToAnAdministrator(t *testing.T) {
	h := testsupport.NewHarness(t)
	grantAuditRead(t, h, true)

	_, err := listTrail(h, api.AdminListAuditRecordsParams{})

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
}

func TestAnEditorWithoutAuditReadIsRefusedTheTrail(t *testing.T) {
	h := testsupport.NewHarnessAs(t, editorID)

	_, err := listTrail(h, api.AdminListAuditRecordsParams{})

	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
}

func TestTheTrailListsNewestFirstAndPaginates(t *testing.T) {
	h := testsupport.NewHarness(t)
	grantAuditRead(t, h, true)
	now := time.Now()
	oldest := seedRecord(t, h, now.Add(-72*time.Hour), "AdminCreateBanner", nil)
	middle := seedRecord(t, h, now.Add(-48*time.Hour), "AdminUpdateBanner", nil)
	newest := seedRecord(t, h, now.Add(-24*time.Hour), "AdminDeleteBanner", nil)

	page, err := listTrail(h, api.AdminListAuditRecordsParams{PerPage: api.NewOptInt32(2)})
	require.NoError(t, err)

	// The records written by this very call and the fixtures are not in the way:
	// the three seeded rows are the only ones in the window except the ones the
	// listing's own refusals left, so assert on relative order.
	var ids []int64
	for _, item := range page.Items {
		ids = append(ids, item.ID)
	}
	require.Len(t, ids, 2)
	assert.Equal(t, int64(newest), ids[0])
	assert.Equal(t, int64(middle), ids[1])
	assert.GreaterOrEqual(t, page.Total, int32(3))
	assert.Equal(t, int32(2), page.PerPage)

	second, err := listTrail(h, api.AdminListAuditRecordsParams{PerPage: api.NewOptInt32(2), Page: api.NewOptInt32(2)})
	require.NoError(t, err)
	require.NotEmpty(t, second.Items)
	assert.Equal(t, int64(oldest), second.Items[0].ID)
}

func TestTheTrailFiltersByActorChannelOperationAndEntity(t *testing.T) {
	h := testsupport.NewHarness(t)
	grantAuditRead(t, h, true)
	now := time.Now().Add(-time.Hour)
	actor := 4242
	byActor := seedRecord(t, h, now, "AdminUpdateBanner", func(e *audit.Entry) { e.Actor = &actor })
	byMCP := seedRecord(t, h, now, "AdminUpdateReview", func(e *audit.Entry) { e.Source.Channel = audit.ChannelMCP })
	byEntity := seedRecord(t, h, now, "AdminUpdateBlogPost", func(e *audit.Entry) {
		e.EntityType = "BlogPost"
		e.EntityID = "77"
	})

	only := func(params api.AdminListAuditRecordsParams) []int64 {
		page, err := listTrail(h, params)
		require.NoError(t, err)
		var ids []int64
		for _, item := range page.Items {
			ids = append(ids, item.ID)
		}
		return ids
	}

	assert.Equal(t, []int64{int64(byActor)}, only(api.AdminListAuditRecordsParams{ActorUserId: api.NewOptInt32(int32(actor))}))
	assert.Equal(t, []int64{int64(byMCP)}, only(api.AdminListAuditRecordsParams{Channel: api.NewOptAuditRecordQueryChannel(api.AuditRecordQueryChannelMcp)}))
	assert.Equal(t, []int64{int64(byEntity)}, only(api.AdminListAuditRecordsParams{Operation: api.NewOptString("AdminUpdateBlogPost")}))
	assert.Equal(t, []int64{int64(byEntity)}, only(api.AdminListAuditRecordsParams{
		EntityType: api.NewOptString("BlogPost"),
		EntityId:   api.NewOptString("77"),
	}))
}

// The window is the date filter, and its default keeps the listing to the
// recent past: a record older than the default window is not listed unless the
// caller asks for it.
func TestTheTrailIsBoundedByAWindowThatDefaultsToThirtyDays(t *testing.T) {
	h := testsupport.NewHarness(t)
	grantAuditRead(t, h, true)
	now := time.Now()
	old := seedRecord(t, h, now.Add(-45*24*time.Hour), "AdminCreateBanner", nil)
	recent := seedRecord(t, h, now.Add(-2*24*time.Hour), "AdminUpdateBanner", nil)

	page, err := listTrail(h, api.AdminListAuditRecordsParams{Operation: api.NewOptString("AdminCreateBanner")})
	require.NoError(t, err)
	assert.Empty(t, page.Items, "a record older than the default window is outside it")

	page, err = listTrail(h, api.AdminListAuditRecordsParams{
		From:      api.NewOptDateTime(now.Add(-50 * 24 * time.Hour)),
		To:        api.NewOptDateTime(now.Add(-40 * 24 * time.Hour)),
		Operation: api.NewOptString("AdminCreateBanner"),
	})
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	assert.Equal(t, int64(old), page.Items[0].ID)

	page, err = listTrail(h, api.AdminListAuditRecordsParams{Operation: api.NewOptString("AdminUpdateBanner")})
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	assert.Equal(t, int64(recent), page.Items[0].ID)
}

func TestTheTrailRefusesAWindowItCannotAffordOrMakeSenseOf(t *testing.T) {
	h := testsupport.NewHarness(t)
	grantAuditRead(t, h, true)
	now := time.Now()

	cases := map[string]struct {
		params api.AdminListAuditRecordsParams
		field  string
	}{
		"wider than 92 days": {api.AdminListAuditRecordsParams{
			From: api.NewOptDateTime(now.Add(-93 * 24 * time.Hour)),
			To:   api.NewOptDateTime(now),
		}, "from"},
		"ends before it starts": {api.AdminListAuditRecordsParams{
			From: api.NewOptDateTime(now),
			To:   api.NewOptDateTime(now.Add(-time.Hour)),
		}, "to"},
		"an entity id without its type": {api.AdminListAuditRecordsParams{
			EntityId: api.NewOptString("5"),
		}, "entityId"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			res, err := h.Client.AdminListAuditRecords(context.Background(), c.params)
			require.NoError(t, err)

			problem, ok := res.(*api.ValidationError)
			require.True(t, ok, "got %T", res)
			assert.Contains(t, problem.Errors, c.field)
			assert.Equal(t, http.StatusUnprocessableEntity, h.LastStatus())
		})
	}
}

func TestTheTrailShowsChangesAndReads(t *testing.T) {
	h := testsupport.NewHarness(t)
	grantAuditRead(t, h, true)
	ctx := context.Background()
	at := time.Now().Add(-time.Hour)
	h.DB.AuditRecord.Create().
		SetOccurredAt(at).SetChannel("web").SetOperation("AdminUpdateBanner").SetOutcome("succeeded").
		SetChanges([]byte(`[{"entity":"Banner","id":"3","op":"update","fields":[{"name":"title","old":"A","new":"B"},{"name":"token","secret":true}]}]`)).
		ExecX(ctx)
	h.DB.AuditRecord.Create().
		SetOccurredAt(at).SetChannel("web").SetOperation("AdminListUsers").SetOutcome("succeeded").
		SetReads([]byte(`{"count":2,"ids":["1","2"]}`)).
		ExecX(ctx)

	page, err := listTrail(h, api.AdminListAuditRecordsParams{Operation: api.NewOptString("AdminUpdateBanner")})
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	change := page.Items[0].Changes[0]
	assert.Equal(t, "Banner", change.Entity)
	assert.Equal(t, "3", change.ID.Value)
	require.Len(t, change.Fields, 2)
	assert.JSONEq(t, `"A"`, string(change.Fields[0].Old))
	assert.JSONEq(t, `"B"`, string(change.Fields[0].New))
	assert.True(t, change.Fields[1].Secret.Value)
	assert.True(t, page.Items[0].Reads.Null)

	page, err = listTrail(h, api.AdminListAuditRecordsParams{Operation: api.NewOptString("AdminListUsers")})
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	assert.Empty(t, page.Items[0].Changes)
	assert.Equal(t, int32(2), page.Items[0].Reads.Value.Count)
	assert.Equal(t, []string{"1", "2"}, page.Items[0].Reads.Value.Ids)
}

// Reading the trail is itself a read of personal data (ADR-0017): the call
// leaves a record of AuditRecord rows being read once the worker has run, and
// the ids are those of the records the page showed.
func TestReadingTheTrailIsRecordedAsAReadOfPersonalData(t *testing.T) {
	h := testsupport.NewHarness(t)
	grantAuditRead(t, h, true)
	testsupport.ArriveFrom(h, "203.0.113.77")
	seeded := seedRecord(t, h, time.Now().Add(-time.Hour), "AdminCreateBanner", nil)

	page, err := listTrail(h, api.AdminListAuditRecordsParams{Operation: api.NewOptString("AdminCreateBanner")})
	require.NoError(t, err)
	require.Len(t, page.Items, 1)

	require.Len(t, testsupport.PublishedEvents[events.StaffReadPersonalData](t, h.Services), 1)
	processReads(t, h)

	rows := records(t, h, "AdminListAuditRecords")
	require.Len(t, rows, 1)
	assert.Equal(t, h.UserID, *rows[0].ActorUserID)
	assert.Equal(t, "AuditRecord", *rows[0].EntityType)
	assert.Equal(t, "203.0.113.77", *rows[0].Address)
	assert.JSONEq(t, `{"count":1,"ids":["`+strconv.Itoa(seeded)+`"]}`, string(rows[0].Reads))
}

// A refused call read nothing, so it leaves no read record.
func TestARefusedReadOfTheTrailIsNotRecordedAsARead(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := listTrail(h, api.AdminListAuditRecordsParams{})
	require.Error(t, err)
	processReads(t, h)

	assert.Empty(t, testsupport.PublishedEvents[events.StaffReadPersonalData](t, h.Services))
}
