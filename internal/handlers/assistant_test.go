package handlers_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/aichat"
	"hexletbasics/ent/aimessage"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/api"
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// The acting learner's hello-world progress carries a fixture chat with one
// question and its answer; the history comes back oldest first, with the
// lesson identity the panel shows.
func TestListAssistantMessagesReturnsTheLessonChat(t *testing.T) {
	t.Parallel()
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
// table in production; the learner only ever sees the conversation. The
// fixture chat on the first lesson holds a system prompt.
func TestListAssistantMessagesHidesProviderRoles(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, firstLessonSlug)
	require.True(t, h.DB.AiMessage.Query().Where(aimessage.RoleEQ("system")).ExistX(t.Context()), "the fixtures hold a system prompt")

	listed := listAssistantChat(t, h, lesson)

	for _, message := range listed.Messages {
		assert.NotEqual(t, "system", message.Role)
	}
}

// Before the first question there is no chat, and reading the panel must not
// create one — nor start the lesson.
func TestListAssistantMessagesIsEmptyBeforeTheFirstQuestion(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, secondLessonSlug)
	h.Assistant.StreamFunc = func(_ context.Context, _ []assistant.Turn, onDelta func(string) error) (assistant.Usage, error) {
		for _, delta := range []string{"Declare it ", "with ", "`let`."} {
			if err := onDelta(delta); err != nil {
				return assistant.Usage{}, err
			}
		}

		return assistant.Usage{InputTokens: 321, OutputTokens: 12}, nil
	}
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
	t.Parallel()
	h := testsupport.NewHarness(t)
	h.Assistant.StreamFunc = echoModel
	lesson := lessonBySlug(t, h, secondLessonSlug)

	prompt := askAssistantWith(t, h, lesson, &api.AssistantMessageInput{
		Message:  "Why does it fail?",
		UserCode: api.NewNilString("let x = 1;"),
		Output:   api.NewNilString("expected 2, got 1"),
	})

	system, question, found := strings.Cut(prompt, "\nuser: ")
	require.True(t, found, "a first question: the system prompt, then the question")
	assert.True(t, strings.HasPrefix(system, "system: "))
	assert.Contains(t, system, jsCourseSlug)
	assert.Contains(t, system, `"Variables"`)
	assert.Contains(t, system, "Answer in English")
	assert.Contains(t, question, "let x = 1;")
	assert.Contains(t, question, "expected 2, got 1")
	assert.Contains(t, question, "Why does it fail?")
}

// echoModel answers with the conversation it was given, one "role: content"
// line per turn, so a test reads the prompt through the response.
func echoModel(_ context.Context, turns []assistant.Turn, onDelta func(string) error) (assistant.Usage, error) {
	for _, turn := range turns {
		if err := onDelta(turn.Role + ": " + turn.Content + "\n"); err != nil {
			return assistant.Usage{}, err
		}
	}

	return assistant.Usage{
		InputTokens:  0,
		OutputTokens: 0,
	}, nil
}

// A learner reading the course in Russian is tutored in Russian, from the
// Russian lesson text.
func TestAskAssistantPromptFollowsTheLocale(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	h.Assistant.StreamFunc = echoModel
	testsupport.SpeakTo(h, "ru")
	lesson := lessonBySlug(t, h, firstLessonSlug)

	prompt := askAssistant(t, h, lesson, "Как вывести строку?")

	assert.Contains(t, prompt, "«Привет, мир!»")
	assert.Contains(t, prompt, "Отвечай на русском языке")
}

// A follow-up question carries the conversation so far, between the system
// prompt and the new question.
func TestAskAssistantSendsTheChatHistory(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	h.Assistant.StreamFunc = echoModel
	lesson := lessonBySlug(t, h, firstLessonSlug)

	prompt := askAssistant(t, h, lesson, "And a number?")

	history := []string{
		"user: How do I print a string?",
		"assistant: Use console.log to print a string.",
		"user: ",
	}

	rest := prompt
	for _, line := range history {
		at := strings.Index(rest, line)
		require.NotEqual(t, -1, at, "%q follows the previous turn", line)
		rest = rest[at+len(line):]
	}

	assert.Contains(t, rest, "And a number?")
}

// Seven questions a day (the fixture learner has asked seven today). Over the quota the operation answers 429 before the
// model is called, stores nothing, and the history reports the disabled state.
func TestAskAssistantRefusesOverTheDailyQuota(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserQuotaUsedToday)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, firstLessonSlug)
	messages := h.DB.AiMessage.Query().CountX(ctx)

	res, err := h.Client.CreateAssistantMessage(ctx,
		&api.AssistantMessageInput{Message: "One more?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)
	require.NoError(t, err)

	refused, ok := res.(*api.ProblemDetails)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, int32(http.StatusTooManyRequests), refused.Status)
	assert.Equal(t, http.StatusTooManyRequests, h.LastStatus())
	assert.Equal(t, messages, h.DB.AiMessage.Query().CountX(ctx))

	assert.True(t, listAssistantChat(t, h, lesson).QuotaExceeded)
}

// The quota is per day: yesterday's questions do not count against today.
func TestAskAssistantQuotaResetsDaily(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserQuotaUsedYesterday)
	lesson := lessonBySlug(t, h, firstLessonSlug)

	answer := askAssistant(t, h, lesson, "A new day?")

	assert.NotEmpty(t, answer, "yesterday's questions do not count, so the model answers")
}

// Any lesson may be asked about, in any order (ADR-0016): the learner has not
// reached the third lesson, and the assistant still answers and stores the turn.
func TestAskAssistantAnswersAboutALessonOutOfOrder(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, thirdLessonSlug)
	messages := h.DB.AiMessage.Query().CountX(t.Context())

	answer := askAssistant(t, h, lesson, "Help?")

	assert.NotEmpty(t, answer)
	assert.Greater(t, h.DB.AiMessage.Query().CountX(t.Context()), messages, "the turn is stored")
}

// When the model fails before its first token the learner gets an error
// status rather than an empty answer, and the question is not stored — it
// does not count against their quota.
func TestAskAssistantReportsAModelFailureBeforeStreaming(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, firstLessonSlug)
	h.Assistant.StreamFunc = func(context.Context, []assistant.Turn, func(string) error) (assistant.Usage, error) {
		return assistant.Usage{}, io.ErrUnexpectedEOF
	}
	messages := h.DB.AiMessage.Query().CountX(ctx)

	_, err := h.Client.CreateAssistantMessage(ctx,
		&api.AssistantMessageInput{Message: "Help?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, h.LastStatus())
	assert.Equal(t, messages, h.DB.AiMessage.Query().CountX(ctx))
}

// The answer is produced on a goroutine of its own, outside the request's
// panic recovery, so a panic in the model client would take the process down.
// It comes back as an ordinary failure instead.
func TestAskAssistantSurvivesAPanicBeforeStreaming(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, firstLessonSlug)
	h.Assistant.StreamFunc = func(context.Context, []assistant.Turn, func(string) error) (assistant.Usage, error) {
		panic("model client blew up")
	}

	_, err := h.Client.CreateAssistantMessage(ctx,
		&api.AssistantMessageInput{Message: "Help?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)

	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, h.LastStatus())
}

// A panic after the answer began ends the answer like any other failure
// mid-answer: the question is spent and nothing is stored for the answer.
func TestAskAssistantSurvivesAPanicMidStream(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, secondLessonSlug)
	h.Assistant.StreamFunc = func(_ context.Context, _ []assistant.Turn, onDelta func(string) error) (assistant.Usage, error) {
		if err := onDelta("Declare it "); err != nil {
			return assistant.Usage{}, err
		}

		panic("model client blew up")
	}

	res, err := h.Client.CreateAssistantMessage(ctx,
		&api.AssistantMessageInput{Message: "How do I declare a variable?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)
	require.NoError(t, err)

	_, _ = io.ReadAll(as[*api.CreateAssistantMessageOK](t, res))

	stored := h.DB.AiMessage.Query().Where(aimessage.AiChatID(chatFor(t, h, lesson).ID)).AllX(ctx)
	require.Len(t, stored, 1, "the question is spent, the broken answer is not stored")
	assert.Equal(t, "user", stored[0].Role)
}

// Once the answer has begun the question is spent, even if the answer never
// completes: otherwise cutting the connection just before the end would be a
// way round the quota. The broken answer itself is not stored.
func TestAskAssistantChargesAQuestionWhoseAnswerBroke(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, secondLessonSlug)
	h.Assistant.StreamFunc = func(_ context.Context, _ []assistant.Turn, onDelta func(string) error) (assistant.Usage, error) {
		if err := onDelta("Declare it "); err != nil {
			return assistant.Usage{}, err
		}

		return assistant.Usage{}, io.ErrUnexpectedEOF
	}

	_, _ = h.Client.CreateAssistantMessage(ctx,
		&api.AssistantMessageInput{Message: "How do I declare a variable?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)

	chat := chatFor(t, h, lesson)
	stored := h.DB.AiMessage.Query().Where(aimessage.AiChatID(chat.ID)).AllX(ctx)
	require.Len(t, stored, 1)
	assert.Equal(t, "user", stored[0].Role)
	assert.Equal(t, "How do I declare a variable?", *stored[0].Content)
}

// The chat is a signed-in feature: a visitor is refused by the contract.
func TestAskAssistantRequiresASignedInLearner(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	lesson := lessonBySlug(t, h, firstLessonSlug)

	_, err := h.Client.CreateAssistantMessage(t.Context(),
		&api.AssistantMessageInput{Message: "Help?"},
		api.CreateAssistantMessageParams{LessonId: int32(lesson.ID)},
	)
	requireProblem(t, err, http.StatusUnauthorized)
}

func listAssistantChat(t *testing.T, h *testsupport.Harness, lesson *ent.CourseLesson) *api.LessonAssistantChat {
	t.Helper()
	res, err := h.Client.ListAssistantMessages(t.Context(), api.ListAssistantMessagesParams{LessonId: int32(lesson.ID)})
	require.NoError(t, err)

	return res
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

	body, err := io.ReadAll(as[*api.CreateAssistantMessageOK](t, res))
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

func assistantMessagesCount(t *testing.T, h *testsupport.Harness) int {
	t.Helper()

	u := h.DB.User.GetX(t.Context(), actingUserID(t, h))
	if u.AssistantMessagesCount == nil {
		return 0
	}

	return *u.AssistantMessagesCount
}
