// Package lessonreviews implements the AI lesson-review loop: the admin review
// actions enqueue one durable job per lesson version info, and the worker-side
// Reviewer summarizes the students' assistant questions for that info into
// `language_lesson_reviews.summary` (legacy ReviewLessonJob).
package lessonreviews

import (
	"context"
	"database/sql"
	"strings"

	"github.com/riverqueue/river"
	"github.com/samber/lo"
	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/aichat"
	"hexletbasics/ent/aimessage"
	"hexletbasics/ent/courselessonreview"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/jobs"
)

// maxQuestions caps how many recent student questions feed one summary,
// matching the legacy job's `limit(100)`.
const maxQuestions = 100

// reviewInstructions is the legacy ReviewLessonJob prompt, verbatim: the
// summaries are content-ops material read by Russian-speaking staff, so the
// prompt language is intentional, not tied to the lesson's locale.
const reviewInstructions = `Проанализируй вопросы, которые задают студенты ассистенту по уроку. Вопросы будут переданы ниже.
Суммируй основные претензиии и пожелания. Предложи как поменять урок.`

// Enqueuer schedules review jobs from the HTTP process (insert-only River
// client). No surrounding transaction: each job is independent and idempotent
// (the worker upserts by its natural key), so partial enqueue on failure is
// harmless — the admin just clicks again.
type Enqueuer struct {
	river *river.Client[*sql.Tx]
}

// NewEnqueuer wires the shared insert-only River client.
func NewEnqueuer(riverClient *river.Client[*sql.Tx]) *Enqueuer {
	return &Enqueuer{river: riverClient}
}

// EnqueueLessonReviews inserts one review job per lesson version info.
func (e *Enqueuer) EnqueueLessonReviews(ctx context.Context, lessonInfoIDs []int) error {
	if len(lessonInfoIDs) == 0 {
		return nil
	}
	params := make([]river.InsertManyParams, len(lessonInfoIDs))
	for i, id := range lessonInfoIDs {
		params[i] = river.InsertManyParams{Args: jobs.ReviewLessonArgs{LessonInfoID: id}}
	}
	if _, err := e.river.InsertMany(ctx, params); err != nil {
		return oops.Wrapf(err, "enqueue lesson reviews")
	}
	return nil
}

// Completer is the LLM seam (implemented by assistant.OpenAI; tests fake it).
type Completer interface {
	Complete(ctx context.Context, instructions, prompt string) (string, error)
}

// Reviewer performs one review job: summarize the lesson's recent student
// questions and upsert the per-(course, lesson, locale) review row.
type Reviewer struct {
	db  *ent.Client
	llm Completer
}

// NewReviewer wires the worker-side dependencies.
func NewReviewer(db *ent.Client, llm Completer) *Reviewer {
	return &Reviewer{db: db, llm: llm}
}

// ReviewLesson mirrors the legacy ReviewLessonJob#perform for one lesson
// version info. A lesson whose assistant chats hold no student questions gets
// an empty summary without an LLM call — the admin list hides those rows, but
// the write still marks the lesson as reviewed.
func (r *Reviewer) ReviewLesson(ctx context.Context, lessonInfoID int) error {
	info, err := r.db.CourseLessonTranslation.Get(ctx, lessonInfoID)
	if err != nil {
		return oops.Wrapf(err, "load lesson info %d", lessonInfoID)
	}

	// The question pool is per LESSON (chats hang off lesson members), shared
	// by every locale's summary — same as legacy `info.lesson.ai_messages`.
	questions, err := r.db.AiMessage.Query().
		Where(
			aimessage.RoleEQ("user"),
			aimessage.HasChatWith(aichat.HasMemberWith(
				lessonprogress.LessonID(info.LanguageLessonID),
			)),
		).
		Order(ent.Desc(aimessage.FieldID)).
		Limit(maxQuestions).
		All(ctx)
	if err != nil {
		return oops.Wrapf(err, "load student questions for lesson %d", info.LanguageLessonID)
	}

	summary := ""
	if len(questions) > 0 {
		summary, err = r.llm.Complete(ctx, reviewInstructions, reviewPrompt(info, questions))
		if err != nil {
			return oops.Wrapf(err, "summarize lesson %d", info.LanguageLessonID)
		}
	}

	return r.upsertReview(ctx, info, summary)
}

// reviewPrompt assembles the legacy user prompt: the lesson content followed by
// the raw question feed.
func reviewPrompt(info *ent.CourseLessonTranslation, questions []*ent.AiMessage) string {
	contents := make([]string, len(questions))
	for i, q := range questions {
		contents[i] = lo.FromPtr(q.Content)
	}
	return strings.Join([]string{
		"Урок (теория и упражнение): " + lo.FromPtr(info.Theory) + "\n\n" + lo.FromPtr(info.Instructions),
		"Вопросы пользователей: " + strings.Join(contents, "\n\n"),
	}, "\n\n")
}

// upsertReview writes the summary keyed by (course, lesson, locale) — the
// legacy find_or_initialize_by — pointing the row at the reviewed version/info.
func (r *Reviewer) upsertReview(ctx context.Context, info *ent.CourseLessonTranslation, summary string) error {
	existing, err := r.db.CourseLessonReview.Query().
		Where(
			courselessonreview.LanguageID(info.LanguageID),
			courselessonreview.LanguageLessonID(info.LanguageLessonID),
			courselessonreview.LocaleEQ(lo.FromPtr(info.Locale)),
		).
		Only(ctx)
	switch {
	case ent.IsNotFound(err):
		_, err = r.db.CourseLessonReview.Create().
			SetLanguageID(info.LanguageID).
			SetLanguageLessonID(info.LanguageLessonID).
			SetLocale(lo.FromPtr(info.Locale)).
			SetSummary(summary).
			SetLanguageLessonVersionID(info.VersionID).
			SetLanguageLessonVersionInfoID(info.ID).
			Save(ctx)
	case err != nil:
		return oops.Wrapf(err, "find review for lesson %d", info.LanguageLessonID)
	default:
		_, err = existing.Update().
			SetSummary(summary).
			SetLanguageLessonVersionID(info.VersionID).
			SetLanguageLessonVersionInfoID(info.ID).
			Save(ctx)
	}
	return oops.Wrapf(err, "save review for lesson %d", info.LanguageLessonID)
}
