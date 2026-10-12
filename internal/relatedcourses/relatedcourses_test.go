package relatedcourses_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpostrelatedcourseitem"
	"hexletbasics/internal/relatedcourses"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport"
)

// answeringModel is the LLM boundary answering with a fixed text.
func answeringModel(answer string) *testsupport.CompleterMock {
	return &testsupport.CompleterMock{
		CompleteFunc: func(context.Context, string, string) (string, error) { return answer, nil },
	}
}

func newSuggester(db *ent.Client, transactor store.Transactor, llm *testsupport.CompleterMock) *relatedcourses.Suggester {
	return relatedcourses.NewSuggester(db, relatedcourses.NewReplacer(transactor), llm,
		slog.New(slog.NewTextHandler(io.Discard, nil)))
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
	db, transactor := testsupport.NewClientWithTransactor(t)
	llm := answeringModel("```json\n[207281424, 999999, 596063838, 617920698, 207281424]\n```")

	require.NoError(t, newSuggester(db, transactor, llm).SuggestRelatedCourses(t.Context(), 6001))

	assert.Equal(t, []string{"ruby", "python"}, relatedSlugs(t, db, 6001))
	assert.Equal(t, 2, db.BlogPost.GetX(t.Context(), 6001).RelatedCourseItemsCount)
}

// The model is offered the published main landing pages of the post's locale
// and nothing else: an en page, a draft, an archived page or a secondary page
// never reaches it.
func TestOfferedCandidatesAreThePublishedMainPagesOfTheLocale(t *testing.T) {
	db := testsupport.NewClient(t)

	offered, err := relatedcourses.OfferedCandidates(t.Context(), db, "ru")
	require.NoError(t, err)

	names := lo.Map(offered, func(c relatedcourses.Candidate, _ int) string {
		return lo.FromPtr(c.Name)
	})
	assert.Contains(t, names, "Курс Ruby")
	assert.NotContains(t, names, "Ruby Course", "an en page is not offered for a ru post")
	assert.NotContains(t, names, "Архивный курс JavaScript", "an archived page is not offered")
	assert.NotContains(t, names, "Курс Node.js", "a draft page is not offered")
	assert.NotContains(t, names, "Курс «Фронтенд»", "a page that is not the main one is not offered")
}

// An empty answer leaves the fixture set (javascript, ruby, python) untouched.
func TestSuggestRelatedCoursesKeepsTheSetOnAnEmptyAnswer(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	llm := answeringModel("```\n[]\n```")

	require.NoError(t, newSuggester(db, transactor, llm).SuggestRelatedCourses(t.Context(), 6001))

	assert.Equal(t, []string{"javascript", "ruby", "python"}, relatedSlugs(t, db, 6001))
	assert.Equal(t, 3, db.BlogPost.GetX(t.Context(), 6001).RelatedCourseItemsCount)
}

// An answer that is not a JSON id list fails the job, so River retries it, and
// nothing is written.
func TestSuggestRelatedCoursesFailsOnAnUnreadableAnswer(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	llm := answeringModel("Here are the courses: ruby, python")

	require.Error(t, newSuggester(db, transactor, llm).SuggestRelatedCourses(t.Context(), 6001))

	assert.Equal(t, []string{"javascript", "ruby", "python"}, relatedSlugs(t, db, 6001))
}
