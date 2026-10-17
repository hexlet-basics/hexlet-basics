package mcpserver_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// english makes the connection speak English, whose pages have no locale prefix.
func english(r *http.Request) { r.Header.Set("Accept-Language", "en") }

// resultsOf is the results array of a search call.
func resultsOf(t *testing.T, result *mcp.CallToolResult) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, item := range structured(t, result)["results"].([]any) {
		out = append(out, item.(map[string]any))
	}
	return out
}

func idsOfResults(results []map[string]any) []string {
	ids := make([]string, len(results))
	for i, r := range results {
		ids[i] = r["id"].(string)
	}
	return ids
}

func TestSearchAnswersInTheShapeChatGPTConnectorsRequire(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)
	schemas := outputSchemas(t, session)

	result := callTool(t, session, "search", map[string]any{"query": "Ruby"})

	content := structured(t, result)
	require.NotEmpty(t, content["results"])
	for _, item := range resultsOf(t, result) {
		assert.NotEmpty(t, item["id"])
		assert.NotEmpty(t, item["title"])
		assert.True(t, strings.HasPrefix(item["url"].(string), testsupport.SiteURL+"/"), "an absolute URL on the configured site: %v", item["url"])
	}
	assert.NoError(t, schemas["search"].Validate(content))

	// The text block is the same JSON as the structured content.
	var fromText map[string]any
	require.NoError(t, json.Unmarshal([]byte(textOf(t, result)), &fromText))
	assert.Equal(t, content, fromText)
}

func TestSearchAndFetchAreReadOnlyAlwaysLoadedAndNeedTheCatalogueScope(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	found := map[string]*mcp.Tool{}
	for _, tool := range tools.Tools {
		found[tool.Name] = tool
	}
	for _, name := range []string{"search", "fetch"} {
		require.Contains(t, found, name)
		assert.True(t, found[name].Annotations.ReadOnlyHint, name)
		assert.Equal(t, true, found[name].Meta["anthropic/alwaysLoad"], "%s must stay loaded for clients that defer tools", name)
	}

	// Without the Scope neither answers.
	progressOnly := connect(t, o, alice, oauthserver.ScopeProgressRead)
	for name, args := range map[string]map[string]any{"search": {"query": "Ruby"}, "fetch": {"id": "post:learning-ruby"}} {
		denied := callTool(t, progressOnly, name, args)
		assert.True(t, denied.IsError, name)
		assert.Contains(t, textOf(t, denied), oauthserver.ScopeCatalogRead)
	}
}

// What a person sees on the website and what an assistant is told are the same
// results, so one cannot disagree with the other.
func TestSearchToolAnswersWhatTheHTTPOperationAnswersAVisitor(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	result := callTool(t, session, "search", map[string]any{"query": "Привет"})

	res := o.Do(http.MethodGet, "/api/search?q="+url.QueryEscape("Привет"), nil, russian)
	defer func() { _ = res.Body.Close() }()
	require.Equal(t, http.StatusOK, res.StatusCode)
	var visitor struct {
		Items []struct{ ID, Title, Path string }
		Total int
	}
	require.NoError(t, json.NewDecoder(res.Body).Decode(&visitor))
	require.NotEmpty(t, visitor.Items)

	tool := resultsOf(t, result)
	require.Len(t, tool, len(visitor.Items))
	for i, item := range visitor.Items {
		assert.Equal(t, item.ID, tool[i]["id"])
		assert.Equal(t, item.Title, tool[i]["title"])
		assert.Equal(t, testsupport.SiteURL+item.Path, tool[i]["url"])
	}
	assert.EqualValues(t, visitor.Total, structured(t, result)["total"])
}

func TestEveryIdSearchReturnsIsAcceptedByFetch(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	seen := map[string][]string{}
	for _, tc := range []struct{ kind, query string }{
		{"course", "JavaScript"}, {"lesson", "первая программа"}, {"post", "Ruby"},
	} {
		results := resultsOf(t, callTool(t, session, "search", map[string]any{"query": tc.query, "type": tc.kind}))
		require.NotEmpty(t, results, tc.kind)
		for _, hit := range results {
			assert.Equal(t, tc.kind, hit["type"])
			id := hit["id"].(string)
			assert.True(t, strings.HasPrefix(id, tc.kind+":"), "the id is typed: %s", id)
			seen[tc.kind] = append(seen[tc.kind], id)

			fetched := callTool(t, session, "fetch", map[string]any{"id": id})

			doc := structured(t, fetched)
			assert.Equal(t, id, doc["id"])
			assert.Equal(t, hit["title"], doc["title"])
			assert.Equal(t, hit["url"], doc["url"], "a search result and the record it names point at the same page: %s", id)
			assert.NotEmpty(t, doc["text"], id)
			assert.Equal(t, tc.kind, doc["metadata"].(map[string]any)["type"])
		}
	}
	assert.Len(t, seen, 3, "every kind of record was searched and fetched")
}

func TestFetchAnswersInTheShapeChatGPTConnectorsRequire(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)
	schemas := outputSchemas(t, session)

	result := callTool(t, session, "fetch", map[string]any{"id": "lesson:" + jsCourse + "/" + firstLesson})

	doc := structured(t, result)
	assert.NoError(t, schemas["fetch"].Validate(doc))
	assert.Equal(t, "lesson:"+jsCourse+"/"+firstLesson, doc["id"])
	assert.Equal(t, testsupport.SiteURL+"/ru/languages/"+jsCourse+"/lessons/"+firstLesson, doc["url"])
	text := doc["text"].(string)
	assert.Contains(t, text, "console.log", "the lesson's theory")
	var fromText map[string]any
	require.NoError(t, json.Unmarshal([]byte(textOf(t, result)), &fromText))
	assert.Equal(t, doc, fromText, "the text block is the same JSON as the structured content")
}

// The text of a lesson carries what the HTTP read gives every visitor, the code
// among it (#851): fetch adds no rule of its own.
func TestFetchOfALessonCarriesWhatGetCourseLessonGivesAVisitor(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	lesson := structured(t, callTool(t, session, "getCourseLesson", map[string]any{"courseSlug": jsCourse, "slug": firstLesson}))["lesson"].(map[string]any)
	doc := structured(t, callTool(t, session, "fetch", map[string]any{"id": "lesson:" + jsCourse + "/" + firstLesson}))

	text := doc["text"].(string)
	for _, field := range []string{"theory", "instructions", "preparedCode", "originalCode", "testCode"} {
		if value, ok := lesson[field].(string); ok && strings.TrimSpace(value) != "" {
			assert.Contains(t, text, strings.TrimSpace(value), field)
		}
	}
	assert.Equal(t, lesson["name"], doc["title"])
}

func TestFetchOfACourseListsItsLessonsWithTheIdsThatFetchThem(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	doc := structured(t, callTool(t, session, "fetch", map[string]any{"id": "course:" + jsLanding}))

	assert.Contains(t, doc["text"], "lesson:"+jsCourse+"/"+firstLesson)
	assert.Equal(t, testsupport.SiteURL+"/ru/languages/"+jsLanding, doc["url"])
}

func TestFetchOfAPostIsItsTextWithoutMarkup(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	doc := structured(t, callTool(t, session, "fetch", map[string]any{"id": "post:learning-ruby"}))

	text := doc["text"].(string)
	assert.Contains(t, text, "Простой синтаксис")
	assert.NotContains(t, text, "<p>")
	assert.NotContains(t, text, "<h2>")
	assert.Contains(t, text, "\n\n", "the paragraphs are kept apart")
}

func TestTheLinksCarryTheLocaleOfTheConnection(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)
	inEnglish := o.MCPClient(tokens.AccessToken, english)

	hits := resultsOf(t, callTool(t, inEnglish, "search", map[string]any{"query": "Hello readers", "type": "post"}))
	require.Equal(t, []string{"post:english-post"}, idsOfResults(hits))
	assert.Equal(t, testsupport.SiteURL+"/blog_posts/english-post", hits[0]["url"], "English is the unprefixed locale")

	doc := structured(t, callTool(t, inEnglish, "fetch", map[string]any{"id": "post:english-post"}))
	assert.Equal(t, testsupport.SiteURL+"/blog_posts/english-post", doc["url"])
}

func TestAnIdThatNamesNothingAVisitorMayReadIsAClearToolError(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	for name, id := range map[string]string{
		"a post that does not exist":         "post:no-such-post",
		"an archived post":                   "post:archived-notes",
		"a post of another locale":           "post:english-post",
		"a course that does not exist":       "course:no-such-course",
		"a draft landing page":               "course:nodejs-ru",
		"a lesson of a course that is not":   "lesson:no-such-course/hello-world",
		"a lesson that is not in the course": "lesson:" + jsCourse + "/no-such-lesson",
	} {
		t.Run(name, func(t *testing.T) {
			result := callTool(t, session, "fetch", map[string]any{"id": id})

			require.True(t, result.IsError)
			assert.Contains(t, textOf(t, result), id)
			assert.Nil(t, result.StructuredContent)
		})
	}
}

// A landing page the lists do not show is not one search returns, and fetch does
// not hand it out either, though the website would serve its URL.
func TestAnUnlistedCourseIsNotFetched(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)
	require.False(t, callTool(t, session, "fetch", map[string]any{"id": "course:ruby-ru"}).IsError, "listed, it is fetched")

	o.DB.CourseLandingPage.Update().Where(courselandingpage.SlugEQ("ruby-ru")).SetListed(false).ExecX(t.Context())

	assert.True(t, callTool(t, session, "fetch", map[string]any{"id": "course:ruby-ru"}).IsError)
}

func TestAnArchivedCourseSaysWhereItMovedInsteadOfAnswering(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	result := callTool(t, session, "fetch", map[string]any{"id": "course:javascript-ru-archived"})

	// Either the page names where it moved (and fetch says so) or it is not
	// found: it is never answered as if it were the current page.
	require.True(t, result.IsError)
	assert.Contains(t, textOf(t, result), "javascript-ru-archived")
}

func TestWhatIsNotAnIdIsAClearToolErrorThatSaysWhatAnIdLooksLike(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	for _, id := range []string{"javascript-ru", "page:home", "lesson:javascript", "post:"} {
		result := callTool(t, session, "fetch", map[string]any{"id": id})

		require.True(t, result.IsError, id)
		assert.Contains(t, textOf(t, result), "course:SLUG")
	}
}

func TestSearchRefusesWhatCannotBeASearch(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	for name, args := range map[string]map[string]any{
		"an empty query":           {"query": "   "},
		"a query that is too long": {"query": strings.Repeat("a", 201)},
		"an unknown type":          {"query": "Ruby", "type": "book"},
		"an unknown argument":      {"query": "Ruby", "locale": "en"},
		"a negative page":          {"query": "Ruby", "page": -1},
	} {
		result := callTool(t, session, "search", args)
		assert.True(t, result.IsError, name)
	}
}

func TestSearchPagesLikeTheOtherCatalogueToolsAndClampsThePageSize(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{MaxPageSize: 2}))
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	first := callTool(t, session, "search", map[string]any{"query": "Привет", "perPage": 100})
	second := callTool(t, session, "search", map[string]any{"query": "Привет", "page": 2, "perPage": 100})

	assert.Len(t, resultsOf(t, first), 2, "a page above the limit is clamped to it")
	assert.EqualValues(t, 2, structured(t, first)["perPage"])
	assert.Greater(t, structured(t, first)["total"], float64(2))
	assert.NotEqual(t, idsOfResults(resultsOf(t, first)), idsOfResults(resultsOf(t, second)))
}

func TestSearchAndFetchSpendTheCatalogueBudget(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{
		RateWindow: time.Minute, User: mcpserver.Rates{Catalogue: 2},
	}))
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	require.False(t, callTool(t, session, "search", map[string]any{"query": "Ruby"}).IsError)
	require.False(t, callTool(t, session, "fetch", map[string]any{"id": "post:learning-ruby"}).IsError)

	rejected(t, callTool(t, session, "search", map[string]any{"query": "Ruby"}))
}

func TestAFetchedTextOverTheResponseLimitIsCutWithTheUsualHintAndStaysTheSameInBothBlocks(t *testing.T) {
	const limit = 3000
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{MaxResponseBytes: limit}))
	session := connect(t, o, alice, oauthserver.ScopeCatalogRead)

	post := o.DB.BlogPost.Query().Where(blogpost.SlugEQ("learning-ruby")).OnlyX(t.Context())
	require.Greater(t, len(post.RichBody), limit/2, "the fixture's post must be longer than half the limit (the document is in both blocks), or the test proves nothing")

	result := callTool(t, session, "fetch", map[string]any{"id": "post:learning-ruby"})

	doc := structured(t, result)
	assert.Contains(t, doc, "truncation")
	assert.Equal(t, true, doc["truncation"].(map[string]any)["truncated"])
	assert.Contains(t, doc["text"], "[cut]")
	assert.LessOrEqual(t, len(textOf(t, result))+mustMarshal(t, doc), limit, "both blocks together are within the limit")
	var fromText map[string]any
	require.NoError(t, json.Unmarshal([]byte(textOf(t, result)), &fromText))
	assert.Equal(t, doc, fromText, "the text block carries the same cut document")
}

func mustMarshal(t *testing.T, v any) int {
	t.Helper()
	raw, err := json.Marshal(v)
	require.NoError(t, err)
	return len(raw)
}
