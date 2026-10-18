package relatedcourse_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpostrelatedcourseitem"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/relatedcourse"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport"
)

// Faults a test injects; each is told apart from the others by identity.
var (
	errModelUnavailable = errors.New("model unavailable")
)

func helloWorldPost() int {
	return fixture.ID(fixture.BlogPostHelloWorld)
}

// answeringModel is the LLM boundary answering with a fixed text.
func answeringModel(answer string) *testsupport.JSONCompleterMock {
	return &testsupport.JSONCompleterMock{
		CompleteJSONFunc: func(context.Context, string, string, string, any) (string, error) { return answer, nil },
	}
}

func newSuggester(db *ent.Client, transactor store.Transactor, llm *testsupport.JSONCompleterMock) *relatedcourse.Suggester {
	return relatedcourse.NewSuggester(db, relatedcourse.NewReplacer(transactor), llm,
		slog.New(slog.DiscardHandler))
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

// blog_post_hello_world is ru. The answer names ruby, an invented id, the
// pre-course (a real course without a landing page, so never offered), python,
// and ruby again: only the offered ids survive, once each, in the LLM's order.
func TestSuggestRelatedCoursesReplacesTheSetInOrder(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	ruby, preCourse, python := fixture.ID(fixture.CourseRuby), fixture.ID(fixture.CoursePreCourseJavascript), fixture.ID(fixture.CoursePython)
	llm := answeringModel(fmt.Sprintf(`{"course_ids": [%d, 999999, %d, %d, %d]}`, ruby, preCourse, python, ruby))

	require.NoError(t, newSuggester(db, transactor, llm).SuggestRelatedCourses(t.Context(), helloWorldPost()))

	assert.Equal(t, []string{"ruby", "python"}, relatedSlugs(t, db, helloWorldPost()))
	assert.Equal(t, 2, db.BlogPost.GetX(t.Context(), helloWorldPost()).RelatedCourseItemsCount)
}

// The model is offered the published main landing pages of the post's locale
// and nothing else: an en page, a draft, an archived page or a secondary page
// never reaches it.
func TestOfferedCandidatesAreThePublishedMainPagesOfTheLocale(t *testing.T) {
	db := testsupport.NewClient(t)

	offered, err := relatedcourse.OfferedCandidates(t.Context(), db, "ru")
	require.NoError(t, err)

	names := lo.Map(offered, func(c relatedcourse.Candidate, _ int) string {
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
	llm := answeringModel(`{"course_ids": []}`)

	require.NoError(t, newSuggester(db, transactor, llm).SuggestRelatedCourses(t.Context(), helloWorldPost()))

	assert.Equal(t, []string{"javascript", "ruby", "python"}, relatedSlugs(t, db, helloWorldPost()))
	assert.Equal(t, 3, db.BlogPost.GetX(t.Context(), helloWorldPost()).RelatedCourseItemsCount)
}

// An answer that is not the structured object fails the job, so River retries it, and
// nothing is written.
func TestSuggestRelatedCoursesFailsOnAnUnreadableAnswer(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	llm := answeringModel("Here are the courses: ruby, python")

	require.Error(t, newSuggester(db, transactor, llm).SuggestRelatedCourses(t.Context(), helloWorldPost()))

	assert.Equal(t, []string{"javascript", "ruby", "python"}, relatedSlugs(t, db, helloWorldPost()))
}

// A model that cannot answer fails the job too, so River retries it later, and
// the post keeps the courses it has.
func TestSuggestRelatedCoursesFailsWhenTheModelCannotAnswer(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	llm := &testsupport.JSONCompleterMock{
		CompleteJSONFunc: func(context.Context, string, string, string, any) (string, error) {
			return "", errModelUnavailable
		},
	}

	err := newSuggester(db, transactor, llm).SuggestRelatedCourses(t.Context(), helloWorldPost())

	require.ErrorContains(t, err, "model unavailable")
	assert.Equal(t, []string{"javascript", "ruby", "python"}, relatedSlugs(t, db, helloWorldPost()))
}

// A post deleted after the job was enqueued is an error the job reports, not a
// silent success, and nothing is written for it.
func TestSuggestRelatedCoursesFailsForAPostThatIsGone(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)

	answer := fmt.Sprintf(`{"course_ids": [%d]}`, fixture.ID(fixture.CourseRuby))
	err := newSuggester(db, transactor, answeringModel(answer)).SuggestRelatedCourses(t.Context(), 999999)

	require.Error(t, err)
	assert.Empty(t, relatedSlugs(t, db, 999999))
}
