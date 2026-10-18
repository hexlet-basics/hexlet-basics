package apiconv

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
)

// The hand-written bridges between nullable legacy columns and the contract's
// types. The generated converters call them for every row, so what a missing or
// half-loaded value becomes is the contract the admin screens rely on.

func TestEnrollmentStateOfARowWithoutOneIsStarted(t *testing.T) {
	t.Parallel()
	assert.Equal(t, api.EnrollmentStateStarted, EnrollmentStateFromPtr(nil))
	assert.Equal(t, api.EnrollmentStateStarted, EnrollmentStateFromPtr(new("")))
	assert.Equal(t, api.EnrollmentState("finished"), EnrollmentStateFromPtr(new("finished")))
}

func TestNullableColumnsBecomeNullOrTheirValue(t *testing.T) {
	t.Parallel()
	assert.True(t, NilReadinessFromPtr(nil).Null)
	readiness := NilReadinessFromPtr(new("completed"))
	assert.False(t, readiness.Null)
	assert.Equal(t, api.CourseReadiness("completed"), readiness.Value)

	assert.True(t, NilReviewStateFromPtr(nil).Null)
	state := NilReviewStateFromPtr(new("published"))
	assert.False(t, state.Null)
	assert.Equal(t, api.ReviewState("published"), state.Value)
}

func TestModuleLessonSlugsKeepOrderAndSkipLessonsWithoutAUrl(t *testing.T) {
	t.Parallel()
	version := func(slug *string) *ent.CourseLessonVersion {
		return &ent.CourseLessonVersion{Edges: ent.CourseLessonVersionEdges{
			Lesson: &ent.CourseLesson{Slug: slug},
		}}
	}
	module := &ent.CourseModuleTranslation{Edges: ent.CourseModuleTranslationEdges{
		Version: &ent.CourseModuleVersion{Edges: ent.CourseModuleVersionEdges{
			LessonVersions: []*ent.CourseLessonVersion{
				version(new("hello-world")),
				version(nil),
				{}, // a version whose lesson was not loaded
				version(new("variables")),
			},
		}},
	}}

	assert.Equal(t, []string{"hello-world", "variables"}, moduleLessonSlugs(module))
}

func TestModuleLessonSlugsOfAModuleWithoutAVersionAreEmptyNotNull(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{}, moduleLessonSlugs(nil))
	assert.Equal(t, []string{}, moduleLessonSlugs(&ent.CourseModuleTranslation{}))
}

func lessonProgressNamed(names ...string) *ent.LessonProgress {
	infos := make([]*ent.CourseLessonTranslation, len(names))
	for i, name := range names {
		infos[i] = &ent.CourseLessonTranslation{Name: new(name)}
	}
	return &ent.LessonProgress{Edges: ent.LessonProgressEdges{
		Lesson: &ent.CourseLesson{Edges: ent.CourseLessonEdges{Infos: infos}},
	}}
}

func TestLessonProgressNameIsTheFirstLoadedInfo(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Hello, World!", lessonProgressLessonName(lessonProgressNamed("Hello, World!", "Привет, мир")))
}

func TestLessonProgressNameIsEmptyWhenNothingNamesTheLesson(t *testing.T) {
	t.Parallel()
	assert.Empty(t, lessonProgressLessonName(nil))
	assert.Empty(t, lessonProgressLessonName(&ent.LessonProgress{}), "lesson not loaded")
	assert.Empty(t, lessonProgressLessonName(lessonProgressNamed()), "lesson with no info in the locale")
}

func TestAssistantMessageNamesTheLessonOfItsChat(t *testing.T) {
	t.Parallel()
	message := &ent.AiMessage{Edges: ent.AiMessageEdges{
		Chat: &ent.AiChat{Edges: ent.AiChatEdges{LessonProgress: lessonProgressNamed("Variables")}},
	}}

	assert.Equal(t, "Variables", assistantMessageLessonName(message))
}

func TestAssistantMessageWithoutALoadedChatHasNoLessonName(t *testing.T) {
	t.Parallel()
	assert.Empty(t, assistantMessageLessonName(nil))
	assert.Empty(t, assistantMessageLessonName(&ent.AiMessage{}), "chat not loaded")
	assert.Empty(t, assistantMessageLessonName(&ent.AiMessage{Edges: ent.AiMessageEdges{Chat: &ent.AiChat{}}}),
		"chat without its lesson progress")
}

func TestCourseDurationIsWholeHoursOfFifteenMinuteLessons(t *testing.T) {
	t.Parallel()
	assert.Equal(t, int32(0), courseDuration(nil), "a course missing from the query reads as no duration")
	assert.Equal(t, int32(0), courseDuration(&ent.Course{LessonsCount: 3}))
	assert.Equal(t, int32(2), courseDuration(&ent.Course{LessonsCount: 9}))
}

// The lessons on either side of one are decided from the list's order alone.
func TestAdjacentLessonsAreTheNeighboursInTheGivenOrder(t *testing.T) {
	t.Parallel()
	lessons := []api.LearnerLesson{
		{Slug: "a", Name: api.NewNilString("A")},
		{Slug: "b", Name: api.NewNilString("B")},
		{Slug: "c", Name: api.NewNilString("C")},
	}

	middle := AdjacentLessons("b", lessons)
	assert.Equal(t, "a", middle.Previous.Value.Slug)
	assert.Equal(t, "C", middle.Next.Value.Name.Value, "named, not just addressed")

	first := AdjacentLessons("a", lessons)
	assert.True(t, first.Previous.Null)
	assert.Equal(t, "b", first.Next.Value.Slug)

	last := AdjacentLessons("c", lessons)
	assert.Equal(t, "b", last.Previous.Value.Slug)
	assert.True(t, last.Next.Null)
}

// A lesson that is not in the list, or the only one in it, has no neighbours:
// both ends are null rather than a guess.
func TestAdjacentLessonsOfALoneOrUnlistedLessonAreNull(t *testing.T) {
	t.Parallel()
	lone := AdjacentLessons("a", []api.LearnerLesson{{Slug: "a"}})
	assert.True(t, lone.Previous.Null)
	assert.True(t, lone.Next.Null)

	unlisted := AdjacentLessons("z", []api.LearnerLesson{{Slug: "a"}, {Slug: "b"}})
	assert.True(t, unlisted.Previous.Null)
	assert.True(t, unlisted.Next.Null)
}

// The Next Lesson is named, and addressed, from the list; with no slug, or a
// lesson the list does not have (no translation in the locale: no page there),
// there is nothing to point at.
func TestNextLessonRefIsNamedFromTheList(t *testing.T) {
	t.Parallel()
	lessons := []api.LearnerLesson{{ID: 7, Slug: "a", Name: api.NewNilString("A")}}

	assert.True(t, NextLessonRef(nil, lessons).Null, "a finished course has none")

	named := NextLessonRef(new("a"), lessons).Value
	assert.Equal(t, int32(7), named.ID)
	assert.Equal(t, "A", named.Name.Value)

	assert.True(t, NextLessonRef(new("z"), lessons).Null)
}
