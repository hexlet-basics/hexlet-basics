package lessonreviews_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/courselessonreview"
	"hexletbasics/internal/lessonreviews"
	"hexletbasics/internal/testsupport"
)

// fakeCompleter records the single call and returns a canned summary.
type fakeCompleter struct {
	calls        int
	instructions string
	prompt       string
}

func (f *fakeCompleter) Complete(_ context.Context, instructions, prompt string) (string, error) {
	f.calls++
	f.instructions = instructions
	f.prompt = prompt
	return "AI summary", nil
}

// Info 2002 = hello-world (lesson 1001) en. Its lesson's assistant chat holds
// one student question, so the LLM runs and the existing (course, lesson, en)
// review row is updated in place — not duplicated.
func TestReviewLessonUpdatesExistingReview(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := context.Background()

	llm := &fakeCompleter{}
	reviewer := lessonreviews.NewReviewer(db, llm)

	require.NoError(t, reviewer.ReviewLesson(ctx, 2002))

	assert.Equal(t, 1, llm.calls)
	assert.Contains(t, llm.instructions, "Проанализируй вопросы")
	assert.Contains(t, llm.prompt, "How do I print a string?", "the student question feeds the prompt")

	review := db.CourseLessonReview.Query().
		Where(
			courselessonreview.CourseID(82481401),
			courselessonreview.CourseLessonID(1001),
			courselessonreview.LocaleEQ("en"),
		).
		OnlyX(ctx)
	assert.Equal(t, "AI summary", review.Summary)
	assert.Equal(t, 2002, review.CourseLessonTranslationID)
}

// Info 2001 = strings (lesson 1003) en. The lesson has no assistant chats, so
// the summary empties WITHOUT an LLM call (legacy behavior: an empty question
// pool still writes, marking the lesson reviewed; the admin list hides it).
func TestReviewLessonWithoutQuestionsSkipsLLM(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := context.Background()

	llm := &fakeCompleter{}
	reviewer := lessonreviews.NewReviewer(db, llm)

	require.NoError(t, reviewer.ReviewLesson(ctx, 2001))

	assert.Zero(t, llm.calls, "no questions -> no LLM call")

	review := db.CourseLessonReview.Query().
		Where(
			courselessonreview.CourseID(82481401),
			courselessonreview.CourseLessonID(1003),
			courselessonreview.LocaleEQ("en"),
		).
		OnlyX(ctx)
	assert.Equal(t, "", review.Summary)
}

// Info 2003 = variables (lesson 1002) en. Deleting the fixture review row
// first exercises the CREATE branch of the upsert: the reviewer recreates the
// (course, lesson, locale) row from scratch, pointing it at the reviewed
// version/info.
func TestReviewLessonCreatesMissingReview(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := context.Background()

	_, err := db.CourseLessonReview.Delete().
		Where(
			courselessonreview.CourseID(82481401),
			courselessonreview.CourseLessonID(1002),
			courselessonreview.LocaleEQ("en"),
		).
		Exec(ctx)
	require.NoError(t, err)

	llm := &fakeCompleter{}
	reviewer := lessonreviews.NewReviewer(db, llm)

	require.NoError(t, reviewer.ReviewLesson(ctx, 2003))

	assert.Equal(t, 1, llm.calls)
	assert.Contains(t, llm.prompt, "What is a variable?")

	review := db.CourseLessonReview.Query().
		Where(
			courselessonreview.CourseID(82481401),
			courselessonreview.CourseLessonID(1002),
			courselessonreview.LocaleEQ("en"),
		).
		OnlyX(ctx)
	assert.Equal(t, "AI summary", review.Summary)
	assert.Equal(t, 965227298, review.CourseLessonVersionID)
}
