package search_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	entgo "entgo.io/ent/dialect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/search"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport"
)

func newCatalogue(t *testing.T) (*search.Catalogue, *ent.Client) {
	t.Helper()
	db := testsupport.NewClient(t)
	return search.NewCatalogue(db), db
}

func ask(t *testing.T, c *search.Catalogue, locale, text string, types ...search.Type) ([]search.Hit, int) {
	t.Helper()
	hits, total, err := c.Search(t.Context(), search.SearchCatalogue{
		Text: text, Locale: locale, Types: types, Page: domain.Page{Limit: 50},
	})
	require.NoError(t, err)
	return hits, total
}

func ids(hits []search.Hit) []string {
	out := make([]string, len(hits))
	for i, hit := range hits {
		out[i] = hit.ID
	}
	return out
}

func TestACourseIsFoundByWhatItsLandingPageSays(t *testing.T) {
	c, _ := newCatalogue(t)

	hits, total := ask(t, c, "ru", "JavaScript", search.Course)

	assert.Contains(t, ids(hits), "course:javascript-ru")
	assert.Equal(t, len(hits), total)
	for _, hit := range hits {
		assert.Equal(t, search.Course, hit.Type)
		assert.Equal(t, "/ru/languages/"+strings.TrimPrefix(hit.ID, "course:"), hit.Path, "the Russian page carries the locale prefix")
		assert.NotEmpty(t, hit.Title)
	}
}

func TestTheEnglishSiteHasNoPrefixAndSearchesEnglishText(t *testing.T) {
	c, _ := newCatalogue(t)

	hits, _ := ask(t, c, "en", "Hello readers", search.Post)

	require.Equal(t, []string{"post:english-post"}, ids(hits))
	assert.Equal(t, "/blog_posts/english-post", hits[0].Path)
	assert.Equal(t, "English post", hits[0].Title)
}

func TestALessonIsFoundByItsTheoryAndNamesItsCourseAndSlug(t *testing.T) {
	c, _ := newCatalogue(t)

	hits, _ := ask(t, c, "ru", "первая программа", search.Lesson)

	require.Contains(t, ids(hits), "lesson:ruby/hello-world-ruby")
	for _, hit := range hits {
		if hit.ID == "lesson:ruby/hello-world-ruby" {
			assert.Equal(t, "/ru/languages/ruby/lessons/hello-world-ruby", hit.Path)
			assert.Equal(t, "Привет, мир!", hit.Title)
			assert.Contains(t, hit.Snippet, "первая программа")
		}
	}
}

func TestAPostIsFoundByItsBodyAndTheSnippetIsPlainText(t *testing.T) {
	c, _ := newCatalogue(t)

	hits, _ := ask(t, c, "ru", "шаблонным кодом", search.Post)

	require.Equal(t, []string{"post:learning-ruby"}, ids(hits))
	assert.Contains(t, hits[0].Snippet, "шаблонным")
	assert.NotContains(t, hits[0].Snippet, "<")
	assert.LessOrEqual(t, len([]rune(hits[0].Snippet)), 301)
}

// Each case is a record that matches and that a visitor cannot read, and the way
// to make it readable: the case proves the words do match it, so that its absence
// is the scope's doing and not a query that finds nothing.
func TestWhatAVisitorCannotReadIsNotFoundUntilItIsPublic(t *testing.T) {
	cases := map[string]struct {
		locale, text, hidden string
		hide, reveal         func(ctx context.Context, db *ent.Client)
	}{
		"an archived post": {locale: "ru", text: "Archived notes", hidden: "post:archived-notes", reveal: func(ctx context.Context, db *ent.Client) {
			db.BlogPost.Update().Where(blogpost.SlugEQ("archived-notes")).SetState("published").ExecX(ctx)
		}},
		"a draft post": {locale: "en", text: "Second post", hidden: "post:second-post", reveal: func(ctx context.Context, db *ent.Client) {
			db.BlogPost.Update().Where(blogpost.SlugEQ("second-post")).SetState("published").ExecX(ctx)
		}},
		"a draft landing page": {locale: "ru", text: "Node.js", hidden: "course:nodejs-ru", reveal: func(ctx context.Context, db *ent.Client) {
			db.CourseLandingPage.Update().Where(courselandingpage.SlugEQ("nodejs-ru")).SetState("published").ExecX(ctx)
		}},
		"an archived landing page": {locale: "ru", text: "архив", hidden: "course:javascript-ru-archived", reveal: func(ctx context.Context, db *ent.Client) {
			db.CourseLandingPage.Update().Where(courselandingpage.SlugEQ("javascript-ru-archived")).SetState("published").ExecX(ctx)
		}},
		"an unlisted landing page": {locale: "ru", text: "Ruby", hidden: "course:ruby-ru",
			hide: func(ctx context.Context, db *ent.Client) {
				db.CourseLandingPage.Update().Where(courselandingpage.SlugEQ("ruby-ru")).SetListed(false).ExecX(ctx)
			},
			reveal: func(ctx context.Context, db *ent.Client) {
				db.CourseLandingPage.Update().Where(courselandingpage.SlugEQ("ruby-ru")).SetListed(true).ExecX(ctx)
			}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c, db := newCatalogue(t)
			if tc.hide != nil {
				tc.hide(t.Context(), db)
			}

			before, _ := ask(t, c, tc.locale, tc.text)
			tc.reveal(t.Context(), db)
			after, _ := ask(t, c, tc.locale, tc.text)

			assert.NotContains(t, ids(before), tc.hidden)
			assert.Contains(t, ids(after), tc.hidden)
		})
	}
}

// The lessons of a course are as readable as the course is: when it is not
// completed, or its landing pages are not published, its lessons leave search
// with it.
func TestALessonIsFoundOnlyWhileItsCourseIsOnTheCourseList(t *testing.T) {
	c, db := newCatalogue(t)
	const lesson = "lesson:ruby/hello-world-ruby"
	found, _ := ask(t, c, "ru", "первая программа", search.Lesson)
	require.Contains(t, ids(found), lesson)

	db.Course.Update().Where(course.SlugEQ("ruby")).SetReadiness("in_development").ExecX(t.Context())
	hits, _ := ask(t, c, "ru", "первая программа", search.Lesson)
	assert.NotContains(t, ids(hits), lesson, "a course in development")

	db.Course.Update().Where(course.SlugEQ("ruby")).SetReadiness("completed").ExecX(t.Context())
	db.CourseLandingPage.Update().Where(courselandingpage.HasCourseWith(course.SlugEQ("ruby"))).SetState("draft").ExecX(t.Context())
	hits, _ = ask(t, c, "ru", "первая программа", search.Lesson)
	assert.NotContains(t, ids(hits), lesson, "a course whose landing pages are drafts")
}

// A lesson that only a retired build of the course had is gone from the course,
// so it is gone from search.
func TestALessonOfARetiredBuildIsNotFound(t *testing.T) {
	c, db := newCatalogue(t)
	retired := db.CourseLessonTranslation.Query().Where(courselessontranslation.NameEQ("Old Hello")).OnlyX(t.Context())
	require.NotNil(t, retired, "the fixtures keep a lesson of a retired build")

	hits, _ := ask(t, c, "en", "Old Hello")

	assert.Empty(t, hits)
}

func TestWhatIsInAnotherLocaleIsNotFound(t *testing.T) {
	c, _ := newCatalogue(t)

	russian, _ := ask(t, c, "ru", "Hello readers")
	english, _ := ask(t, c, "en", "Hello readers")

	assert.NotContains(t, ids(russian), "post:english-post")
	assert.Contains(t, ids(english), "post:english-post")
}

func TestTheTypeFilterNarrowsTheAnswerAndItsTotal(t *testing.T) {
	c, _ := newCatalogue(t)

	all, allTotal := ask(t, c, "ru", "Привет")
	lessons, lessonTotal := ask(t, c, "ru", "Привет", search.Lesson)

	require.NotEmpty(t, lessons)
	for _, hit := range lessons {
		assert.Equal(t, search.Lesson, hit.Type)
	}
	assert.Equal(t, len(lessons), lessonTotal)
	assert.GreaterOrEqual(t, allTotal, lessonTotal)
	assert.GreaterOrEqual(t, len(all), len(lessons))
}

func TestPagesOfAnAnswerAreTheAnswerInOrderWithoutGapsOrRepeats(t *testing.T) {
	c, _ := newCatalogue(t)
	whole, total := ask(t, c, "ru", "Привет")
	require.Greater(t, total, 3, "the fixtures have several lessons that greet")

	var paged []search.Hit
	for offset := 0; offset < total; offset += 3 {
		hits, pagedTotal, err := c.Search(t.Context(), search.SearchCatalogue{
			Text: "Привет", Locale: "ru", Page: domain.Page{Offset: offset, Limit: 3},
		})
		require.NoError(t, err)
		assert.Equal(t, total, pagedTotal)
		paged = append(paged, hits...)
	}

	assert.Equal(t, ids(whole), ids(paged))
	assert.Len(t, paged, total)
}

func TestBetterMatchesComeFirst(t *testing.T) {
	c, _ := newCatalogue(t)

	hits, _ := ask(t, c, "ru", "Ruby")

	require.NotEmpty(t, hits)
	assert.Contains(t, hits[0].Title, "Ruby", "a title that is the word outranks a body that mentions it")
}

func TestAnythingTheVisitorTypesIsASearchNotAnError(t *testing.T) {
	c, _ := newCatalogue(t)

	for _, text := range []string{`"unclosed`, "a & b | !c", "(", "-", "':*", "or or or", "   ", "\x00", strings.Repeat("слово ", 400)} {
		_, _, err := c.Search(t.Context(), search.SearchCatalogue{Text: text, Locale: "ru", Page: domain.Page{Limit: 10}})
		assert.NoError(t, err, "%q", text)
	}
}

func TestAPageBeyondTheWindowIsEmptyButTheTotalIsTold(t *testing.T) {
	c, _ := newCatalogue(t)

	hits, total, err := c.Search(t.Context(), search.SearchCatalogue{
		Text: "Привет", Locale: "ru", Page: domain.Page{Offset: 5000, Limit: 10},
	})

	require.NoError(t, err)
	assert.Empty(t, hits)
	assert.Positive(t, total)
}

func TestNoLearnerWrittenTextIsSearched(t *testing.T) {
	db := testsupport.NewClient(t)
	// The review's author wrote the body, and the site shows it on a course page.
	review := db.Review.GetX(t.Context(), fixtures.ID(fixtures.ReviewAdminLearnerRubyPublished))
	require.Equal(t, "Отличный курс, всё понятно", *review.Body)

	client, queries, _ := recording(t)
	c := search.NewCatalogue(client)

	hits, total := ask(t, c, "ru", "Отличный")

	assert.Empty(t, hits)
	assert.Zero(t, total)
	assert.NotEmpty(t, *queries)
	// And structurally: the tables of learner content are not read at all.
	for _, q := range *queries {
		for _, table := range []string{`"reviews"`, `"ai_messages"`, `"ai_chats"`, `"leads"`, `"users"`, `"book_requests"`, `"language_lesson_reviews"`} {
			assert.NotContains(t, q.query, table)
		}
	}
}

// capturedQuery is a statement the search sent, with its arguments.
type capturedQuery struct {
	query string
	args  []any
}

// recorder is a driver that remembers what was asked of it.
type recorder struct {
	entgo.Driver
	queries *[]capturedQuery
}

func (r recorder) Query(ctx context.Context, query string, args, v any) error {
	list, _ := args.([]any)
	*r.queries = append(*r.queries, capturedQuery{query: query, args: list})
	return r.Driver.Query(ctx, query, args, v)
}

// recording is an ent client over the test transaction that records every
// SELECT it sends, in the slice it returns, with the transaction itself.
func recording(t *testing.T) (*ent.Client, *[]capturedQuery, *sql.Tx) {
	t.Helper()
	_, tx := testsupport.NewClientWithTx(t)
	captured := &[]capturedQuery{}
	return ent.NewClient(ent.Driver(recorder{Driver: store.NewTxDriver(tx), queries: captured})), captured, tx
}

// The statements of a search use the GIN indexes of the migration, for each
// language and each table: the plan names the index. The fixtures are a handful
// of rows, for which PostgreSQL rightly prefers reading the table, so the test
// forbids a sequential scan for its own transaction; what it proves is that the
// query's expression is the one the index holds (a query that spelled the
// document differently would have no index to use, and would fail here), not
// that the planner chooses it at this size.
func TestEverySearchStatementCanUseItsIndex(t *testing.T) {
	// The indexes are named for the locale they serve.
	for _, locale := range []string{"ru", "en"} {
		suffix := locale
		t.Run(locale, func(t *testing.T) {
			client, queries, tx := recording(t)
			_, err := tx.ExecContext(t.Context(), "SET LOCAL enable_seqscan = off")
			require.NoError(t, err)
			// The lessons are reached through their locale and version, and for a
			// table this small that index alone is the cheaper way in, whatever the
			// text predicate says. Without it (in this transaction only) the choice
			// is between the text index and nothing.
			_, err = tx.ExecContext(t.Context(), `DROP INDEX index_language_lesson_version_infos_on_version_id_and_locale`)
			require.NoError(t, err)

			hits, _ := ask(t, search.NewCatalogue(client), locale, "hello")
			require.NotEmpty(t, hits)

			indexes := map[string]bool{
				"index_language_landing_pages_on_search_" + suffix:        false,
				"index_language_lesson_version_infos_on_search_" + suffix: false,
				"index_blog_posts_on_search_" + suffix:                    false,
			}
			for _, q := range *queries {
				if !strings.Contains(q.query, "@@") {
					continue
				}
				plan := explain(t, tx, q)
				for index := range indexes {
					if strings.Contains(plan, index) {
						indexes[index] = true
					}
				}
			}
			for index, used := range indexes {
				assert.True(t, used, "no statement of the search used %s", index)
			}
		})
	}
}
