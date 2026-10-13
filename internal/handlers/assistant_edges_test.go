package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/aichat"
	"hexletbasics/ent/aimessage"
	"hexletbasics/ent/course"
	"hexletbasics/internal/api"
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/testsupport"
)

// A chat read for a lesson that does not exist is a not-found, not an empty chat.
func TestListAssistantMessagesIsNotFoundForAnUnknownLesson(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.ListAssistantMessages(t.Context(), api.ListAssistantMessagesParams{LessonId: 999999})

	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

// A learner's question is answered from the lesson's text in their own language;
// a lesson not written in it is nothing to ask about.
func TestAskAssistantRefusesALessonNotWrittenInTheLocale(t *testing.T) {
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, secondLessonSlug)
	testsupport.SpeakTo(h, "es")
	messages := h.DB.AiMessage.Query().CountX(t.Context())

	res, err := h.Client.CreateAssistantMessage(t.Context(),
		&api.AssistantMessageInput{Message: "¿Ayuda?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)

	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
	assert.Equal(t, messages, h.DB.AiMessage.Query().CountX(t.Context()), "nothing is stored")
}

// A course with nothing built has no lesson text to tutor from.
func TestAskAssistantRefusesACourseWithNothingBuilt(t *testing.T) {
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, secondLessonSlug)
	h.DB.Course.Update().Where(course.SlugEQ(jsCourseSlug)).ClearCurrentVersionID().ExecX(t.Context())

	res, err := h.Client.CreateAssistantMessage(t.Context(),
		&api.AssistantMessageInput{Message: "Help?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)

	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

// When storing fails the learner is told it was a server error, and a question
// that could not be stored is not answered as if it had been.
func TestAskAssistantAnswersAStorageFailureAsAServerError(t *testing.T) {
	noAnswer := func(context.Context, []assistant.Turn, func(string) error) (assistant.Usage, error) {
		return assistant.Usage{}, nil
	}
	cases := []struct {
		name    string
		lesson  string
		model   func(context.Context, []assistant.Turn, func(string) error) (assistant.Usage, error)
		storage func(*testing.T, *testsupport.Harness)
	}{
		{"today's questions cannot be counted", firstLessonSlug, nil, reads(aimessage.Table)},
		{"the earlier conversation cannot be read", firstLessonSlug, nil, column(aimessage.Table, "content")},
		{"the chat cannot be created", secondLessonSlug, nil, writes(aichat.Table, "INSERT")},
		{"the question cannot be stored", firstLessonSlug, nil, writes(aimessage.Table, "INSERT")},
		{"an empty answer's question cannot be stored", firstLessonSlug, noAnswer, writes(aimessage.Table, "INSERT")},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			lesson := lessonBySlug(t, h, c.lesson)
			if c.model != nil {
				h.Assistant.StreamFunc = c.model
			}
			c.storage(t, h)

			_, err := h.Client.CreateAssistantMessage(t.Context(),
				&api.AssistantMessageInput{Message: "Help?"},
				api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
			)

			require.Error(t, err)
			assert.Equal(t, http.StatusInternalServerError, h.LastStatus())
		})
	}
}
