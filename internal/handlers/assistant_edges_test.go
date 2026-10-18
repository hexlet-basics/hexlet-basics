package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"hexletbasics/ent/course"
	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// A chat read for a lesson that does not exist is a not-found, not an empty chat.
func TestListAssistantMessagesIsNotFoundForAnUnknownLesson(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.ListAssistantMessages(t.Context(), api.ListAssistantMessagesParams{LessonId: 999999})

	requireProblem(t, err, http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

// A learner's question is answered from the lesson's text in their own language;
// a lesson not written in it is nothing to ask about.
func TestAskAssistantRefusesALessonNotWrittenInTheLocale(t *testing.T) {
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, secondLessonSlug)
	testsupport.SpeakTo(h, "es")
	messages := h.DB.AiMessage.Query().CountX(t.Context())

	_, err := h.Client.CreateAssistantMessage(t.Context(),
		&api.AssistantMessageInput{Message: "¿Ayuda?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)

	requireProblem(t, err, http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
	assert.Equal(t, messages, h.DB.AiMessage.Query().CountX(t.Context()), "nothing is stored")
}

// A course with nothing built has no lesson text to tutor from.
func TestAskAssistantRefusesACourseWithNothingBuilt(t *testing.T) {
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, secondLessonSlug)
	h.DB.Course.Update().Where(course.SlugEQ(jsCourseSlug)).ClearCurrentVersionID().ExecX(t.Context())

	_, err := h.Client.CreateAssistantMessage(t.Context(),
		&api.AssistantMessageInput{Message: "Help?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)

	requireProblem(t, err, http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}
