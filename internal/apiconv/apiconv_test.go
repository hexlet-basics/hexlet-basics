package apiconv

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
)

// The hand-written bridges between nullable legacy columns and the contract's
// types. The generated converters call them for every row, so what a missing or
// half-loaded value becomes is the contract the admin screens rely on.

func TestEnrollmentStateOfARowWithoutOneIsStarted(t *testing.T) {
	assert.Equal(t, api.EnrollmentStateStarted, EnrollmentStateFromPtr(nil))
	assert.Equal(t, api.EnrollmentStateStarted, EnrollmentStateFromPtr(lo.ToPtr("")))
	assert.Equal(t, api.EnrollmentState("finished"), EnrollmentStateFromPtr(lo.ToPtr("finished")))
}

func TestNullableColumnsBecomeNullOrTheirValue(t *testing.T) {
	assert.True(t, NilReadinessFromPtr(nil).Null)
	readiness := NilReadinessFromPtr(lo.ToPtr("completed"))
	assert.False(t, readiness.Null)
	assert.Equal(t, api.CourseReadiness("completed"), readiness.Value)

	assert.True(t, NilReviewStateFromPtr(nil).Null)
	state := NilReviewStateFromPtr(lo.ToPtr("published"))
	assert.False(t, state.Null)
	assert.Equal(t, api.ReviewState("published"), state.Value)
}

func TestModuleLessonSlugsKeepOrderAndSkipLessonsWithoutAUrl(t *testing.T) {
	version := func(slug *string) *ent.CourseLessonVersion {
		return &ent.CourseLessonVersion{Edges: ent.CourseLessonVersionEdges{
			Lesson: &ent.CourseLesson{Slug: slug},
		}}
	}
	module := &ent.CourseModuleTranslation{Edges: ent.CourseModuleTranslationEdges{
		Version: &ent.CourseModuleVersion{Edges: ent.CourseModuleVersionEdges{
			LessonVersions: []*ent.CourseLessonVersion{
				version(lo.ToPtr("hello-world")),
				version(nil),
				{}, // a version whose lesson was not loaded
				version(lo.ToPtr("variables")),
			},
		}},
	}}

	assert.Equal(t, []string{"hello-world", "variables"}, moduleLessonSlugs(module))
}

func TestModuleLessonSlugsOfAModuleWithoutAVersionAreEmptyNotNull(t *testing.T) {
	assert.Equal(t, []string{}, moduleLessonSlugs(nil))
	assert.Equal(t, []string{}, moduleLessonSlugs(&ent.CourseModuleTranslation{}))
}

func lessonProgressNamed(names ...string) *ent.LessonProgress {
	infos := make([]*ent.CourseLessonTranslation, len(names))
	for i, name := range names {
		infos[i] = &ent.CourseLessonTranslation{Name: lo.ToPtr(name)}
	}
	return &ent.LessonProgress{Edges: ent.LessonProgressEdges{
		Lesson: &ent.CourseLesson{Edges: ent.CourseLessonEdges{Infos: infos}},
	}}
}

func TestLessonProgressNameIsTheFirstLoadedInfo(t *testing.T) {
	assert.Equal(t, "Hello, World!", lessonProgressLessonName(lessonProgressNamed("Hello, World!", "Привет, мир")))
}

func TestLessonProgressNameIsEmptyWhenNothingNamesTheLesson(t *testing.T) {
	assert.Empty(t, lessonProgressLessonName(nil))
	assert.Empty(t, lessonProgressLessonName(&ent.LessonProgress{}), "lesson not loaded")
	assert.Empty(t, lessonProgressLessonName(lessonProgressNamed()), "lesson with no info in the locale")
}

func TestAssistantMessageNamesTheLessonOfItsChat(t *testing.T) {
	message := &ent.AiMessage{Edges: ent.AiMessageEdges{
		Chat: &ent.AiChat{Edges: ent.AiChatEdges{LessonProgress: lessonProgressNamed("Variables")}},
	}}

	assert.Equal(t, "Variables", assistantMessageLessonName(message))
}

func TestAssistantMessageWithoutALoadedChatHasNoLessonName(t *testing.T) {
	assert.Empty(t, assistantMessageLessonName(nil))
	assert.Empty(t, assistantMessageLessonName(&ent.AiMessage{}), "chat not loaded")
	assert.Empty(t, assistantMessageLessonName(&ent.AiMessage{Edges: ent.AiMessageEdges{Chat: &ent.AiChat{}}}),
		"chat without its lesson progress")
}

func TestCourseDurationIsWholeHoursOfFifteenMinuteLessons(t *testing.T) {
	assert.Equal(t, int32(0), courseDuration(nil), "a course missing from the query reads as no duration")
	assert.Equal(t, int32(0), courseDuration(&ent.Course{LessonsCount: 3}))
	assert.Equal(t, int32(2), courseDuration(&ent.Course{LessonsCount: 9}))
}
