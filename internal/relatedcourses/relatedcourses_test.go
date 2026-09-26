package relatedcourses_test

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpostrelatedcourseitem"
	"hexletbasics/internal/relatedcourses"
	"hexletbasics/internal/testsupport"
)

// fakeCompleter records the single call and returns a canned answer.
type fakeCompleter struct {
	answer       string
	calls        int
	instructions string
	prompt       string
}

func (f *fakeCompleter) Complete(_ context.Context, instructions, prompt string) (string, error) {
	f.calls++
	f.instructions = instructions
	f.prompt = prompt
	return f.answer, nil
}

func newSuggester(db *ent.Client, llm *fakeCompleter) *relatedcourses.Suggester {
	return relatedcourses.NewSuggester(db, llm, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

// relatedSlugs is the post's stored related courses in display order, by slug.
func relatedSlugs(t *testing.T, db *ent.Client, postID int) []string {
	t.Helper()
	items := db.BlogPostRelatedCourseItem.Query().
		Where(blogpostrelatedcourseitem.BlogPostID(postID)).
		Order(ent.Asc(blogpostrelatedcourseitem.FieldOrder), ent.Asc(blogpostrelatedcourseitem.FieldID)).
		WithCourse().
		AllX(t.Context())
	return lo.Map(items, func(i *ent.BlogPostRelatedCourseItem, _ int) string {
		return lo.FromPtr(i.Edges.Course.Slug)
	})
}

// Post 6001 is ru. The fenced answer names ruby, an invented id, the
// pre-course (a real course without a landing page, so never offered), python,
// and ruby again: only the offered ids survive, once each, in the LLM's order.
func TestSuggestRelatedCoursesReplacesTheSetInOrder(t *testing.T) {
	db := testsupport.NewClient(t)
	llm := &fakeCompleter{answer: "```json\n[207281424, 999999, 596063838, 617920698, 207281424]\n```"}

	require.NoError(t, newSuggester(db, llm).SuggestRelatedCourses(t.Context(), 6001))

	assert.Equal(t, 1, llm.calls)
	assert.Equal(t, []string{"ruby", "python"}, relatedSlugs(t, db, 6001))
	assert.Equal(t, 2, db.BlogPost.GetX(t.Context(), 6001).RelatedCourseItemsCount)
}

// The prompt is the legacy one: the post's plain text, then the published main
// landing pages of the post's locale as {id, name} JSON.
func TestSuggestRelatedCoursesBuildsTheLegacyPrompt(t *testing.T) {
	db := testsupport.NewClient(t)
	llm := &fakeCompleter{answer: "[]"}

	require.NoError(t, newSuggester(db, llm).SuggestRelatedCourses(t.Context(), 6001))

	assert.Contains(t, llm.instructions, "Ты — ассистент, который помогает подобрать курсы.")
	assert.True(t, strings.HasPrefix(llm.prompt, "Текст статьи: Hello world from the blog"), llm.prompt)
	assert.Contains(t, llm.prompt, "\n\nСписок курсов: [")
	assert.Contains(t, llm.prompt, `{"id":207281424,"name":"Курс Ruby"}`)
	assert.NotContains(t, llm.prompt, "Ruby Course", "an en landing page is not offered for a ru post")
	assert.NotContains(t, llm.prompt, "Архивный курс JavaScript", "an archived landing page is not offered")
}

// Legacy `truncate(2000)`: 1997 characters plus "...", counted in characters,
// not bytes — the posts are Russian.
func TestSuggestRelatedCoursesTruncatesTheText(t *testing.T) {
	db := testsupport.NewClient(t)
	db.BlogPost.UpdateOneID(6001).SetRichBody("<p>" + strings.Repeat("я", 3000) + "</p>").ExecX(t.Context())
	llm := &fakeCompleter{answer: "[]"}

	require.NoError(t, newSuggester(db, llm).SuggestRelatedCourses(t.Context(), 6001))

	assert.Contains(t, llm.prompt, "Текст статьи: "+strings.Repeat("я", 1997)+"...\n\n")
}

// An empty answer leaves the fixture set (javascript, ruby, python) untouched.
func TestSuggestRelatedCoursesKeepsTheSetOnAnEmptyAnswer(t *testing.T) {
	db := testsupport.NewClient(t)
	llm := &fakeCompleter{answer: "```\n[]\n```"}

	require.NoError(t, newSuggester(db, llm).SuggestRelatedCourses(t.Context(), 6001))

	assert.Equal(t, []string{"javascript", "ruby", "python"}, relatedSlugs(t, db, 6001))
	assert.Equal(t, 3, db.BlogPost.GetX(t.Context(), 6001).RelatedCourseItemsCount)
}

// An answer that is not a JSON id list fails the job, so River retries it, and
// nothing is written.
func TestSuggestRelatedCoursesFailsOnAnUnreadableAnswer(t *testing.T) {
	db := testsupport.NewClient(t)
	llm := &fakeCompleter{answer: "Here are the courses: ruby, python"}

	require.Error(t, newSuggester(db, llm).SuggestRelatedCourses(t.Context(), 6001))

	assert.Equal(t, []string{"javascript", "ruby", "python"}, relatedSlugs(t, db, 6001))
}
