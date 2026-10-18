package mcpserver_test

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/auditrecord"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/eventhandler"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// The audit trail is the same operation over MCP as on the admin page
// (adminListAuditRecords), dispatched through the same generated server, so the
// page's rules (the explicit audit.read grant, the window, the order) apply to an
// assistant unchanged. These tests cover what is specific to the tool.

// trailReader connects an administrator who is also staff (user_admin_staff, on
// the editor role) with the admin Scope and the given grants on the role.
func trailReader(t *testing.T, o *testsupport.OAuth, grants ...string) (*mcp.ClientSession, oauthserver.Token) {
	t.Helper()
	o.Grant(editor, grants...)
	return admin(t, o, adminStaff)
}

// seedTrail writes one Audit Record through the production writer, placed in the
// past so that its position in a newest-first listing is known.
func seedTrail(t *testing.T, o *testsupport.OAuth, ago time.Duration, operation string, mutate func(*audit.Entry)) int {
	t.Helper()
	entry := audit.Entry{
		OccurredAt: time.Now().Add(-ago),
		Source:     audit.Source{Channel: audit.ChannelWeb},
		Operation:  operation,
		Outcome:    audit.OutcomeSucceeded,
	}
	if mutate != nil {
		mutate(&entry)
	}
	require.NoError(t, audit.Write(t.Context(), o.DB, entry, nil))
	return o.DB.AuditRecord.Query().Order(ent.Desc(auditrecord.FieldID)).FirstX(t.Context()).ID
}

func itemIDs(t *testing.T, page map[string]any) []int {
	t.Helper()
	listed := page["items"].([]any)
	ids := make([]int, 0, len(listed))
	for _, item := range listed {
		ids = append(ids, int(item.(map[string]any)["id"].(float64)))
	}
	return ids
}

func TestTrailToolListsNewestFirstPaginatesAndFilters(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session, _ := trailReader(t, o, "audit.read", "mcp.read")
	oldest := seedTrail(t, o, 72*time.Hour, "AdminCreateBanner", nil)
	middle := seedTrail(t, o, 48*time.Hour, "AdminUpdateBanner", func(e *audit.Entry) {
		e.Source.Channel = audit.ChannelMCP
	})
	newest := seedTrail(t, o, 24*time.Hour, "AdminDeleteBanner", func(e *audit.Entry) {
		e.EntityType, e.EntityID = "Banner", "77"
	})

	first := structured(t, call(t, session, "adminListAuditRecords", map[string]any{"locale": "en", "perPage": 2}))
	assert.Equal(t, []int{newest, middle}, itemIDs(t, first))
	assert.EqualValues(t, 2, first["perPage"])

	second := structured(t, call(t, session, "adminListAuditRecords", map[string]any{"locale": "en", "perPage": 2, "page": 2}))
	require.NotEmpty(t, itemIDs(t, second))
	assert.Equal(t, oldest, itemIDs(t, second)[0])

	byChannel := structured(t, call(t, session, "adminListAuditRecords", map[string]any{"locale": "en", "channel": "mcp", "operation": "AdminUpdateBanner"}))
	assert.Equal(t, []int{middle}, itemIDs(t, byChannel))

	byEntity := structured(t, call(t, session, "adminListAuditRecords", map[string]any{"locale": "en", "entityType": "Banner", "entityId": "77"}))
	assert.Equal(t, []int{newest}, itemIDs(t, byEntity))

	from := time.Now().Add(-60 * time.Hour).UTC().Format(time.RFC3339)
	to := time.Now().Add(-36 * time.Hour).UTC().Format(time.RFC3339)
	window := structured(t, call(t, session, "adminListAuditRecords", map[string]any{"locale": "en", "from": from, "to": to}))
	assert.Equal(t, []int{middle}, itemIDs(t, window))

	actor := seedTrail(t, o, 12*time.Hour, "AdminUpdateBanner", func(e *audit.Entry) {
		id := bob
		e.Actor = &id
	})
	byActor := structured(t, call(t, session, "adminListAuditRecords", map[string]any{"locale": "en", "actorUserId": bob}))
	assert.Equal(t, []int{actor}, itemIDs(t, byActor))
}

func TestTrailToolWindowRulesAreThePagesRules(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session, _ := trailReader(t, o, "audit.read", "mcp.read")
	seedTrail(t, o, 40*24*time.Hour, "AdminCreateBanner", nil)

	tooWide := call(t, session, "adminListAuditRecords", map[string]any{
		"locale": "en", "from": time.Now().Add(-93 * 24 * time.Hour).UTC().Format(time.RFC3339),
	})
	require.True(t, tooWide.IsError)
	assert.Contains(t, textOf(t, tooWide), "HTTP 422")
	assert.Contains(t, textOf(t, tooWide), "from")

	reversed := call(t, session, "adminListAuditRecords", map[string]any{
		"locale": "en",
		"from":   time.Now().Add(-time.Hour).UTC().Format(time.RFC3339),
		"to":     time.Now().Add(-2 * time.Hour).UTC().Format(time.RFC3339),
	})
	require.True(t, reversed.IsError)
	assert.Contains(t, textOf(t, reversed), "HTTP 422")

	orphan := call(t, session, "adminListAuditRecords", map[string]any{"locale": "en", "entityId": "7"})
	require.True(t, orphan.IsError)
	assert.Contains(t, textOf(t, orphan), "HTTP 422")

	// The default window is 30 days: the record of 40 days ago is outside it, and
	// a narrower explicit window reaches back to it.
	assert.Empty(t, structured(t, call(t, session, "adminListAuditRecords", map[string]any{"locale": "en"}))["items"])
	back := structured(t, call(t, session, "adminListAuditRecords", map[string]any{
		"locale": "en", "from": time.Now().Add(-45 * 24 * time.Hour).UTC().Format(time.RFC3339),
	}))
	assert.NotEmpty(t, back["items"])
}

func TestTrailToolNeedsAuditReadAndMCPReadAndTheAdminScope(t *testing.T) {
	args := map[string]any{"locale": "en"}

	t.Run("administrator without audit.read", func(t *testing.T) {
		o := testsupport.NewOAuth(t)
		session, token := trailReader(t, o, "mcp.read")

		assertRefused(t, call(t, session, "adminListAuditRecords", args))

		records := o.Records("AdminListAuditRecords")
		require.Len(t, records, 1)
		assert.Equal(t, "forbidden", records[0].Outcome)
		assert.Equal(t, "mcp", records[0].Channel)
		assert.Equal(t, token.ID, *records[0].TokenID)
	})

	t.Run("audit.read without mcp.read", func(t *testing.T) {
		o := testsupport.NewOAuth(t)
		session, _ := trailReader(t, o, "audit.read")

		assertRefused(t, call(t, session, "adminListAuditRecords", args))

		records := o.Records("AdminListAuditRecords")
		require.Len(t, records, 1)
		assert.Equal(t, "forbidden", records[0].Outcome)
	})

	t.Run("an administrator with no staff record", func(t *testing.T) {
		o := testsupport.NewOAuth(t)
		o.Grant(editor, "audit.read", "mcp.read")
		session, _ := admin(t, o, alice)

		assertRefused(t, call(t, session, "adminListAuditRecords", args))
	})

	t.Run("a token without the admin Scope", func(t *testing.T) {
		o := testsupport.NewOAuth(t)
		o.Grant(editor, "audit.read", "mcp.read")
		tokens := o.Connect(adminStaff, oauthserver.ScopeCatalogRead)
		session := o.MCPClient(tokens.AccessToken)

		assertRefused(t, call(t, session, "adminListAuditRecords", args))

		records := o.Records("AdminListAuditRecords")
		require.Len(t, records, 1)
		assert.Equal(t, "forbidden", records[0].Outcome)
	})

	t.Run("both grants", func(t *testing.T) {
		o := testsupport.NewOAuth(t)
		session, _ := trailReader(t, o, "audit.read", "mcp.read")

		result := call(t, session, "adminListAuditRecords", args)

		assert.False(t, result.IsError, "%+v", result.Content)
	})
}

// Reading the trail is a read of personal data: the records carry addresses and
// old values. Over MCP it is recorded like the page's read, on the mcp channel
// with the Connected App and the token.
func TestReadingTheTrailOverMCPIsRecordedAsAReadOfPersonalData(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session, token := trailReader(t, o, "audit.read", "mcp.read")
	seeded := seedTrail(t, o, time.Hour, "AdminUpdateBanner", nil)

	page := structured(t, call(t, session, "adminListAuditRecords", map[string]any{"locale": "ru"}))
	require.Contains(t, itemIDs(t, page), seeded)

	// Recorded after the request, by the consumer of the event.
	assert.Empty(t, o.Records("AdminListAuditRecords"))
	testsupport.Deliver(t, o.Services, eventhandler.StaffReadPersonalData(o.DB))

	records := o.Records("AdminListAuditRecords")
	require.Len(t, records, 1)
	record := records[0]
	assert.Equal(t, "mcp", record.Channel)
	assert.Equal(t, adminStaff, *record.ActorUserID)
	assert.Equal(t, "succeeded", record.Outcome)
	assert.Equal(t, "AuditRecord", *record.EntityType)
	assert.Equal(t, "ru", *record.Locale)
	assert.Equal(t, token.ConnectedAppID, *record.ConnectedAppID)
	assert.Equal(t, token.ID, *record.TokenID)

	var reads struct {
		Count int      `json:"count"`
		IDs   []string `json:"ids"`
	}
	require.NoError(t, json.Unmarshal(record.Reads, &reads))
	assert.Equal(t, len(itemIDs(t, page)), reads.Count)
	assert.Contains(t, reads.IDs, strconv.Itoa(seeded))
}

func TestReadingTheTrailSpendsThePersonalDataRateClass(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{
		User: mcpserver.Rates{Personal: 1, Write: 100},
	}))
	session, _ := trailReader(t, o, "audit.read", "mcp.read")
	args := map[string]any{"locale": "en"}

	require.False(t, call(t, session, "adminListAuditRecords", args).IsError)
	rejected(t, call(t, session, "adminListAuditRecords", args))
}

func TestTrailToolDeclaresTheFiltersAndARequiredLocale(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session, _ := trailReader(t, o, "audit.read", "mcp.read")
	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	var tool *mcp.Tool
	for _, candidate := range tools.Tools {
		if candidate.Name == "adminListAuditRecords" {
			tool = candidate
		}
	}
	require.NotNil(t, tool)
	assert.True(t, tool.Annotations.ReadOnlyHint)

	raw, err := json.Marshal(tool.InputSchema)
	require.NoError(t, err)
	var input struct {
		Required   []string       `json:"required"`
		Properties map[string]any `json:"properties"`
	}
	require.NoError(t, json.Unmarshal(raw, &input))
	assert.Equal(t, []string{"locale"}, input.Required, "only the locale is required")
	for _, name := range []string{"actorUserId", "channel", "operation", "entityType", "entityId", "from", "to", "page", "perPage"} {
		assert.Contains(t, input.Properties, name)
	}
	assert.NotContains(t, input.Properties, "body")
}

// The rows carry `changes` JSON, so a long listing is cut by whole records and
// the answer stays valid JSON that still matches the output schema.
func TestTrailToolResponseOverTheLimitIsCutAndStaysValid(t *testing.T) {
	const limit = 6000
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{MaxResponseBytes: limit}))
	session, _ := trailReader(t, o, "audit.read", "mcp.read")
	for i := range 40 {
		seedTrail(t, o, time.Duration(i+1)*time.Minute, "AdminUpdateBanner", func(e *audit.Entry) {
			e.UserAgent = "padding-padding-padding-padding-padding-padding-padding-padding"
		})
	}
	schemas := outputSchemas(t, session)

	result := call(t, session, "adminListAuditRecords", map[string]any{"locale": "en", "perPage": 100})

	content := structured(t, result)
	raw, err := json.Marshal(content)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(raw)+len(textOf(t, result)), limit)
	require.NoError(t, schemas["adminListAuditRecords"].Validate(content))
	hint, ok := content["truncation"].(map[string]any)
	require.True(t, ok, "the cut is announced")
	assert.Equal(t, true, hint["truncated"])
	assert.Contains(t, textOf(t, result), "truncation")
	assert.Less(t, len(content["items"].([]any)), 40)
	assert.GreaterOrEqual(t, content["total"], float64(40), "the total is not cut")
}

// A secret field is recorded by name only. A change of a password digest, listed
// through the tool, shows the field and says it is a secret, with no value on
// either side, and the digest is nowhere in the answer.
func TestTrailToolShowsASecretFieldByNameWithoutValues(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session, _ := trailReader(t, o, "audit.read", "mcp.read")

	const digest = "$2a$10$brand-new-secret-digest-value"
	ctx, recorder := audit.WithRecorder(t.Context())
	user := fixture.ID(fixture.UserWithPassword)
	require.NoError(t, o.DB.User.UpdateOneID(user).SetPasswordDigest(digest).SetFirstName("Renamed").Exec(ctx))
	require.NoError(t, audit.Write(t.Context(), o.DB, audit.Entry{
		Source:    audit.Source{Channel: audit.ChannelWeb},
		Operation: "AdminUpdateUser",
		Outcome:   audit.OutcomeSucceeded,
	}, recorder))

	result := call(t, session, "adminListAuditRecords", map[string]any{"locale": "en", "operation": "AdminUpdateUser"})

	content := structured(t, result)
	items := content["items"].([]any)
	require.Len(t, items, 1)
	changes := items[0].(map[string]any)["changes"].([]any)
	require.Len(t, changes, 1)

	fields := map[string]map[string]any{}
	for _, f := range changes[0].(map[string]any)["fields"].([]any) {
		field := f.(map[string]any)
		fields[field["name"].(string)] = field
	}
	secret, ok := fields["password_digest"]
	require.True(t, ok, "a changed secret is still listed")
	assert.Equal(t, true, secret["secret"])
	assert.NotContains(t, secret, "old")
	assert.NotContains(t, secret, "new")
	assert.Equal(t, "Renamed", fields["first_name"]["new"], "an ordinary field keeps its values")

	raw, err := json.Marshal(result)
	require.NoError(t, err)
	assert.NotContains(t, string(raw), digest)
}
