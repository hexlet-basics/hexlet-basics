package assistant

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/aichat"
	"hexletbasics/ent/aimessage"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/store"
)

// DailyQuota is how many questions a learner may ask per UTC day. The panel
// promises seven; legacy's policy compared `count <= 7` before storing and so
// actually admitted an eighth, which is the one place this port does not copy
// it.
const DailyQuota = 7

// ErrQuotaExceeded refuses a question once today's quota is used up.
var ErrQuotaExceeded = errors.New("assistant daily quota exceeded")

// The roles the chat stores and sends. Production rows also hold RubyLLM's
// system and tool roles; only user and assistant turns are the conversation.
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

// Turn is one message sent to the model.
type Turn struct {
	Role    string
	Content string
}

// Usage is the token count of one answer, stored on it for cost control.
type Usage struct {
	InputTokens  int
	OutputTokens int
}

// Streamer is the LLM seam: implemented by OpenAI, faked in tests with canned
// deltas.
type Streamer interface {
	Stream(ctx context.Context, turns []Turn, onDelta func(string) error) (Usage, error)
}

// LessonStarter starts the lesson a question is about (progress.Progress).
type LessonStarter interface {
	StartLesson(ctx context.Context, learner progress.Learner, lessonID int, locale string) (*progress.CourseState, error)
}

// Chat is the in-lesson assistant: a learner's conversation per lesson, over
// the legacy ai_chats/ai_messages tables, answered synchronously as a stream —
// no job, no out-of-band delivery (see the assistant port design).
type Chat struct {
	db      *ent.Client
	tx      store.Transactor
	lessons LessonStarter
	llm     Streamer
	i18n    *localization.Translator
	now     func() time.Time
}

// NewChat wires the chat's dependencies.
func NewChat(
	db *ent.Client,
	txStore store.Transactor,
	lessons LessonStarter,
	llm Streamer,
	translator *localization.Translator,
) *Chat {
	return &Chat{db: db, tx: txStore, lessons: lessons, llm: llm, i18n: translator, now: time.Now}
}

// History is a learner's chat on one lesson.
type History struct {
	// Messages are the user and assistant turns, oldest first, with the chat's
	// lesson progress, course and lesson loaded, and the lesson's infos in the
	// request locale — what the message projection reads.
	Messages      []*ent.AiMessage
	QuotaExceeded bool
}

// History reads the learner's chat on lessonID without creating anything: a
// lesson they never asked about has no chat yet, and an empty one reads the
// same. Returns an ent not-found error for a lesson that does not exist.
func (c *Chat) History(ctx context.Context, userID, lessonID int) (*History, error) {
	if _, err := c.db.CourseLesson.Get(ctx, lessonID); err != nil {
		return nil, err
	}

	locale := c.i18n.Locale(ctx)
	messages, err := c.db.AiMessage.Query().
		Where(
			aimessage.RoleIn(RoleUser, RoleAssistant),
			aimessage.HasChatWith(aichat.HasLessonProgressWith(
				lessonprogress.UserID(userID),
				lessonprogress.LessonID(lessonID),
			)),
		).
		WithChat(func(q *ent.AiChatQuery) {
			q.WithLessonProgress(func(q *ent.LessonProgressQuery) {
				q.WithCourse().WithLesson(func(q *ent.CourseLessonQuery) {
					q.WithInfos(func(q *ent.CourseLessonTranslationQuery) {
						q.Where(courselessontranslation.LocaleEQ(locale)).
							Order(ent.Asc(courselessontranslation.FieldID))
					})
				})
			})
		}).
		Order(ent.Asc(aimessage.FieldID)).
		All(ctx)
	if err != nil {
		return nil, oops.Wrapf(err, "load assistant chat for lesson %d", lessonID)
	}

	asked, err := c.askedToday(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &History{Messages: messages, QuotaExceeded: asked >= DailyQuota}, nil
}

// Question is one learner question with the editor state it is about.
type Question struct {
	UserID   int
	LessonID int
	Message  string
	UserCode string
	Output   string
}

// Ask streams the model's answer to q. Everything that decides the outcome —
// the quota, the gate, the lesson's text — is settled before it returns, and
// so is the model's first token: a model that fails straight away is an error
// here rather than an empty 200. The returned reader yields the answer as the
// model produces it; both turns are stored, with the answer's usage, before it
// reports EOF, so a finished response means a stored conversation.
//
// A failed or abandoned answer stores nothing: the learner did not get it, so
// it does not count against their quota.
//
// Returns ErrQuotaExceeded, progress.ErrLessonNotAvailable, or an ent
// not-found error for a lesson outside the course's current build.
func (c *Chat) Ask(ctx context.Context, q Question) (io.ReadCloser, error) {
	asked, err := c.askedToday(ctx, q.UserID)
	if err != nil {
		return nil, err
	}
	if asked >= DailyQuota {
		return nil, ErrQuotaExceeded
	}

	// Asking is doing the lesson: nothing else in the player starts it before
	// the first check, and the chat hangs off the learner's lesson progress.
	locale := c.i18n.Locale(ctx)
	learner := progress.Learner{UserID: q.UserID}
	if _, err := c.lessons.StartLesson(ctx, learner, q.LessonID, locale); err != nil {
		return nil, err
	}

	turns, lessonProgressID, err := c.turns(ctx, q, locale)
	if err != nil {
		return nil, err
	}

	answer, writer := io.Pipe()
	started := make(chan error, 1)
	go func() {
		var text strings.Builder
		first := true
		usage, err := c.llm.Stream(ctx, turns, func(delta string) error {
			if first {
				first = false
				started <- nil
			}
			text.WriteString(delta)
			_, err := writer.Write([]byte(delta))
			return err
		})
		if err != nil {
			if first {
				started <- err
			}
			writer.CloseWithError(err)
			return
		}
		if first {
			// An empty answer is still an answer.
			started <- nil
		}
		// The answer is complete; storing it must not depend on the learner's
		// connection outliving the last token.
		err = c.save(context.WithoutCancel(ctx), q, lessonProgressID, text.String(), usage)
		writer.CloseWithError(err)
	}()

	if err := <-started; err != nil {
		return nil, err
	}
	return answer, nil
}

// turns builds what the model is sent: the lesson's system prompt, the
// conversation so far, and the new question wrapped with the editor state.
func (c *Chat) turns(ctx context.Context, q Question, locale string) ([]Turn, int, error) {
	lp, err := c.db.LessonProgress.Query().
		Where(lessonprogress.UserID(q.UserID), lessonprogress.LessonID(q.LessonID)).
		WithCourse().
		Only(ctx)
	if err != nil {
		return nil, 0, oops.Wrapf(err, "load lesson progress for lesson %d", q.LessonID)
	}
	crs := lp.Edges.Course
	if crs.CurrentVersionID == nil {
		return nil, 0, &ent.NotFoundError{}
	}

	// The text the learner is reading: the current build, in their locale.
	info, err := c.db.CourseLessonTranslation.Query().
		Where(
			courselessontranslation.CourseLessonID(q.LessonID),
			courselessontranslation.CourseVersionID(*crs.CurrentVersionID),
			courselessontranslation.LocaleEQ(locale),
		).
		Only(ctx)
	if err != nil {
		return nil, 0, oops.Wrapf(err, "load lesson %d info in %s", q.LessonID, locale)
	}

	previous, err := c.db.AiMessage.Query().
		Where(
			aimessage.RoleIn(RoleUser, RoleAssistant),
			aimessage.HasChatWith(aichat.LessonProgressID(lp.ID)),
		).
		Order(ent.Asc(aimessage.FieldID)).
		All(ctx)
	if err != nil {
		return nil, 0, oops.Wrapf(err, "load assistant chat for lesson %d", q.LessonID)
	}

	turns := make([]Turn, 0, len(previous)+2)
	turns = append(turns, Turn{Role: RoleSystem, Content: c.i18n.TextWith(ctx, localization.AssistantInstructions, map[string]string{
		"Course":       lo.FromPtr(crs.Slug),
		"Lesson":       lo.FromPtr(info.Name),
		"Theory":       lo.FromPtr(info.Theory),
		"Instructions": lo.FromPtr(info.Instructions),
	})})
	for _, message := range previous {
		turns = append(turns, Turn{Role: message.Role, Content: lo.FromPtr(message.Content)})
	}
	turns = append(turns, Turn{Role: RoleUser, Content: c.i18n.TextWith(ctx, localization.AssistantQuestion, map[string]string{
		"Code":     q.UserCode,
		"Output":   q.Output,
		"Question": q.Message,
	})})
	return turns, lp.ID, nil
}

// save writes a completed exchange: the chat on first use, the question as
// asked (the wrapped prompt is rebuilt per request, so storing it would show
// the learner — and the lesson reviews — boilerplate), and the answer with its
// usage. The question is attributed to the learner and counted on
// users.assistant_messages_count, as legacy's counter cache did.
func (c *Chat) save(ctx context.Context, q Question, lessonProgressID int, answer string, usage Usage) error {
	return c.tx.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		chat, err := db.AiChat.Query().Where(aichat.LessonProgressID(lessonProgressID)).Only(ctx)
		if ent.IsNotFound(err) {
			chat, err = db.AiChat.Create().
				SetUserID(q.UserID).
				SetLessonProgressID(lessonProgressID).
				Save(ctx)
		}
		if err != nil {
			return oops.Wrapf(err, "find or create assistant chat")
		}

		if err := db.AiMessage.Create().
			SetAiChatID(chat.ID).
			SetRole(RoleUser).
			SetContent(q.Message).
			SetUserID(q.UserID).
			Exec(ctx); err != nil {
			return oops.Wrapf(err, "store assistant question")
		}
		if err := db.AiMessage.Create().
			SetAiChatID(chat.ID).
			SetRole(RoleAssistant).
			SetContent(answer).
			SetInputTokens(usage.InputTokens).
			SetOutputTokens(usage.OutputTokens).
			Exec(ctx); err != nil {
			return oops.Wrapf(err, "store assistant answer")
		}
		return oops.Wrapf(
			db.User.UpdateOneID(q.UserID).AddAssistantMessagesCount(1).Exec(ctx),
			"count assistant question",
		)
	})
}

// askedToday counts the learner's questions since the start of the UTC day —
// legacy's Date.current under the default UTC time zone.
func (c *Chat) askedToday(ctx context.Context, userID int) (int, error) {
	dayStart := c.now().UTC().Truncate(24 * time.Hour)
	count, err := c.db.AiMessage.Query().
		Where(
			aimessage.UserID(userID),
			aimessage.RoleEQ(RoleUser),
			aimessage.CreatedAtGTE(dayStart),
		).
		Count(ctx)
	return count, oops.Wrapf(err, "count today's assistant questions")
}
