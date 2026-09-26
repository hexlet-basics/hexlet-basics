// Package relatedcourses owns a blog post's related (promoted) courses: the
// replace operation shared by the admin's hand-picked set and the AI pick, and
// the AI pick itself — the admin action enqueues one durable job per post and
// the worker-side Suggester asks the LLM for the closest courses (legacy
// FindRelatedCoursesForBlogPostJob).
package relatedcourses

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/riverqueue/river"
	"github.com/samber/lo"
	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpostrelatedcourseitem"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/htmltext"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/store"
)

// maxTextLength caps the post text in the prompt, counted in runes with the
// "..." suffix included — Rails `String#truncate(2000)`, which legacy used.
const maxTextLength = 2000

// suggestInstructions is the legacy job's system prompt, verbatim (including
// the heredoc's trailing newline). The prompt language is intentional, not the
// post's locale: it is what the legacy job sent for every post.
const suggestInstructions = `Ты — ассистент, который помогает подобрать курсы.
У тебя есть текст статьи блога и список курсов.
Выбери пять курсов подходящих под тему статьи в порядке приоритета. Первый наиболее близок, последний - наименее.
Верни результат в виде JSON-массива идентификаторов курсов (по полю ` + "`id`" + `) отсортированный по похожести.
Первыми должны идти наиболее близкие курсы.
`

// Replacer swaps a post's related-courses set: the admin's hand-picked set and
// the AI pick both go through it.
type Replacer struct {
	store store.Transactor
}

// NewReplacer builds the replacer over the application's transaction seam.
func NewReplacer(txStore store.Transactor) *Replacer {
	return &Replacer{store: txStore}
}

// Replace swaps the post's related-courses set for courseIDs, keeping their
// order as the display order and the counter column in sync. Duplicates keep
// their first position. Order is 0-based (legacy's job wrote 1-based, its
// admin set nothing); readers sort by it ascending, so only the relative order
// matters. The counter is written outright rather than counted up: legacy's
// counter_culture drifted because delete_all skips its callbacks.
//
// The delete, the insert and the counter run in one transaction. Legacy ran
// them bare, so a failed insert (an unknown course id fails the FK) left the
// post with no related courses at all; here the previous set survives. The
// admin set surfaces that failure as 409, and the Suggester passes only the
// ids it offered.
func (r *Replacer) Replace(ctx context.Context, postID int, courseIDs []int) error {
	courseIDs = lo.Uniq(courseIDs)
	return r.store.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		if _, err := db.BlogPostRelatedCourseItem.Delete().
			Where(blogpostrelatedcourseitem.BlogPostID(postID)).Exec(ctx); err != nil {
			return oops.Wrapf(err, "clear related courses of post %d", postID)
		}

		builders := lo.Map(courseIDs, func(courseID int, i int) *ent.BlogPostRelatedCourseItemCreate {
			return db.BlogPostRelatedCourseItem.Create().
				SetBlogPostID(postID).
				SetCourseID(courseID).
				SetOrder(i)
		})
		if _, err := db.BlogPostRelatedCourseItem.CreateBulk(builders...).Save(ctx); err != nil {
			return oops.Wrapf(err, "store related courses of post %d", postID)
		}

		if err := db.BlogPost.UpdateOneID(postID).
			SetRelatedCourseItemsCount(len(courseIDs)).
			Exec(ctx); err != nil {
			return oops.Wrapf(err, "count related courses of post %d", postID)
		}
		return nil
	})
}

// Enqueuer schedules the suggestion job from the HTTP process (insert-only
// River client). The job is idempotent — it replaces the whole set — so a
// repeated click only costs another LLM call.
type Enqueuer struct {
	river *river.Client[*sql.Tx]
}

// NewEnqueuer wires the shared insert-only River client.
func NewEnqueuer(riverClient *river.Client[*sql.Tx]) *Enqueuer {
	return &Enqueuer{river: riverClient}
}

// EnqueueRelatedCoursesSuggestion inserts one suggestion job for the post.
func (e *Enqueuer) EnqueueRelatedCoursesSuggestion(ctx context.Context, blogPostID int) error {
	if _, err := e.river.Insert(ctx, jobs.SuggestRelatedCoursesArgs{BlogPostID: blogPostID}, nil); err != nil {
		return oops.Wrapf(err, "enqueue related courses suggestion for post %d", blogPostID)
	}
	return nil
}

// Completer is the LLM seam (implemented by assistant.OpenAI; tests fake it).
type Completer interface {
	Complete(ctx context.Context, instructions, prompt string) (string, error)
}

// Suggester performs one suggestion job.
type Suggester struct {
	db       *ent.Client
	replacer *Replacer
	llm      Completer
	logger   *slog.Logger
}

// NewSuggester wires the worker-side dependencies.
func NewSuggester(db *ent.Client, replacer *Replacer, llm Completer, logger *slog.Logger) *Suggester {
	return &Suggester{db: db, replacer: replacer, llm: llm, logger: logger}
}

// candidate is one course offered to the LLM, serialized exactly like the
// legacy `{id: lp.language.id, name: lp.header}` hash.
type candidate struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}

// SuggestRelatedCourses mirrors legacy FindRelatedCoursesForBlogPostJob#perform:
// the post's text plus its locale's published main landing pages go in, an
// ordered id list comes out and replaces the related courses. An empty answer
// leaves the set untouched; an unparsable one fails the job so River retries
// it. Divergence, on purpose: ids the LLM invents (not among the offered
// courses) are dropped, where legacy failed on the FK mid-replace and left
// the post with a partial set.
func (s *Suggester) SuggestRelatedCourses(ctx context.Context, blogPostID int) error {
	post, err := s.db.BlogPost.Get(ctx, blogPostID)
	if err != nil {
		return oops.Wrapf(err, "load blog post %d", blogPostID)
	}

	// Not filtered by `listed`, like legacy; ordered by id only so the prompt
	// is deterministic (legacy left the order to the database).
	pages, err := s.db.LandingPage.Query().
		Where(
			landingpage.StateEQ("published"),
			landingpage.Main(true),
			landingpage.LocaleEQ(lo.FromPtr(post.Locale)),
		).
		Order(ent.Asc(landingpage.FieldID)).
		All(ctx)
	if err != nil {
		return oops.Wrapf(err, "load course candidates for post %d", blogPostID)
	}
	candidates := lo.Map(pages, func(p *ent.LandingPage, _ int) candidate {
		return candidate{ID: p.CourseID, Name: p.Header}
	})

	prompt, err := suggestPrompt(post.RichBody, candidates)
	if err != nil {
		return err
	}
	answer, err := s.llm.Complete(ctx, suggestInstructions, prompt)
	if err != nil {
		return oops.Wrapf(err, "suggest related courses for post %d", blogPostID)
	}

	courseIDs, err := parseCourseIDs(answer)
	if err != nil {
		return oops.Wrapf(err, "parse related courses for post %d", blogPostID)
	}
	s.logger.InfoContext(ctx, "related courses suggested", "blog_post_id", blogPostID, "course_ids", courseIDs)
	if len(courseIDs) == 0 {
		return nil
	}

	offered := lo.SliceToMap(candidates, func(c candidate) (int, struct{}) { return c.ID, struct{}{} })
	known := lo.Filter(courseIDs, func(id int, _ int) bool { _, ok := offered[id]; return ok })
	return s.replacer.Replace(ctx, blogPostID, known)
}

// suggestPrompt assembles the legacy user message: the truncated plain text,
// a blank line, then the candidates as JSON.
func suggestPrompt(richBody string, candidates []candidate) (string, error) {
	courses, err := json.Marshal(candidates)
	if err != nil {
		return "", oops.Wrapf(err, "encode course candidates")
	}
	// Whitespace collapses first so markup indentation does not eat the budget.
	text := lo.Ellipsis(strings.Join(strings.Fields(htmltext.PlainText(richBody)), " "), maxTextLength)
	return "Текст статьи: " + text + "\n\n" + "Список курсов: " + string(courses), nil
}

// parseCourseIDs reads the LLM's JSON id array, tolerating the Markdown code
// fence (```json … ```) models like to wrap it in — the legacy
// `gsub(/\A```(?:json)?|```\z/, "").strip`.
func parseCourseIDs(answer string) ([]int, error) {
	raw := strings.TrimSpace(answer)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")

	var ids []int
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &ids); err != nil {
		return nil, oops.Wrapf(err, "decode course id list %q", answer)
	}
	return ids, nil
}
