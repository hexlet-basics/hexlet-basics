package handlers_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/auditrecord"
	"hexletbasics/internal/api"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/events"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// The audit trail page is closed by a capability Permission: audit.read has to
// be granted to the Staff Role of whoever reads it, an administrator included.

// seedRecord writes one Audit Record through the production writer, with the
// instant and fields a test needs: the writer takes the time of a read of
// personal data, which is recorded after the request, so a record can be placed
// in the past without a back door.
func seedRecord(t *testing.T, h *testsupport.Harness, at time.Time, operation string, mutate func(*audit.Entry)) int {
	t.Helper()
	entry := audit.Entry{
		OccurredAt: at,
		Source:     audit.Source{Channel: audit.ChannelWeb},
		Operation:  operation,
		Outcome:    audit.OutcomeSucceeded,
	}
	if mutate != nil {
		mutate(&entry)
	}
	require.NoError(t, audit.Write(t.Context(), h.DB, entry, nil))
	return newestRecord(t, h)
}

// newestRecord is the id of the Audit Record written last in the test's transaction.
func newestRecord(t *testing.T, h *testsupport.Harness) int {
	t.Helper()
	return h.DB.AuditRecord.Query().Order(ent.Desc(auditrecord.FieldID)).FirstX(t.Context()).ID
}

// auditReader is a signed-in administrator who is also staff (user_admin_staff,
// on the editor role), the way an operator who may be delegated the capability
// is: with the grant on the role when withGrant is true, without it otherwise.
func auditReader(t *testing.T, withGrant bool) *testsupport.Harness {
	t.Helper()
	h := testsupport.NewHarnessAs(t, fixture.UserAdminStaff)
	if withGrant {
		h.Grant(t, int(roleEditorID), "audit.read")
	}
	return h
}

func listTrail(h *testsupport.Harness, params api.AdminListAuditRecordsParams) (*api.AuditRecordPage, error) {
	res, err := h.Client.AdminListAuditRecords(context.Background(), params)
	if err != nil {
		return nil, err
	}
	page := res
	return page, nil
}

func TestAnAdministratorWithoutAStaffRecordIsRefusedTheTrail(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	_, err := listTrail(h, api.AdminListAuditRecordsParams{})

	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
}

func TestAnAdministratorWhoseRoleLacksAuditReadIsRefusedTheTrail(t *testing.T) {
	t.Parallel()
	h := auditReader(t, false)

	_, err := listTrail(h, api.AdminListAuditRecordsParams{})

	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
}

func TestGrantingAuditReadOpensTheTrailToAnAdministrator(t *testing.T) {
	t.Parallel()
	h := auditReader(t, true)

	_, err := listTrail(h, api.AdminListAuditRecordsParams{})

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
}

func TestAnEditorWithoutAuditReadIsRefusedTheTrail(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)

	_, err := listTrail(h, api.AdminListAuditRecordsParams{})

	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
}

func TestTheTrailListsNewestFirstAndPaginates(t *testing.T) {
	t.Parallel()
	h := auditReader(t, true)
	now := time.Now()
	oldest := seedRecord(t, h, now.Add(-72*time.Hour), "AdminCreateBanner", nil)
	middle := seedRecord(t, h, now.Add(-48*time.Hour), "AdminUpdateBanner", nil)
	newest := seedRecord(t, h, now.Add(-24*time.Hour), "AdminDeleteBanner", nil)

	page, err := listTrail(h, api.AdminListAuditRecordsParams{PerPage: api.NewOptInt32(2)})
	require.NoError(t, err)

	// The records written by this very call and the fixtures are not in the way:
	// the three seeded rows are the only ones in the window except the ones the
	// listing's own refusals left, so assert on relative order.
	ids := make([]int64, 0, len(page.Items))
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
	t.Parallel()
	h := auditReader(t, true)
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
		ids := make([]int64, 0, len(page.Items))
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
	t.Parallel()
	h := auditReader(t, true)
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
	t.Parallel()
	h := auditReader(t, true)
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
			_, err := h.Client.AdminListAuditRecords(context.Background(), c.params)

			problem := problemWithStatus(t, err, http.StatusUnprocessableEntity)
			assert.Contains(t, problem.Response.Errors.Value, c.field)
			assert.Equal(t, http.StatusUnprocessableEntity, h.LastStatus())
		})
	}
}

func TestTheTrailShowsChangesAndReads(t *testing.T) {
	t.Parallel()
	h := auditReader(t, true)
	at := time.Now().Add(-time.Hour)
	banner, carol := fixture.ID(fixture.BannerEnDraft), fixture.ID(fixture.UserWithPassword)

	// A change is what the operation's own writes left in the Recorder: a banner
	// text and a password digest, which the trail keeps by name only.
	ctx, changes := audit.WithRecorder(t.Context())
	h.DB.Banner.UpdateOneID(banner).SetBody("Cyber Monday sale").ExecX(ctx)
	h.DB.User.UpdateOneID(carol).SetPasswordDigest("digest-the-trail-must-not-keep").ExecX(ctx)
	require.NoError(t, audit.Write(t.Context(), h.DB, audit.Entry{
		OccurredAt: at, Source: audit.Source{Channel: audit.ChannelWeb},
		Operation: "AdminUpdateBanner", Outcome: audit.OutcomeSucceeded,
	}, changes))
	require.NoError(t, audit.Write(t.Context(), h.DB, audit.Entry{
		OccurredAt: at, Source: audit.Source{Channel: audit.ChannelWeb},
		Operation: "AdminListUsers", Outcome: audit.OutcomeSucceeded,
		Reads: &audit.Reads{Count: 2, IDs: []string{"1", "2"}},
	}, nil))

	page, err := listTrail(h, api.AdminListAuditRecordsParams{Operation: api.NewOptString("AdminUpdateBanner")})
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	require.Len(t, page.Items[0].Changes, 2)
	change := page.Items[0].Changes[0]
	assert.Equal(t, "Banner", change.Entity)
	assert.Equal(t, strconv.Itoa(banner), change.ID.Value)
	require.Len(t, change.Fields, 1)
	assert.JSONEq(t, `"Black Friday sale"`, string(change.Fields[0].Old))
	assert.JSONEq(t, `"Cyber Monday sale"`, string(change.Fields[0].New))
	secret := page.Items[0].Changes[1].Fields[0]
	assert.Equal(t, "password_digest", secret.Name)
	assert.True(t, secret.Secret.Value)
	assert.Empty(t, secret.New, "a secret keeps its name and no value")
	assert.True(t, page.Items[0].Reads.Null)

	page, err = listTrail(h, api.AdminListAuditRecordsParams{Operation: api.NewOptString("AdminListUsers")})
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	assert.Empty(t, page.Items[0].Changes)
	assert.Equal(t, int32(2), page.Items[0].Reads.Value.Count)
	assert.Equal(t, []string{"1", "2"}, page.Items[0].Reads.Value.Ids)
}

// Reading the trail is itself a read of personal data (ADR-0020): the call
// leaves a record of AuditRecord rows being read once the worker has run, and
// the ids are those of the records the page showed.
func TestReadingTheTrailIsRecordedAsAReadOfPersonalData(t *testing.T) {
	t.Parallel()
	h := auditReader(t, true)
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
	t.Parallel()
	h := testsupport.NewHarness(t)

	_, err := listTrail(h, api.AdminListAuditRecordsParams{})
	require.Error(t, err)
	processReads(t, h)

	assert.Empty(t, testsupport.PublishedEvents[events.StaffReadPersonalData](t, h.Services))
}

func capabilitiesOf(t *testing.T, h *testsupport.Harness) []string {
	t.Helper()
	res, err := h.Client.AdminGetMyCapabilities(context.Background())
	require.NoError(t, err)
	held := res
	return held.Held
}

// The menu asks what the caller holds: the administrator flag adds nothing, the
// grant on the role is the only way in. Each case has its own test because a
// harness holds an open transaction, and two of them inserting a staff record
// for the same user would wait on each other.
func TestAnAdministratorHoldsNoCapabilityByBeingOne(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	assert.Empty(t, capabilitiesOf(t, h))
}

func TestARoleWithTheGrantHoldsTheCapability(t *testing.T) {
	t.Parallel()
	h := auditReader(t, true)

	assert.Equal(t, []string{"audit.read"}, capabilitiesOf(t, h))
}

func TestARoleWithoutTheGrantHoldsNoCapability(t *testing.T) {
	t.Parallel()
	h := auditReader(t, false)

	assert.Empty(t, capabilitiesOf(t, h))
}

func TestAnyStaffMemberMayAskWhatTheyHold(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)

	assert.Empty(t, capabilitiesOf(t, h))
}
