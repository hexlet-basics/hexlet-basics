// Package relatedcourse owns a blog post's related (promoted) courses: the
// replace operation shared by the admin's hand-picked set and the AI pick, and
// the AI pick itself — the admin action enqueues one durable job per post and
// the worker-side Suggester asks the LLM for the closest courses (legacy
// FindRelatedCoursesForBlogPostJob).
package relatedcourse

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/samber/lo"
	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpostrelatedcourseitem"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/internal/htmltext"
	"hexletbasics/internal/landingpage"
	"hexletbasics/internal/store"
)

// maxTextLength caps the post text in the prompt, counted in runes with the
// "..." suffix included — Rails `String#truncate(2000)`, which legacy used.
const maxTextLength = 2000

// suggestInstructions is the legacy job's system prompt (including the
// heredoc's trailing newline), except the output sentence: the answer is now
// constrained by suggestSchema, so it asks for the `course_ids` field instead
// of a bare array. The prompt language is intentional, not the post's locale:
// it is what the legacy job sent for every post.
const suggestInstructions = `Ты — ассистент, который помогает подобрать курсы.
У тебя есть текст статьи блога и список курсов.
Выбери пять курсов подходящих под тему статьи в порядке приоритета. Первый наиболее близок, последний - наименее.
Верни идентификаторы курсов (по полю ` + "`id`" + `) в поле ` + "`course_ids`" + `, отсортированные по похожести.
Первыми должны идти наиболее близкие курсы.
`

// suggestSchema is the structured shape of the model's answer. The API demands
// an object at the root, hence the wrapper around the id list.
func suggestSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"course_ids": map[string]any{
				"type":  "array",
				"items": map[string]any{"type": "integer"},
			},
		},
		"required":             []string{"course_ids"},
		"additionalProperties": false,
	}
}

// suggestion is the decoded answer.
type suggestion struct {
	CourseIDs []int `json:"course_ids"`
}

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
// counter_culture drifted because delete_all skips its callbacks. A missing post
// is ent not-found, reported from inside the transaction.
//
// The delete, the insert and the counter run in one transaction. Legacy ran
// them bare, so a failed insert (an unknown course id fails the FK) left the
// post with no related courses at all; here the previous set survives. The
// admin set surfaces that failure as 409, and the Suggester passes only the
// ids it offered.
func (r *Replacer) Replace(ctx context.Context, postID int, courseIDs []int) error {
	courseIDs = lo.Uniq(courseIDs)

	return r.store.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		// First, so a missing post is ent not-found before anything is written
		// (the insert would otherwise fail its FK and read as a conflict). Left
		// unwrapped: the central handler maps ent's not-found to 404.
		err := db.BlogPost.UpdateOneID(postID).
			SetRelatedCourseItemsCount(len(courseIDs)).
			Exec(ctx)
		if err != nil {
			return errs.Wrapf(err, "save blog post")
		}

		_, err = db.BlogPostRelatedCourseItem.Delete().
			Where(blogpostrelatedcourseitem.BlogPostID(postID)).Exec(ctx)
		if err != nil {
			return oops.Wrapf(err, "clear related courses of post %d", postID)
		}

		builders := lo.Map(courseIDs, func(courseID int, i int) *ent.BlogPostRelatedCourseItemCreate {
			return db.BlogPostRelatedCourseItem.Create().
				SetBlogPostID(postID).
				SetCourseID(courseID).
				SetOrder(i)
		})

		_, err = db.BlogPostRelatedCourseItem.CreateBulk(builders...).Save(ctx)
		if err != nil {
			return oops.Wrapf(err, "store related courses of post %d", postID)
		}

		return nil
	})
}

// JSONCompleter is the LLM seam (implemented by assistant.OpenAI; tests fake
// it): one completion constrained to a JSON Schema, answered as JSON text.
type JSONCompleter interface {
	CompleteJSON(ctx context.Context, instructions, prompt, name string, schema any) (string, error)
}

// Suggester performs one suggestion job.
type Suggester struct {
	db       *ent.Client
	replacer *Replacer
	llm      JSONCompleter
	logger   *slog.Logger
}

// NewSuggester wires the worker-side dependencies.
func NewSuggester(db *ent.Client, replacer *Replacer, llm JSONCompleter, logger *slog.Logger) *Suggester {
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

	candidates, err := offeredCandidates(ctx, s.db, lo.FromPtr(post.Locale))
	if err != nil {
		return oops.Wrapf(err, "load course candidates for post %d", blogPostID)
	}

	prompt, err := suggestPrompt(post.RichBody, candidates)
	if err != nil {
		return err
	}

	answer, err := s.llm.CompleteJSON(ctx, suggestInstructions, prompt, "related_courses", suggestSchema())
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
	known := lo.Filter(courseIDs, func(id int, _ int) bool {
		_, ok := offered[id]

		return ok
	})

	return s.replacer.Replace(ctx, blogPostID, known)
}

// offeredCandidates are the courses the model may choose from: the published
// main landing pages of the post's locale. Not filtered by `listed`, like
// legacy; ordered by id only so the prompt is deterministic (legacy left the
// order to the database).
func offeredCandidates(ctx context.Context, db *ent.Client, locale string) ([]candidate, error) {
	pages, err := db.CourseLandingPage.Query().
		Where(
			courselandingpage.StateEQ(landingpage.StatePublished),
			courselandingpage.Main(true),
			courselandingpage.LocaleEQ(locale),
		).
		Order(ent.Asc(courselandingpage.FieldID)).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list course landing pages")
	}

	return lo.Map(pages, func(p *ent.CourseLandingPage, _ int) candidate {
		return candidate{ID: p.CourseID, Name: p.Header}
	}), nil
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

// parseCourseIDs decodes the structured answer into the ordered id list.
func parseCourseIDs(answer string) ([]int, error) {
	var parsed suggestion

	err := json.Unmarshal([]byte(answer), &parsed)
	if err != nil {
		return nil, oops.Wrapf(err, "decode course id list %q", answer)
	}

	return parsed.CourseIDs, nil
}
