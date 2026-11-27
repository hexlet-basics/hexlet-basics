package apiconv

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/transition"
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

func TestTheStateVocabulariesAreTheContractsEnums(t *testing.T) {
	t.Parallel()

	for _, state := range transition.PublicationStates.States() {
		require.NoError(t, PublicationStateOf(transition.PublicationState(state)).Validate(), state)
	}

	for _, readiness := range transition.CourseReadiness.States() {
		require.NoError(t, ReadinessOf(transition.Readiness(readiness)).Validate(), readiness)
	}
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
// An id or count past the contract's int32 must not wrap into a negative or
// unrelated number on the wire: it saturates at the bound instead.
func TestInt32FromIntSaturatesInsteadOfWrapping(t *testing.T) {
	t.Parallel()
	assert.Equal(t, int32(42), Int32FromInt(42))
	assert.Equal(t, int32(math.MaxInt32), Int32FromInt(math.MaxInt32+1))
	assert.Equal(t, int32(math.MinInt32), Int32FromInt(math.MinInt32-1))
	assert.Equal(t, int32(7), Int32FromPtr(new(7)))
	assert.Equal(t, int32(math.MaxInt32), *nilInt32From(new(math.MaxInt32 + 1)))
}
