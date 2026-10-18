package mcpserver

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/oauthserver"
)

func TestEveryClassHasItsOwnConfigurableBudgetPerUserAndPerClient(t *testing.T) {
	now := time.Now()
	m := newMeter(Limits{
		RateWindow: time.Minute,
		User:       Rates{Catalogue: 3, Personal: 2, Write: 1},
		Client:     Rates{Catalogue: 6, Personal: 5, Write: 4},
		Now:        func() time.Time { return now },
	}.withDefaults())

	for class, userLimit := range map[Class]int{ClassCatalogue: 3, ClassPersonal: 2, ClassWrite: 1} {
		for range userLimit {
			ok, _ := m.allow(class, 1, "app")
			require.True(t, ok, class)
		}

		ok, wait := m.allow(class, 1, "app")
		assert.False(t, ok, "%s: the user's own limit applies", class)
		assert.Equal(t, time.Minute, wait, class)
	}

	// Each class counted on its own: exhausting one left another's budget alone,
	// and another user of the same app is unaffected by the first's.
	ok, _ := m.allow(ClassCatalogue, 2, "app")
	assert.True(t, ok)
}

func TestClientLimitOfAClassStopsAnAppWhateverTheUser(t *testing.T) {
	now := time.Now()
	m := newMeter(Limits{
		User:   Rates{Catalogue: 100, Personal: 100, Write: 100},
		Client: Rates{Catalogue: 100, Personal: 100, Write: 2},
		Now:    func() time.Time { return now },
	}.withDefaults())

	first, _ := m.allow(ClassWrite, 1, "app")
	second, _ := m.allow(ClassWrite, 2, "app")
	third, _ := m.allow(ClassWrite, 3, "app")
	otherApp, _ := m.allow(ClassWrite, 3, "other")
	reads, _ := m.allow(ClassCatalogue, 3, "app")

	assert.True(t, first)
	assert.True(t, second)
	assert.False(t, third)
	assert.True(t, otherApp)
	assert.True(t, reads, "a write limit does not touch reads")
}

func TestDefaultLimitsAreStricterForPersonalReadsAndWrites(t *testing.T) {
	d := DefaultLimits()

	assert.Equal(t, 50, d.MaxPageSize)
	assert.Equal(t, 100*1024, d.MaxResponseBytes)
	assert.Equal(t, 30*time.Second, d.ToolTimeout)

	for _, r := range []Rates{d.User, d.Client} {
		assert.Greater(t, r.Catalogue, r.Personal)
		assert.Greater(t, r.Personal, r.Write)
	}

	filled := Limits{MaxPageSize: 3, User: Rates{Write: 1}}.withDefaults()
	assert.Equal(t, 3, filled.MaxPageSize, "a value that is set stays")
	assert.Equal(t, 1, filled.User.Write)
	assert.Equal(t, d.User.Catalogue, filled.User.Catalogue, "a value that is not set takes its default")
	assert.Equal(t, d.ToolTimeout, filled.ToolTimeout)
}

func requestFor(userID, client string, scopes ...string) *mcp.CallToolRequest {
	return &mcp.CallToolRequest{Extra: &mcp.RequestExtra{TokenInfo: &auth.TokenInfo{
		UserID: userID, Scopes: scopes, Extra: map[string]any{"client_id": client},
	}}}
}

func TestAdmitAppliesTheRateLimitOfTheClassTheToolNames(t *testing.T) {
	now := time.Now()
	s := &Server{limits: Limits{Now: func() time.Time { return now }}.withDefaults()}
	s.limits.User = Rates{Catalogue: 3, Personal: 2, Write: 1}
	s.meter = newMeter(s.limits)
	req := requestFor("7", "app", oauthserver.ScopeCatalogRead)

	for class, limit := range map[Class]int{ClassCatalogue: 3, ClassPersonal: 2, ClassWrite: 1} {
		for range limit {
			userID, client, denied := s.admit(req, oauthserver.ScopeCatalogRead, class)
			require.Nil(t, denied, class)
			assert.Equal(t, 7, userID)
			assert.Equal(t, "app", client)
		}

		_, _, denied := s.admit(req, oauthserver.ScopeCatalogRead, class)
		require.NotNil(t, denied, class)
		assert.True(t, denied.IsError)
		text := denied.Content[0].(*mcp.TextContent).Text
		assert.Contains(t, text, "rate limit exceeded for "+class.String()+" tools")
		assert.Contains(t, text, "retry in 60 seconds")
		assert.Equal(t, 60, denied.Meta["retryAfterSeconds"])
		assert.Equal(t, class.String(), denied.Meta["rateLimitClass"])
	}
}

func TestAdmitRefusesAMissingScopeOrSubjectBeforeCounting(t *testing.T) {
	s := &Server{limits: DefaultLimits()}
	s.meter = newMeter(s.limits)

	for name, req := range map[string]*mcp.CallToolRequest{
		"no token":    {Extra: &mcp.RequestExtra{}},
		"not a user":  requestFor("client-3", "app", oauthserver.ScopeCatalogRead),
		"wrong scope": requestFor("7", "app", oauthserver.ScopeProgressRead),
	} {
		_, _, denied := s.admit(req, oauthserver.ScopeCatalogRead, ClassCatalogue)
		require.NotNil(t, denied, name)
		assert.True(t, denied.IsError, name)
	}
}

func TestClassNamesAreStableAndUnknownClassIsNamedNotHidden(t *testing.T) {
	assert.Equal(t, "catalogue", ClassCatalogue.String())
	assert.Equal(t, "personal", ClassPersonal.String())
	assert.Equal(t, "write", ClassWrite.String())
	assert.Equal(t, "class9", Class(9).String())
	assert.False(t, annotations(ClassWrite).ReadOnlyHint)
	assert.True(t, annotations(ClassPersonal).ReadOnlyHint)
}

// sizeOf is the size fit measures: the JSON of the data and the text.
func sizeOf(t *testing.T, data map[string]any, textLen int) int {
	t.Helper()

	raw, err := json.Marshal(data)
	require.NoError(t, err)

	return len(raw) + textLen
}

func records(n int) []any {
	out := make([]any, n)
	for i := range out {
		out[i] = map[string]any{"slug": strings.Repeat("x", 100), "i": i}
	}

	return out
}

func TestFitLeavesAResultWithinTheLimitAlone(t *testing.T) {
	data := map[string]any{"items": records(2)}

	assert.True(t, fit(data, 10, 100_000, nil))

	assert.NotContains(t, data, truncationKey)
	assert.Len(t, data["items"], 2)
}

func TestFitCutsTheLargestListAndSaysHowToContinueFromTheNextPage(t *testing.T) {
	data := map[string]any{"items": records(100), "small": records(2), "total": 100}

	ok := fit(data, 20, 4000, &listInfo{key: "items", page: 3, perPage: 100})

	require.True(t, ok)
	assert.LessOrEqual(t, sizeOf(t, data, 20), 4000)
	assert.Len(t, data["small"], 2, "the small list was not touched")
	kept := len(data["items"].([]any))
	assert.Less(t, kept, 100)

	hint := data[truncationKey].(map[string]any)
	assert.Equal(t, true, hint["truncated"])
	assert.Equal(t, kept, hint["suggestedPerPage"])
	// Page 3 of 100 starts at offset 200, the page of that size containing it.
	assert.Equal(t, 200/kept+1, hint["suggestedPage"])
	assert.Contains(t, hint["hint"], "perPage=")
}

func TestFitShortensALongTextOnARuneBoundary(t *testing.T) {
	text := strings.Repeat("Привет, мир! ", 1000)
	data := map[string]any{"lesson": map[string]any{"theory": text, "slug": "a"}}

	ok := fit(data, 0, 3000, nil)

	require.True(t, ok)

	theory := data["lesson"].(map[string]any)["theory"].(string)
	assert.True(t, utf8.ValidString(theory), "a cut must not split a character")
	assert.True(t, strings.HasSuffix(theory, "[cut]"))
	assert.LessOrEqual(t, sizeOf(t, data, 0), 3000)
	hint := data[truncationKey].(map[string]any)
	cuts := hint["cuts"].([]any)
	require.Len(t, cuts, 1)
	assert.Equal(t, "lesson.theory", cuts[0].(map[string]any)["path"])
	assert.Equal(t, "bytes", cuts[0].(map[string]any)["unit"])
	assert.Contains(t, hint["hint"], "lesson.theory")
}

func TestFitCutsListsInsideARecordAndReportsEachPlace(t *testing.T) {
	data := map[string]any{"lessons": records(60), "modules": records(60), "name": "course"}

	ok := fit(data, 0, 5000, nil)

	require.True(t, ok)
	assert.LessOrEqual(t, sizeOf(t, data, 0), 5000)
	hint := data[truncationKey].(map[string]any)
	assert.NotEmpty(t, hint["cuts"])
	assert.Contains(t, hint["hint"], "shortened")
	assert.NotContains(t, hint, "suggestedPage", "not a page of a list")
}

func TestFitGivesUpWhenNothingCanBeShortened(t *testing.T) {
	data := map[string]any{}
	for i := range 200 {
		data[strings.Repeat("k", 10)+string(rune('A'+i%26))+strings.Repeat("v", i)] = 1
	}

	assert.False(t, fit(data, 0, 2000, nil))
}

func TestWithTruncationDeclaresTheHintWithoutChangingTheContractSchema(t *testing.T) {
	contract := map[string]any{
		"type":                 "object",
		"properties":           map[string]any{"items": map[string]any{"type": "array"}},
		"additionalProperties": false,
	}

	out := withTruncation(contract)

	assert.Contains(t, out["properties"], truncationKey)
	assert.Contains(t, out["properties"], "items")
	assert.NotContains(t, contract["properties"], truncationKey, "the original schema is not edited")
	assert.Equal(t, false, out["additionalProperties"])
}

func TestAllowNullKeepsEveryKindOfNullableSchemaValid(t *testing.T) {
	assert.Equal(t, map[string]any{"type": []any{"string", "null"}}, allowNull(map[string]any{"type": "string"}))

	enum := allowNull(map[string]any{"type": "string", "enum": []any{"a"}})
	assert.Equal(t, []any{"a", nil}, enum["enum"])

	wrapped := allowNull(map[string]any{"type": "object", "allOf": []any{map[string]any{"type": "object"}}})
	assert.Contains(t, wrapped, "anyOf", "a schema made of others takes null as an alternative")
}
