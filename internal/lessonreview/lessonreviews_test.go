package lessonreview_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/courselessonreview"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/lessonreview"
	"hexletbasics/internal/testsupport"
)

// reviewingModel is the LLM boundary answering with a fixed summary.
func reviewingModel() *testsupport.CompleterMock {
	return &testsupport.CompleterMock{
		CompleteFunc: func(context.Context, string, string) (string, error) { return "AI summary", nil },
	}
}

// echoingModel answers with the instructions and prompt it was given, so the
// stored summary shows what the reviewer asked the model.
func echoingModel() *testsupport.CompleterMock {
	return &testsupport.CompleterMock{
		CompleteFunc: func(_ context.Context, instructions, prompt string) (string, error) {
			return instructions + "\n" + prompt, nil
		},
	}
}

// The en info of lesson_javascript_hello_world. Its lesson's assistant chat holds
// one student question, so the LLM runs and the existing (course, lesson, en)
// review row is updated in place — not duplicated.
func TestReviewLessonUpdatesExistingReview(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := context.Background()

	reviewer := lessonreview.NewReviewer(db, echoingModel())

	require.NoError(t, reviewer.ReviewLesson(ctx, fixture.ID(fixture.LessonVersionInfoJavascriptHelloWorldCurrentEn)))

	review := db.CourseLessonReview.Query().
		Where(
			courselessonreview.CourseID(fixture.ID(fixture.CourseJavascript)),
			courselessonreview.CourseLessonID(fixture.ID(fixture.LessonJavascriptHelloWorld)),
			courselessonreview.LocaleEQ("en"),
		).
		OnlyX(ctx)
	assert.Contains(t, review.Summary, "Проанализируй вопросы")
	assert.Contains(t, review.Summary, "How do I print a string?", "the student question feeds the prompt")
	assert.Equal(t, fixture.ID(fixture.LessonVersionInfoJavascriptHelloWorldCurrentEn), review.CourseLessonTranslationID)
}

// The en info of lesson_javascript_strings. The lesson has no assistant chats, so
// the summary empties WITHOUT an LLM call (legacy behavior: an empty question
// pool still writes, marking the lesson reviewed; the admin list hides it).
func TestReviewLessonWithoutQuestionsSkipsLLM(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := context.Background()

	reviewer := lessonreview.NewReviewer(db, reviewingModel())

	require.NoError(t, reviewer.ReviewLesson(ctx, fixture.ID(fixture.LessonVersionInfoJavascriptStringsCurrentEn)))

	review := db.CourseLessonReview.Query().
		Where(
			courselessonreview.CourseID(fixture.ID(fixture.CourseJavascript)),
			courselessonreview.CourseLessonID(fixture.ID(fixture.LessonJavascriptStrings)),
			courselessonreview.LocaleEQ("en"),
		).
		OnlyX(ctx)
	assert.Empty(t, review.Summary, "no questions: the model is not asked, so its answer is not stored")
}

// The en info of lesson_javascript_variables. Deleting the fixture review row
// first exercises the CREATE branch of the upsert: the reviewer recreates the
// (course, lesson, locale) row from scratch, pointing it at the reviewed
// version/info.
func TestReviewLessonCreatesMissingReview(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := context.Background()

	_, err := db.CourseLessonReview.Delete().
		Where(
			courselessonreview.CourseID(fixture.ID(fixture.CourseJavascript)),
			courselessonreview.CourseLessonID(fixture.ID(fixture.LessonJavascriptVariables)),
			courselessonreview.LocaleEQ("en"),
		).
		Exec(ctx)
	require.NoError(t, err)

	reviewer := lessonreview.NewReviewer(db, echoingModel())

	require.NoError(t, reviewer.ReviewLesson(ctx, fixture.ID(fixture.LessonVersionInfoJavascriptVariablesCurrentEn)))

	review := db.CourseLessonReview.Query().
		Where(
			courselessonreview.CourseID(fixture.ID(fixture.CourseJavascript)),
			courselessonreview.CourseLessonID(fixture.ID(fixture.LessonJavascriptVariables)),
			courselessonreview.LocaleEQ("en"),
		).
		OnlyX(ctx)
	assert.Contains(t, review.Summary, "What is a variable?")
	assert.Equal(t, fixture.ID(fixture.LessonVersionJavascriptVariablesCurrent), review.CourseLessonVersionID, "the version of the lesson reviewed")
}
