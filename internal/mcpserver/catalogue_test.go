package mcpserver_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/fixture"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

func bob() int {
	return fixture.ID(fixture.UserStaffEditor)
}

const (
	// jsLanding is the landing page slug of the JavaScript course in Russian, the
	// one getCourse takes; jsCourse is the course's own slug, which
	// getCourseLesson takes.
	jsLanding   = "javascript-ru"
	jsCourse    = "javascript"
	firstLesson = "hello-world"
)

// russian makes the connection speak Russian: the fixtures' course content is.
func russian(r *http.Request) { r.Header.Set("Accept-Language", "ru") }

// callTool calls a tool and fails the test on a protocol error.
func callTool(t *testing.T, session *mcp.ClientSession, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	result, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: name, Arguments: args})
	require.NoError(t, err)

	return result
}

// structured is the structured content of a successful call.
func structured(t *testing.T, result *mcp.CallToolResult) map[string]any {
	t.Helper()
	require.False(t, result.IsError, "%+v", result.Content)
	content, ok := result.StructuredContent.(map[string]any)
	require.True(t, ok, "%T", result.StructuredContent)

	return content
}

// textOf is the text of a result, which every result has.
func textOf(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	require.Len(t, result.Content, 1)
	text, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)

	return text.Text
}

// outputSchemas indexes the advertised output schemas by tool name.
func outputSchemas(t *testing.T, session *mcp.ClientSession) map[string]*jsonschema.Resolved {
	t.Helper()
	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	schemas := map[string]*jsonschema.Resolved{}

	for _, tool := range tools.Tools {
		raw, err := json.Marshal(tool.OutputSchema)
		require.NoError(t, err)

		var schema jsonschema.Schema
		require.NoError(t, json.Unmarshal(raw, &schema))
		resolved, err := schema.Resolve(nil)
		require.NoError(t, err)

		schemas[tool.Name] = resolved
	}

	return schemas
}

// courseCount is how many courses the catalogue has in the connection's locale.
func courseCount(t *testing.T, session *mcp.ClientSession) int {
	t.Helper()

	return int(as[float64](t, structured(t, callTool(t, session, "listCourses", map[string]any{"perPage": 1}))["total"]))
}

func connect(t *testing.T, o *testsupport.OAuth, userID int, scopes ...string) *mcp.ClientSession {
	t.Helper()

	tokens := o.Connect(userID, scopes...)

	return o.MCPClient(tokens.AccessToken, russian)
}

func TestEveryCatalogueToolAnswersWithStructuredContentThatMatchesItsContractSchema(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)
	schemas := outputSchemas(t, session)

	categories := as[[]any](t, structured(t, callTool(t, session, "listPublicCourseCategories", nil))["categories"])
	require.NotEmpty(t, categories)
	categorySlug := as[string](t, as[map[string]any](t, categories[0])["slug"])

	for name, args := range map[string]map[string]any{
		"listCourses":                {},
		"listPublicCourseCategories": {},
		"getPublicCourseCategory":    {"slug": categorySlug},
		"getCourse":                  {"slug": jsLanding},
		"getCourseLesson":            {"courseSlug": jsCourse, "slug": firstLesson},
		"listBlogPosts":              {},
		"getBlogPost":                {"slug": "learning-ruby"},
	} {
		t.Run(name, func(t *testing.T) {
			result := callTool(t, session, name, args)

			content := structured(t, result)
			require.Contains(t, schemas, name)
			require.NoError(t, schemas[name].Validate(content), "the content must match the advertised schema")
			assert.NotEmpty(t, textOf(t, result), "a short text accompanies the content")
		})
	}
}

func TestEveryCatalogueToolIsReadOnlyAndNamedAfterItsOperation(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	var names []string

	for _, tool := range tools.Tools {
		if strings.HasPrefix(tool.Name, "admin") || isLearnerTool(tool.Name) || tool.Name == "findTools" {
			continue // the admin tools are tested in admin_test.go, the lookup in discovery_test.go
		}

		names = append(names, tool.Name)
		assert.True(t, tool.Annotations.ReadOnlyHint, tool.Name)
	}

	assert.ElementsMatch(t, []string{
		"listCourses", "listPublicCourseCategories", "getPublicCourseCategory",
		"getCourse", "getCourseLesson", "listBlogPosts", "getBlogPost",
	}, names)
}

func TestCourseToolReturnsTheCourseWithItsLessonsAsTheHTTPOperationDoes(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	content := structured(t, callTool(t, session, "getCourse", map[string]any{"slug": jsLanding}))

	res := o.Do(http.MethodGet, "/api/languages/"+jsLanding, nil, russian)
	defer func() { _ = res.Body.Close() }()

	var visitor map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&visitor))
	assert.Equal(t, visitor, content, "a tool and the operation cannot answer differently")
	assert.NotEmpty(t, content["lessons"])
}

func TestLessonToolReturnsTheContentAndTestsTheHTTPOperationGivesAVisitor(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	content := structured(t, callTool(t, session, "getCourseLesson", map[string]any{"courseSlug": jsCourse, "slug": firstLesson}))

	res := o.Do(http.MethodGet, "/api/languages/"+jsCourse+"/lessons/"+firstLesson, nil, russian)
	defer func() { _ = res.Body.Close() }()

	var visitor map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&visitor))
	assert.Equal(t, visitor, content, "the same fields, the same values, for the same visitor")

	lesson := as[map[string]any](t, content["lesson"])
	assert.Equal(t, firstLesson, lesson["slug"])
	assert.Contains(t, lesson, "testCode", "the lesson tests are part of the lesson content")
	assert.Contains(t, lesson, "preparedCode", "the starter code follows the HTTP operation's rule for visitors")
	assert.Contains(t, lesson, "originalCode", "so does the reference solution")
}

func TestAnsweringAsAVisitorNeverCarriesWhatOnlyASignedInLearnerSees(t *testing.T) {
	// Alice is enrolled in the course and has finished its first lesson. The
	// token identifies her to the endpoint, not to the handler: the handler sees
	// no session, so her Enrollment and progress must not appear.
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	course := structured(t, callTool(t, session, "getCourse", map[string]any{"slug": jsLanding}))
	assert.Nil(t, course["enrollment"], "no Enrollment for a visitor")
	progress, _ := course["progress"].(map[string]any)
	assert.NotContains(t, progress, "enrollment")
	assert.EqualValues(t, 0, progress["completion"], "a visitor stands at the beginning, whatever the token's user finished")

	lesson := structured(t, callTool(t, session, "getCourseLesson", map[string]any{"courseSlug": jsCourse, "slug": firstLesson}))
	lessonProgress, _ := lesson["progress"].(map[string]any)
	assert.EqualValues(t, 0, lessonProgress["completion"])

	raw, err := json.Marshal(lesson)
	require.NoError(t, err)
	assert.NotContains(t, string(raw), `"finished":true`, "no lesson reads as finished for a visitor")
}

func TestOnlyPublishedPostsOfTheLocaleAreReadable(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	posts := structured(t, callTool(t, session, "listBlogPosts", nil))
	listed := as[[]any](t, posts["items"])

	slugs := make([]string, 0, len(listed))
	for _, p := range listed {
		slugs = append(slugs, as[string](t, as[map[string]any](t, p)["slug"]))
	}

	assert.ElementsMatch(t, []string{"hello-world", "learning-ruby", "oldest-post"}, slugs, "no draft, archived or English post")

	for _, slug := range []string{"second-post", "archived-notes", "english-post"} {
		result := callTool(t, session, "getBlogPost", map[string]any{"slug": slug})
		assert.True(t, result.IsError, slug)
		assert.Contains(t, textOf(t, result), "not found", slug)
	}
}

func TestAMissingRecordIsAToolErrorTheAgentCanReadNotAFailedRequest(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	for name, args := range map[string]map[string]any{
		"getPublicCourseCategory": {"slug": "no-such-category"},
		"getCourse":               {"slug": "no-such-course"},
		"getCourseLesson":         {"courseSlug": jsCourse, "slug": "no-such-lesson"},
		"getBlogPost":             {"slug": "no-such-post"},
	} {
		t.Run(name, func(t *testing.T) {
			result := callTool(t, session, name, args)

			assert.True(t, result.IsError)
			assert.Contains(t, textOf(t, result), "not found")
			assert.Nil(t, result.StructuredContent)
		})
	}
}

func TestToolRefusesArgumentsItDoesNotKnow(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	for name, args := range map[string]map[string]any{
		"listCourses": {"per_page": 5},
		"getCourse":   {"slug": 5},
	} {
		result := callTool(t, session, name, args)

		assert.True(t, result.IsError, name)
		assert.Contains(t, textOf(t, result), "invalid arguments", name)
	}
}

func TestEveryCatalogueToolRefusesAConnectionWithoutTheCatalogReadScope(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice(), oauthserver.ScopeProgressRead)

	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	for _, tool := range tools.Tools {
		if strings.HasPrefix(tool.Name, "admin") || isLearnerTool(tool.Name) || tool.Name == "findTools" {
			continue
		}

		result := callTool(t, session, tool.Name, map[string]any{"slug": "x", "courseSlug": "x"})

		assert.True(t, result.IsError, tool.Name)
		assert.Contains(t, textOf(t, result), oauthserver.ScopeCatalogRead, tool.Name)
		assert.Nil(t, result.StructuredContent, tool.Name)
	}
}

func TestPageSizeAboveTheLimitIsClampedAndSaidSo(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{MaxPageSize: 2}))
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	t.Run("a list the tool pages itself", func(t *testing.T) {
		result := callTool(t, session, "listCourses", map[string]any{"perPage": 500})

		content := structured(t, result)
		assert.Len(t, content["courses"], 2)
		assert.EqualValues(t, 2, content["perPage"])
		assert.EqualValues(t, courseCount(t, session), content["total"], "the total still says how many there are")
		assert.Contains(t, textOf(t, result), "clamped to 2")
	})

	t.Run("a list the handler pages", func(t *testing.T) {
		result := callTool(t, session, "listBlogPosts", map[string]any{"perPage": 100})

		content := structured(t, result)
		assert.Len(t, content["items"], 2)
		assert.EqualValues(t, 3, content["total"])
		assert.Contains(t, textOf(t, result), "clamped to 2")
	})

	t.Run("the default page is the limit", func(t *testing.T) {
		content := structured(t, callTool(t, session, "listPublicCourseCategories", nil))

		assert.Len(t, content["categories"], 2)
	})
}

func TestPagesFollowEachOtherWithoutGapsAndPastTheEndIsEmpty(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	var slugs []string

	total := courseCount(t, session)
	for page := 1; page <= (total+1)/2; page++ {
		content := structured(t, callTool(t, session, "listCourses", map[string]any{"page": page, "perPage": 2}))
		for _, c := range as[[]any](t, content["courses"]) {
			slugs = append(slugs, as[string](t, as[map[string]any](t, c)["slug"]))
		}
	}

	all := as[[]any](t, structured(t, callTool(t, session, "listCourses", nil))["courses"])
	require.Len(t, all, total)

	want := make([]string, 0, len(all))
	for _, c := range all {
		want = append(want, as[string](t, as[map[string]any](t, c)["slug"]))
	}

	assert.Equal(t, want, slugs)

	past := structured(t, callTool(t, session, "listCourses", map[string]any{"page": 99}))
	assert.Empty(t, past["courses"])

	for _, args := range []map[string]any{{"page": -1}, {"perPage": -3}} {
		assert.True(t, callTool(t, session, "listCourses", args).IsError, "%v", args)
	}
}

func TestResponseOverTheLimitIsCutWithAPaginationHintAndStaysValid(t *testing.T) {
	const limit = 3000

	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{MaxResponseBytes: limit}))
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)
	schemas := outputSchemas(t, session)

	result := callTool(t, session, "listCourses", nil)

	content := structured(t, result)
	raw, err := json.Marshal(content)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(raw)+len(textOf(t, result)), limit)
	require.NoError(t, schemas["listCourses"].Validate(content), "a cut result still matches the schema")

	hint, ok := content["truncation"].(map[string]any)
	require.True(t, ok, "the cut is announced inside the structured content")
	assert.Equal(t, true, hint["truncated"])
	assert.EqualValues(t, limit, hint["limitBytes"])
	assert.Contains(t, textOf(t, result), "truncation")

	kept := len(as[[]any](t, content["courses"]))
	total := courseCount(t, session)
	assert.Less(t, kept, total)
	assert.EqualValues(t, total, content["total"], "the total is not cut")
	assert.EqualValues(t, kept, hint["suggestedPerPage"])

	// Following the hint reads the rest, and the pages together hold every course.
	var slugs []string

	perPage := int(as[float64](t, hint["suggestedPerPage"]))
	for page := 1; len(slugs) < total && page < 20; page++ {
		next := structured(t, callTool(t, session, "listCourses", map[string]any{"page": page, "perPage": perPage}))
		for _, c := range as[[]any](t, next["courses"]) {
			slugs = append(slugs, as[string](t, as[map[string]any](t, c)["slug"]))
		}
	}

	assert.Len(t, slugs, total)
}

func TestResponseUnderTheLimitCarriesNoHint(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	content := structured(t, callTool(t, session, "listCourses", nil))

	assert.NotContains(t, content, "truncation")
}

// The cutting of records is covered by the unit tests of fit, where a record can
// be made as long as needed; here, a fixture lesson has no part that can be
// shortened, so a limit under its size must be a clear refusal naming the limit,
// never a half-cut or invalid result.
func TestARecordThatCannotBeShortenedToTheLimitIsRefusedByName(t *testing.T) {
	const limit = 2000

	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{MaxResponseBytes: limit}))
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	result := callTool(t, session, "getCourseLesson", map[string]any{"courseSlug": jsCourse, "slug": firstLesson})

	assert.True(t, result.IsError)
	assert.Contains(t, textOf(t, result), "2000 byte response limit")
	assert.Nil(t, result.StructuredContent)
}

func TestCallThatOverrunsTheTimeoutIsCancelledAndToldSo(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{ToolTimeout: 1}))
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	result := callTool(t, session, "listCourses", nil)

	assert.True(t, result.IsError)
	assert.Contains(t, textOf(t, result), "took longer than")
	assert.Nil(t, result.StructuredContent)
}
