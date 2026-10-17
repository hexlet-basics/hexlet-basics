package mcpserver_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// The fixtures: JavaScript has the lessons hello-world, variables and strings.
// Alice (1001) is enrolled and has finished hello-world; Bob (1002) has a
// finished Enrollment and has only started variables; Carol (1003) has nothing.
// carol is a fixture user.
const carol = 1003

func learnerSession(t *testing.T, o *testsupport.OAuth, userID int, scopes ...string) *mcp.ClientSession {
	t.Helper()
	if len(scopes) == 0 {
		scopes = []string{oauthserver.ScopeProgressRead}
	}
	// Russian, because the fixtures give the JavaScript course its main
	// catalogue card in that locale, and the dashboard lists only courses with one.
	return o.MCPClient(o.Connect(userID, scopes...).AccessToken, russian)
}

// isLearnerTool tells the learner tools from the catalogue's by their prefix.
func isLearnerTool(name string) bool { return strings.HasPrefix(name, "getMy") }

func lessonStates(t *testing.T, course map[string]any) map[string]string {
	t.Helper()
	states := map[string]string{}
	for _, l := range course["lessons"].([]any) {
		lesson := l.(map[string]any)
		states[lesson["slug"].(string)] = lesson["state"].(string)
	}
	return states
}

func TestLearnerToolsAreReadOnlyAndTakeNoParameterNamingAUser(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tools, err := learnerSession(t, o, alice).ListTools(t.Context(), nil)
	require.NoError(t, err)

	names := map[string][]string{}
	for _, tool := range tools.Tools {
		if !isLearnerTool(tool.Name) {
			continue
		}
		assert.True(t, tool.Annotations.ReadOnlyHint, tool.Name)
		props, _ := tool.InputSchema.(map[string]any)["properties"].(map[string]any)
		names[tool.Name] = nil
		for name := range props {
			names[tool.Name] = append(names[tool.Name], name)
		}
	}

	assert.Contains(t, names, "getMyDashboard")
	assert.Empty(t, names["getMyDashboard"])
	assert.ElementsMatch(t, []string{"courseSlug"}, names["getMyCourseProgress"])
	assert.ElementsMatch(t, []string{"courseSlug", "lessonSlug"}, names["getMyLessonProgress"])
}

func TestDashboardToolAnswersWithTheTokensUsersEnrollments(t *testing.T) {
	o := testsupport.NewOAuth(t)

	aliceBoard := structured(t, callTool(t, learnerSession(t, o, alice), "getMyDashboard", nil))
	bobBoard := structured(t, callTool(t, learnerSession(t, o, bob), "getMyDashboard", nil))
	carolBoard := structured(t, callTool(t, learnerSession(t, o, carol), "getMyDashboard", nil))

	assert.Len(t, aliceBoard["startedEnrollments"], 1)
	assert.Empty(t, aliceBoard["finishedEnrollments"])
	assert.Equal(t, float64(33), aliceBoard["startedEnrollments"].([]any)[0].(map[string]any)["completion"])
	assert.Empty(t, bobBoard["startedEnrollments"])
	assert.Len(t, bobBoard["finishedEnrollments"], 1)
	assert.Empty(t, carolBoard["startedEnrollments"])
	assert.Empty(t, carolBoard["finishedEnrollments"])
}

func TestCourseProgressToolTellsFinishedStartedAndNotStarted(t *testing.T) {
	o := testsupport.NewOAuth(t)
	args := map[string]any{"courseSlug": jsCourse}

	a := structured(t, callTool(t, learnerSession(t, o, alice), "getMyCourseProgress", args))
	b := structured(t, callTool(t, learnerSession(t, o, bob), "getMyCourseProgress", args))
	c := structured(t, callTool(t, learnerSession(t, o, carol), "getMyCourseProgress", args))

	assert.Equal(t, "started", a["state"])
	assert.Equal(t, float64(33), a["completion"])
	assert.Equal(t, "variables", a["nextLesson"])
	assert.Equal(t, map[string]string{"hello-world": "finished", "variables": "not_started", "strings": "not_started"}, lessonStates(t, a))

	// Bob's rows are his: he started variables, and alice's finished lesson is
	// not his.
	assert.Equal(t, map[string]string{"hello-world": "not_started", "variables": "started", "strings": "not_started"}, lessonStates(t, b))
	assert.Equal(t, "hello-world", b["nextLesson"])

	assert.Equal(t, "not_started", c["state"])
	assert.Equal(t, float64(0), c["completion"])
}

func TestLessonProgressToolNeverShowsAnotherUsersLesson(t *testing.T) {
	o := testsupport.NewOAuth(t)

	// Bob alone has progressed in variables; Alice alone has finished hello-world.
	// A new Connect replaces the previous token of the client, so each read opens
	// its own session.
	lesson := func(userID int, slug string) map[string]any {
		return structured(t, callTool(t, learnerSession(t, o, userID), "getMyLessonProgress", map[string]any{"courseSlug": jsCourse, "lessonSlug": slug}))
	}

	assert.Equal(t, "started", lesson(bob, "variables")["state"])
	assert.Equal(t, false, lesson(bob, "variables")["passedCheck"])
	assert.Equal(t, "not_started", lesson(alice, "variables")["state"], "bob's lesson is not alice's")
	assert.Equal(t, "finished", lesson(alice, "hello-world")["state"])
	assert.Equal(t, true, lesson(alice, "hello-world")["passedCheck"])
	assert.Equal(t, "not_started", lesson(bob, "hello-world")["state"])
}

func TestLearnerToolsRefuseAnArgumentThatNamesAUser(t *testing.T) {
	o := testsupport.NewOAuth(t)
	s := learnerSession(t, o, alice)

	for _, args := range []map[string]any{
		{"courseSlug": jsCourse, "userId": bob},
		{"courseSlug": jsCourse, "user_id": bob},
		{"courseSlug": jsCourse, "email": "bob@example.com"},
	} {
		result := callTool(t, s, "getMyCourseProgress", args)
		assert.True(t, result.IsError, "%v", args)
		assert.Nil(t, result.StructuredContent)
	}
	assert.True(t, callTool(t, s, "getMyDashboard", map[string]any{"userId": bob}).IsError)
	assert.True(t, callTool(t, s, "getMyLessonProgress", map[string]any{"courseSlug": jsCourse, "lessonSlug": "variables", "userId": bob}).IsError)
}

func TestLearnerToolsReportWhatIsMissing(t *testing.T) {
	o := testsupport.NewOAuth(t)
	s := learnerSession(t, o, alice)

	for name, args := range map[string]map[string]any{
		"getMyCourseProgress": {"courseSlug": "no-such-course"},
		"getMyLessonProgress": {"courseSlug": jsCourse, "lessonSlug": "no-such-lesson"},
	} {
		assert.True(t, callTool(t, s, name, args).IsError, name)
	}
	assert.True(t, callTool(t, s, "getMyCourseProgress", nil).IsError, "an empty course slug is refused")
	assert.True(t, callTool(t, s, "getMyLessonProgress", map[string]any{"courseSlug": jsCourse}).IsError)
	assert.True(t, callTool(t, s, "getMyCourseProgress", map[string]any{"courseSlug": 7}).IsError, "a malformed argument")
}

func TestLearnerToolsRefuseATokenWithoutTheProgressReadScope(t *testing.T) {
	o := testsupport.NewOAuth(t)
	s := learnerSession(t, o, alice, oauthserver.ScopeCatalogRead)

	for name, args := range map[string]map[string]any{
		"getMyDashboard":      nil,
		"getMyCourseProgress": {"courseSlug": jsCourse},
		"getMyLessonProgress": {"courseSlug": jsCourse, "lessonSlug": "variables"},
	} {
		result := callTool(t, s, name, args)
		assert.True(t, result.IsError, name)
		assert.Nil(t, result.StructuredContent, name)
	}
}

// forbidden are the names of what a learner's assistant must never be offered:
// the solution code, the assistant chat, the email and the rest of the profile.
var forbidden = []string{"code", "solution", "message", "chat", "email", "password", "firstname", "lastname", "phone"}

func keys(node any, into *[]string) {
	switch n := node.(type) {
	case map[string]any:
		for k, v := range n {
			*into = append(*into, k)
			keys(v, into)
		}
	case []any:
		for _, v := range n {
			keys(v, into)
		}
	}
}

func TestLearnerSchemasAndResultsCarryNoSolutionChatOrProfileData(t *testing.T) {
	o := testsupport.NewOAuth(t)
	s := learnerSession(t, o, alice)

	tools, err := s.ListTools(t.Context(), nil)
	require.NoError(t, err)
	var learnerTools int
	for _, tool := range tools.Tools {
		if !isLearnerTool(tool.Name) {
			continue
		}
		learnerTools++
		var names []string
		keys(tool.InputSchema, &names)
		keys(tool.OutputSchema, &names)
		for _, name := range names {
			for _, bad := range forbidden {
				assert.NotContains(t, strings.ToLower(name), bad, "%s schema property %s", tool.Name, name)
			}
		}
	}
	require.Equal(t, 3, learnerTools)

	for name, args := range map[string]map[string]any{
		"getMyDashboard":      nil,
		"getMyCourseProgress": {"courseSlug": jsCourse},
		"getMyLessonProgress": {"courseSlug": jsCourse, "lessonSlug": "hello-world"},
	} {
		result := callTool(t, s, name, args)
		var names []string
		keys(result.StructuredContent, &names)
		for _, n := range names {
			for _, bad := range forbidden {
				assert.NotContains(t, strings.ToLower(n), bad, "%s result field %s", name, n)
			}
		}
		raw, err := json.Marshal(result)
		require.NoError(t, err)
		assert.NotContains(t, string(raw), "alice@example.com", name)
	}
}

// The learner tools draw on the personal-data budget, not the catalogue's: a
// caller that exhausted one can still read the other.
func TestLearnerToolsAreLimitedAsPersonalDataReads(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{
		User: mcpserver.Rates{Personal: 1, Catalogue: 5},
	}))
	session := o.MCPClient(o.Connect(alice, oauthserver.ScopeProgressRead, oauthserver.ScopeCatalogRead).AccessToken, russian)

	structured(t, callTool(t, session, "getMyDashboard", nil))
	limited := callTool(t, session, "getMyCourseProgress", map[string]any{"courseSlug": jsCourse})
	assert.True(t, limited.IsError)
	assert.Contains(t, textOf(t, limited), "rate limit exceeded")

	structured(t, callTool(t, session, "listCourses", nil))
}
