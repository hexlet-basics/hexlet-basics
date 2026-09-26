package handlers_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/aichat"
	"hexletbasics/ent/aimessage"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/api"
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/testsupport"
)

// The acting learner's hello-world progress carries a fixture chat with one
// question and its answer; the history comes back oldest first, with the
// lesson identity the panel shows.
func TestListAssistantMessagesReturnsTheLessonChat(t *testing.T) {
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, firstLessonSlug)

	chat := listAssistantChat(t, h, lesson)

	require.Len(t, chat.Messages, 2)
	assert.Equal(t, "user", chat.Messages[0].Role)
	assert.Equal(t, "How do I print a string?", chat.Messages[0].Content)
	assert.Equal(t, "assistant", chat.Messages[1].Role)
	assert.Equal(t, "Use console.log to print a string.", chat.Messages[1].Content)
	assert.Equal(t, jsCourseSlug, chat.Messages[1].CourseSlug)
	assert.Equal(t, firstLessonSlug, chat.Messages[1].CourseLessonSlug)
	assert.Equal(t, "Hello, World!", chat.Messages[1].CourseLessonName)
	assert.False(t, chat.QuotaExceeded)
}

// Provider bookkeeping rows (system prompts, tool calls) live in the same
// table in production; the learner only ever sees the conversation.
func TestListAssistantMessagesHidesProviderRoles(t *testing.T) {
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, firstLessonSlug)
	chat := chatFor(t, h, lesson)
	h.DB.AiMessage.Create().
		SetAiChatID(chat.ID).
		SetRole("system").
		SetContent("You help people learn").
		ExecX(t.Context())

	listed := listAssistantChat(t, h, lesson)

	for _, message := range listed.Messages {
		assert.NotEqual(t, "system", message.Role)
	}
}

// Before the first question there is no chat, and reading the panel must not
// create one — nor start the lesson.
func TestListAssistantMessagesIsEmptyBeforeTheFirstQuestion(t *testing.T) {
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, secondLessonSlug)
	chats := h.DB.AiChat.Query().CountX(t.Context())

	chat := listAssistantChat(t, h, lesson)

	assert.Empty(t, chat.Messages)
	assert.False(t, chat.QuotaExceeded)
	assert.Equal(t, chats, h.DB.AiChat.Query().CountX(t.Context()))
	assert.False(t, hasLessonProgress(t, h, lesson))
}

// Asking streams the model's tokens back as the body and, once the answer is
// complete, stores the question and the answer with its token usage. The
// learner never started the lesson; asking starts it.
func TestAskAssistantStreamsTheAnswerAndStoresBothTurns(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, secondLessonSlug)
	h.Assistant.Deltas = []string{"Declare it ", "with ", "`let`."}
	h.Assistant.Usage = assistant.Usage{InputTokens: 321, OutputTokens: 12}
	counterBefore := assistantMessagesCount(t, h)

	answer := askAssistant(t, h, lesson, "How do I declare a variable?")

	assert.Equal(t, "Declare it with `let`.", answer)
	assert.True(t, hasLessonProgress(t, h, lesson), "asking is doing the lesson")

	chat := chatFor(t, h, lesson)
	assert.Equal(t, actingUserID(t, h), chat.UserID)
	stored := h.DB.AiMessage.Query().Where(aimessage.AiChatID(chat.ID)).Order(ent.Asc(aimessage.FieldID)).AllX(ctx)
	require.Len(t, stored, 2)

	assert.Equal(t, "user", stored[0].Role)
	assert.Equal(t, "How do I declare a variable?", *stored[0].Content, "the question as asked, not the wrapped prompt")
	require.NotNil(t, stored[0].UserID)
	assert.Equal(t, actingUserID(t, h), *stored[0].UserID)

	assert.Equal(t, "assistant", stored[1].Role)
	assert.Equal(t, "Declare it with `let`.", *stored[1].Content)
	assert.Nil(t, stored[1].UserID)
	assert.Equal(t, 321, *stored[1].InputTokens)
	assert.Equal(t, 12, *stored[1].OutputTokens)

	assert.Equal(t, counterBefore+1, assistantMessagesCount(t, h), "the users counter counts questions, as the legacy counter cache did")
}

// The model sees the lesson it tutors (system prompt from the lesson's info in
// the request locale), and the question arrives wrapped with the editor state.
func TestAskAssistantBuildsThePromptFromTheLesson(t *testing.T) {
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, secondLessonSlug)

	askAssistantWith(t, h, lesson, &api.AssistantMessageInput{
		Message:  "Why does it fail?",
		UserCode: api.NewNilString("let x = 1;"),
		Output:   api.NewNilString("expected 2, got 1"),
	})

	require.Len(t, h.Assistant.Requests, 1)
	turns := h.Assistant.Requests[0]
	require.Len(t, turns, 2, "a first question: the system prompt and the question")

	system := turns[0]
	assert.Equal(t, assistant.RoleSystem, system.Role)
	assert.Contains(t, system.Content, jsCourseSlug)
	assert.Contains(t, system.Content, `"Variables"`)
	assert.Contains(t, system.Content, "Answer in English")

	question := turns[1]
	assert.Equal(t, assistant.RoleUser, question.Role)
	assert.Contains(t, question.Content, "let x = 1;")
	assert.Contains(t, question.Content, "expected 2, got 1")
	assert.Contains(t, question.Content, "Why does it fail?")
}

// A learner reading the course in Russian is tutored in Russian, from the
// Russian lesson text.
func TestAskAssistantPromptFollowsTheLocale(t *testing.T) {
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")
	lesson := lessonBySlug(t, h, firstLessonSlug)

	askAssistant(t, h, lesson, "Как вывести строку?")

	require.Len(t, h.Assistant.Requests, 1)
	system := h.Assistant.Requests[0][0].Content
	assert.Contains(t, system, "«Привет, мир!»")
	assert.Contains(t, system, "Отвечай на русском языке")
}

// A follow-up question carries the conversation so far, between the system
// prompt and the new question.
func TestAskAssistantSendsTheChatHistory(t *testing.T) {
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, firstLessonSlug)

	askAssistant(t, h, lesson, "And a number?")

	require.Len(t, h.Assistant.Requests, 1)
	turns := h.Assistant.Requests[0]
	require.Len(t, turns, 4)
	assert.Equal(t, assistant.Turn{Role: assistant.RoleUser, Content: "How do I print a string?"}, turns[1])
	assert.Equal(t, assistant.Turn{Role: assistant.RoleAssistant, Content: "Use console.log to print a string."}, turns[2])
	assert.Contains(t, turns[3].Content, "And a number?")
}

// Seven questions a day. Over the quota the operation answers 429 before the
// model is called, stores nothing, and the history reports the disabled state.
func TestAskAssistantRefusesOverTheDailyQuota(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, firstLessonSlug)
	askedAt(t, h, lesson, assistant.DailyQuota, time.Now())
	messages := h.DB.AiMessage.Query().CountX(ctx)

	res, err := h.Client.CreateAssistantMessage(ctx,
		&api.AssistantMessageInput{Message: "One more?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)
	require.NoError(t, err)
	assert.IsType(t, &api.CreateAssistantMessageTooManyRequests{}, res)
	assert.Equal(t, http.StatusTooManyRequests, h.LastStatus())
	assert.Empty(t, h.Assistant.Requests, "the model is never asked")
	assert.Equal(t, messages, h.DB.AiMessage.Query().CountX(ctx))

	assert.True(t, listAssistantChat(t, h, lesson).QuotaExceeded)
}

// The quota is per day: yesterday's questions do not count against today.
func TestAskAssistantQuotaResetsDaily(t *testing.T) {
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, firstLessonSlug)
	askedAt(t, h, lesson, assistant.DailyQuota, time.Now().UTC().Truncate(24*time.Hour).Add(-time.Minute))

	askAssistant(t, h, lesson, "A new day?")

	assert.Len(t, h.Assistant.Requests, 1)
}

// Asking about a lesson beyond the gate is refused as starting it would be,
// and nothing is stored.
func TestAskAssistantRefusesALessonBeyondTheGate(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, thirdLessonSlug)
	messages := h.DB.AiMessage.Query().CountX(ctx)

	res, err := h.Client.CreateAssistantMessage(ctx,
		&api.AssistantMessageInput{Message: "Help?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)
	require.NoError(t, err)
	assert.IsType(t, &api.CreateAssistantMessageConflict{}, res)
	assert.Equal(t, http.StatusConflict, h.LastStatus())
	assert.Empty(t, h.Assistant.Requests)
	assert.Equal(t, messages, h.DB.AiMessage.Query().CountX(ctx))
}

// When the model fails before its first token the learner gets an error
// status rather than an empty answer, and the question is not stored — it
// does not count against their quota.
func TestAskAssistantReportsAModelFailureBeforeStreaming(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, firstLessonSlug)
	h.Assistant.Err = io.ErrUnexpectedEOF
	messages := h.DB.AiMessage.Query().CountX(ctx)

	_, err := h.Client.CreateAssistantMessage(ctx,
		&api.AssistantMessageInput{Message: "Help?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, h.LastStatus())
	assert.Equal(t, messages, h.DB.AiMessage.Query().CountX(ctx))
}

// The chat is a signed-in feature: a visitor is refused by the contract.
func TestAskAssistantRequiresASignedInLearner(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	lesson := lessonBySlug(t, h, firstLessonSlug)

	res, err := h.Client.CreateAssistantMessage(t.Context(),
		&api.AssistantMessageInput{Message: "Help?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)
	require.NoError(t, err)
	assert.IsType(t, &api.CreateAssistantMessageUnauthorized{}, res)
	assert.Empty(t, h.Assistant.Requests)
}

func listAssistantChat(t *testing.T, h *testsupport.Harness, lesson *ent.CourseLesson) *api.LessonAssistantChat {
	t.Helper()
	res, err := h.Client.ListAssistantMessages(t.Context(), api.ListAssistantMessagesParams{LessonId: int32(lesson.ID)})
	require.NoError(t, err)
	require.IsType(t, &api.LessonAssistantChat{}, res)
	return res.(*api.LessonAssistantChat)
}

func askAssistant(t *testing.T, h *testsupport.Harness, lesson *ent.CourseLesson, question string) string {
	t.Helper()
	return askAssistantWith(t, h, lesson, &api.AssistantMessageInput{Message: question})
}

func askAssistantWith(t *testing.T, h *testsupport.Harness, lesson *ent.CourseLesson, input *api.AssistantMessageInput) string {
	t.Helper()
	res, err := h.Client.CreateAssistantMessage(t.Context(), input,
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)
	require.NoError(t, err)
	require.IsType(t, &api.CreateAssistantMessageOK{}, res)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	body, err := io.ReadAll(res.(*api.CreateAssistantMessageOK))
	require.NoError(t, err)
	return string(body)
}

// chatFor is the acting learner's chat on lesson, found the way the handlers
// find it: through their lesson progress.
func chatFor(t *testing.T, h *testsupport.Harness, lesson *ent.CourseLesson) *ent.AiChat {
	t.Helper()
	return h.DB.AiChat.Query().
		Where(aichat.HasLessonProgressWith(
			lessonprogress.LessonID(lesson.ID),
			lessonprogress.UserID(actingUserID(t, h)),
		)).
		OnlyX(t.Context())
}

func hasLessonProgress(t *testing.T, h *testsupport.Harness, lesson *ent.CourseLesson) bool {
	t.Helper()
	return h.DB.LessonProgress.Query().
		Where(lessonprogress.LessonID(lesson.ID), lessonprogress.UserID(actingUserID(t, h))).
		ExistX(t.Context())
}

// askedAt records n questions by the acting learner at the given time.
func askedAt(t *testing.T, h *testsupport.Harness, lesson *ent.CourseLesson, n int, at time.Time) {
	t.Helper()
	chat := chatFor(t, h, lesson)
	for i := range n {
		h.DB.AiMessage.Create().
			SetAiChatID(chat.ID).
			SetRole("user").
			SetContent(strings.Repeat("?", i+1)).
			SetUserID(actingUserID(t, h)).
			SetCreatedAt(at).
			SetUpdatedAt(at).
			ExecX(t.Context())
	}
}

func assistantMessagesCount(t *testing.T, h *testsupport.Harness) int {
	t.Helper()
	u := h.DB.User.GetX(t.Context(), actingUserID(t, h))
	if u.AssistantMessagesCount == nil {
		return 0
	}
	return *u.AssistantMessagesCount
}
