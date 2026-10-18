package lessonreview_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/lessonreview"
	"hexletbasics/internal/testsupport"
)

func newScheduler(t *testing.T) (*lessonreview.Scheduler, *testsupport.Services, func(int) int) {
	t.Helper()
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	services := testsupport.NewServices(t, tx, transactor)
	translationsOfLesson := func(lessonID int) int {
		return db.CourseLessonTranslation.Query().Where(courselessontranslation.CourseLessonID(lessonID)).CountX(t.Context())
	}

	return lessonreview.NewScheduler(db, services.Jobs), services, translationsOfLesson
}

func TestReviewLessonSchedulesOneJobPerTranslation(t *testing.T) {
	scheduler, services, translationsOfLesson := newScheduler(t)
	lesson := fixture.ID(fixture.LessonJavascriptHelloWorld)

	require.NoError(t, scheduler.ReviewLesson(t.Context(), lesson))

	queued := testsupport.EnqueuedJobs[jobs.ReviewLessonArgs](t, services)
	assert.Len(t, queued, translationsOfLesson(lesson))
}

func TestReviewOfAMissingLessonOrCourseIsNotFoundAndSchedulesNothing(t *testing.T) {
	scheduler, services, _ := newScheduler(t)

	require.ErrorIs(t, scheduler.ReviewLesson(t.Context(), 0), domain.ErrNotFound)
	require.ErrorIs(t, scheduler.ReviewCourse(t.Context(), 0), domain.ErrNotFound)

	assert.Empty(t, testsupport.EnqueuedJobs[jobs.ReviewLessonArgs](t, services))
}

func TestReviewCourseSchedulesTheCurrentVersionsTranslations(t *testing.T) {
	scheduler, services, _ := newScheduler(t)

	require.NoError(t, scheduler.ReviewCourse(t.Context(), fixture.ID(fixture.CourseJavascript)))

	assert.NotEmpty(t, testsupport.EnqueuedJobs[jobs.ReviewLessonArgs](t, services))
}
