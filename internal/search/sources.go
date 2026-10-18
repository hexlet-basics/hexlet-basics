package search

import (
	"context"
	"fmt"
	"html"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courseversion"
	"hexletbasics/ent/predicate"
	"hexletbasics/internal/blog"
	"hexletbasics/internal/htmltext"
	"hexletbasics/internal/landingpage"
)

// The expressions below are the document of each table, as
// catalogue_search_vector(config, title, summary, body) in the migration that
// indexes them: a query that spells it differently would be right but would not
// use the index, and the EXPLAIN test is what catches that.

// vector is the search document of a row, for the configuration.
func vector(config, title, summary, body string) string {
	return fmt.Sprintf("catalogue_search_vector('%s'::regconfig, %s, %s, %s)", config, title, summary, body)
}

// tsquery is the visitor's text as a web search query in the configuration.
func (m match) tsquery() string {
	return fmt.Sprintf("websearch_to_tsquery('%s'::regconfig, ?)", m.config)
}

// matches holds for a row whose document contains the words.
func (m match) matches(document string) *sql.Predicate {
	return sql.P(withArgs(document+" @@ "+m.tsquery(), m.text))
}

// rank scores how well a document matches (higher is better); title words
// outweigh summary words, which outweigh body words, by the weights the document
// gives them.
func (m match) rank(document string) sql.Querier {
	return sql.ExprFunc(withArgs("ts_rank("+document+", "+m.tsquery()+")::float8", m.text))
}

// headline is the stretch of text around the match. PostgreSQL marks the words
// with <b>, which tidy strips along with the markup of a post, so the snippet is
// plain text.
func (m match) headline(text string) sql.Querier {
	return sql.ExprFunc(withArgs(fmt.Sprintf("ts_headline('%s'::regconfig, %s, %s, 'MaxFragments=1, MinWords=12, MaxWords=35')", m.config, text, m.tsquery()), m.text))
}

// withArgs writes text with each `?` as a bound argument. The builder numbers
// them for PostgreSQL ($1, $2, ...), which a `?` left in a raw expression would
// not be.
func withArgs(text string, args ...any) func(*sql.Builder) {
	parts := strings.Split(text, "?")

	return func(b *sql.Builder) {
		for i, part := range parts {
			b.WriteString(part)

			if i < len(args) {
				b.Arg(args[i])
			}
		}
	}
}

// summaryAndBody is the text a snippet is taken from: the title is not, it is
// already the hit's title.
func summaryAndBody(summary, body string) string {
	return "coalesce(" + summary + ", '') || ' ' || coalesce(" + body + ", '')"
}

// withColumn selects one more column, computed from the selector, under a name
// the row's Value reads back.
func withColumn(name string, build func(*sql.Selector) sql.Querier) func(*sql.Selector) {
	return func(s *sql.Selector) { s.AppendSelectExprAs(build(s), name) }
}

// textOf reads a text column selected with withColumn.
func textOf(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return ""
	}
}

// htmlText is the plain text of a fragment that may be HTML.
func htmlText(fragment string) string {
	if !strings.ContainsAny(fragment, "<&") {
		return fragment
	}

	return html.UnescapeString(htmltext.PlainText(fragment))
}

// landingPages are the courses: the landing pages the catalogue lists.
type landingPages struct{ db *ent.Client }

func (landingPages) vector(m match, s *sql.Selector) string {
	return vector(m.config,
		"coalesce("+s.C(courselandingpage.FieldName)+", '') || ' ' || coalesce("+s.C(courselandingpage.FieldHeader)+", '')",
		s.C(courselandingpage.FieldDescription),
		s.C(courselandingpage.FieldOutcomesDescription))
}

// scope is what the course list shows: listed pages of completed courses, in the
// locale (landingpage.Listed), that have a slug to be addressed by.
//
// An unrestricted search (an editor's) keeps only the locale: a draft, archived
// or unlisted page, and the page of a course still in development, are records
// too.
func (l landingPages) scope(m match) []predicate.CourseLandingPage {
	visible := landingpage.Listed(m.locale)
	if m.unrestricted {
		visible = courselandingpage.LocaleEQ(m.locale)
	}

	return []predicate.CourseLandingPage{
		visible,
		courselandingpage.SlugNotNil(),
		func(s *sql.Selector) { s.Where(m.matches(l.vector(m, s))) },
	}
}

func (l landingPages) count(ctx context.Context, m match) (int, error) {
	result, err := l.db.CourseLandingPage.Query().Where(l.scope(m)...).Count(ctx)
	if err != nil {
		return result, errs.Wrapf(err, "count course landing pages")
	}

	return result, nil
}

func (l landingPages) top(ctx context.Context, m match, n int) ([]candidate, error) {
	rows, err := l.db.CourseLandingPage.Query().
		Where(l.scope(m)...).
		Select(courselandingpage.FieldID, courselandingpage.FieldSlug, courselandingpage.FieldName, courselandingpage.FieldHeader,
			courselandingpage.FieldState, courselandingpage.FieldCourseID).
		WithCourse(func(q *ent.CourseQuery) { q.Select(course.FieldID, course.FieldReadiness) }).
		Order(courselandingpage.OrderOption(ranked(func(s *sql.Selector) sql.Querier { return m.rank(l.vector(m, s)) }))).
		Limit(n).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list course landing pages")
	}

	out := make([]candidate, 0, len(rows))
	for _, row := range rows {
		rank, err := rankOf(row.Value)
		if err != nil {
			return nil, err
		}

		slug := *row.Slug
		out = append(out, candidate{rank: rank, typ: Course, id: row.ID, hit: Hit{
			Type:  Course,
			ID:    CourseID(slug),
			Title: titleOf(slug, row.Name, row.Header),
			Path:  pagePath(m.locale, "/languages/"+slug),
			State: lo.FromPtr(row.State),
			// The course's own readiness, which says whether the page is a course
			// that can be learned yet.
			Readiness: readinessOf(row.Edges.Course),
			Snippet:   "",
		}})
	}

	return out, nil
}

func (l landingPages) snippets(ctx context.Context, m match, ids []int) (map[int]string, error) {
	rows, err := l.db.CourseLandingPage.Query().
		Where(
			courselandingpage.IDIn(ids...),
			predicate.CourseLandingPage(withColumn("snippet", func(s *sql.Selector) sql.Querier {
				return m.headline(summaryAndBody(s.C(courselandingpage.FieldDescription), s.C(courselandingpage.FieldOutcomesDescription)))
			})),
		).
		Select(courselandingpage.FieldID).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list course landing pages")
	}

	return snippetsOf(rows, func(row *ent.CourseLandingPage) int { return row.ID })
}

// lessons are the lessons of the current build of the courses the catalogue
// lists, one translation per locale.
type lessons struct{ db *ent.Client }

func (lessons) vector(m match, s *sql.Selector) string {
	return vector(m.config,
		s.C(courselessontranslation.FieldName),
		s.C(courselessontranslation.FieldDescription),
		"coalesce("+s.C(courselessontranslation.FieldTheory)+", '') || ' ' || coalesce("+s.C(courselessontranslation.FieldInstructions)+", '')")
}

// scope is a lesson of the current Version of a completed course that the course
// list shows in the locale, translated into it: what the lesson page opens for a
// visitor who got there from the course page.
//
// An unrestricted search takes the lessons of any course, in the newest build the
// lesson has a translation in: the build being worked on, which is the current
// one once a build is released. A lesson has a translation per build, so without
// the choice one lesson would be one hit per build under the same id.
func (l lessons) scope(m match) []predicate.CourseLessonTranslation {
	where := []predicate.CourseLessonTranslation{courselessontranslation.LocaleEQ(m.locale)}
	if m.unrestricted {
		where = append(where, newestBuild(m.locale))
	} else {
		where = append(where, courselessontranslation.HasCourseVersionWith(courseversion.HasCurrentCoursesWith(
			landingpage.Completed(),
			course.HasLandingPagesWith(landingpage.Listed(m.locale)),
		)))
	}

	return append(where, func(s *sql.Selector) { s.Where(m.matches(l.vector(m, s))) })
}

// newestBuild holds for a translation that no later build of the same lesson
// translates into the same locale.
func newestBuild(locale string) predicate.CourseLessonTranslation {
	return func(s *sql.Selector) {
		later := sql.Table(courselessontranslation.Table).As("later")
		s.Where(sql.NotExists(sql.Select(later.C(courselessontranslation.FieldID)).From(later).Where(sql.And(
			sql.ColumnsEQ(later.C(courselessontranslation.FieldCourseLessonID), s.C(courselessontranslation.FieldCourseLessonID)),
			sql.EQ(later.C(courselessontranslation.FieldLocale), locale),
			sql.ColumnsGT(later.C(courselessontranslation.FieldCourseVersionID), s.C(courselessontranslation.FieldCourseVersionID)),
		))))
	}
}

func (l lessons) count(ctx context.Context, m match) (int, error) {
	result, err := l.db.CourseLessonTranslation.Query().Where(l.scope(m)...).Count(ctx)
	if err != nil {
		return result, errs.Wrapf(err, "count course lesson translations")
	}

	return result, nil
}

func (l lessons) top(ctx context.Context, m match, n int) ([]candidate, error) {
	rows, err := l.db.CourseLessonTranslation.Query().
		Where(l.scope(m)...).
		Select(courselessontranslation.FieldID, courselessontranslation.FieldName,
			courselessontranslation.FieldCourseID, courselessontranslation.FieldCourseLessonID).
		Order(courselessontranslation.OrderOption(ranked(func(s *sql.Selector) sql.Querier { return m.rank(l.vector(m, s)) }))).
		Limit(n).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list course lesson translations")
	}

	if len(rows) == 0 {
		return nil, nil
	}

	// The URL of a lesson names the course and the lesson by slug; both are one
	// small query for the whole window rather than an edge loaded per row.
	var lessonIDs, courseIDs []int
	for _, row := range rows {
		lessonIDs = append(lessonIDs, row.CourseLessonID)
		courseIDs = append(courseIDs, row.CourseID)
	}

	lessonRows, err := l.db.CourseLesson.Query().
		Where(courselesson.IDIn(lessonIDs...)).
		Select(courselesson.FieldID, courselesson.FieldSlug, courselesson.FieldState).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list course lessons")
	}

	courseRows, err := l.db.Course.Query().
		Where(course.IDIn(courseIDs...)).
		Select(course.FieldID, course.FieldSlug, course.FieldReadiness).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list courses")
	}

	lessonSlugs, lessonStates := map[int]string{}, map[int]string{}

	for _, row := range lessonRows {
		if row.Slug != nil {
			lessonSlugs[row.ID] = *row.Slug
			lessonStates[row.ID] = lo.FromPtr(row.State)
		}
	}

	courseSlugs, courseReadiness := map[int]string{}, map[int]string{}

	for _, row := range courseRows {
		if row.Slug != nil {
			courseSlugs[row.ID] = *row.Slug
			courseReadiness[row.ID] = lo.FromPtr(row.Readiness)
		}
	}

	out := make([]candidate, 0, len(rows))
	for _, row := range rows {
		lessonSlug, courseSlug := lessonSlugs[row.CourseLessonID], courseSlugs[row.CourseID]
		if lessonSlug == "" || courseSlug == "" {
			continue
		}

		rank, err := rankOf(row.Value)
		if err != nil {
			return nil, err
		}

		out = append(out, candidate{rank: rank, typ: Lesson, id: row.ID, hit: Hit{
			Type:  Lesson,
			ID:    LessonID(courseSlug, lessonSlug),
			Title: titleOf(lessonSlug, row.Name),
			Path:  pagePath(m.locale, "/languages/"+courseSlug+"/lessons/"+lessonSlug),
			State: lessonStates[row.CourseLessonID], Readiness: courseReadiness[row.CourseID],
			Snippet: "",
		}})
	}

	return out, nil
}

func (l lessons) snippets(ctx context.Context, m match, ids []int) (map[int]string, error) {
	rows, err := l.db.CourseLessonTranslation.Query().
		Where(
			courselessontranslation.IDIn(ids...),
			predicate.CourseLessonTranslation(withColumn("snippet", func(s *sql.Selector) sql.Querier {
				body := "coalesce(" + s.C(courselessontranslation.FieldTheory) + ", '') || ' ' || coalesce(" + s.C(courselessontranslation.FieldInstructions) + ", '')"

				return m.headline(summaryAndBody(s.C(courselessontranslation.FieldDescription), body))
			})),
		).
		Select(courselessontranslation.FieldID).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list course lesson translations")
	}

	return snippetsOf(rows, func(row *ent.CourseLessonTranslation) int { return row.ID })
}

// posts are the published blog posts of the locale.
type posts struct{ db *ent.Client }

func (posts) vector(m match, s *sql.Selector) string {
	return vector(m.config, s.C(blogpost.FieldName), s.C(blogpost.FieldDescription), s.C(blogpost.FieldRichBody))
}

// scope is the public blog's own scope (blog.Published), for posts that
// have a slug to be addressed by.
//
// An unrestricted search takes the posts of the locale in any state.
func (p posts) scope(m match) []predicate.BlogPost {
	visible := blog.Published(m.locale)
	if m.unrestricted {
		visible = blogpost.LocaleEQ(m.locale)
	}

	return []predicate.BlogPost{
		visible,
		blogpost.SlugNotNil(),
		func(s *sql.Selector) { s.Where(m.matches(p.vector(m, s))) },
	}
}

func (p posts) count(ctx context.Context, m match) (int, error) {
	result, err := p.db.BlogPost.Query().Where(p.scope(m)...).Count(ctx)
	if err != nil {
		return result, errs.Wrapf(err, "count blog posts")
	}

	return result, nil
}

func (p posts) top(ctx context.Context, m match, n int) ([]candidate, error) {
	rows, err := p.db.BlogPost.Query().
		Where(p.scope(m)...).
		Select(blogpost.FieldID, blogpost.FieldSlug, blogpost.FieldName, blogpost.FieldState).
		Order(blogpost.OrderOption(ranked(func(s *sql.Selector) sql.Querier { return m.rank(p.vector(m, s)) }))).
		Limit(n).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list blog posts")
	}

	out := make([]candidate, 0, len(rows))
	for _, row := range rows {
		rank, err := rankOf(row.Value)
		if err != nil {
			return nil, err
		}

		slug := *row.Slug
		out = append(out, candidate{rank: rank, typ: Post, id: row.ID, hit: Hit{
			Type:      Post,
			ID:        PostID(slug),
			Title:     titleOf(slug, row.Name),
			Path:      pagePath(m.locale, "/blog_posts/"+slug),
			State:     lo.FromPtr(row.State),
			Snippet:   "",
			Readiness: "",
		}})
	}

	return out, nil
}

func (p posts) snippets(ctx context.Context, m match, ids []int) (map[int]string, error) {
	rows, err := p.db.BlogPost.Query().
		Where(
			blogpost.IDIn(ids...),
			predicate.BlogPost(withColumn("snippet", func(s *sql.Selector) sql.Querier {
				return m.headline(summaryAndBody(s.C(blogpost.FieldDescription), s.C(blogpost.FieldRichBody)))
			})),
		).
		Select(blogpost.FieldID).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list blog posts")
	}

	return snippetsOf(rows, func(row *ent.BlogPost) int { return row.ID })
}

// snippetsOf reads the "snippet" column the query selected beside each row's
// id, by id; both sources that headline a match select it the same way.
func snippetsOf[R interface {
	Value(name string) (ent.Value, error)
}](rows []R, id func(R) int) (map[int]string, error) {
	out := make(map[int]string, len(rows))
	for _, row := range rows {
		value, err := row.Value("snippet")
		if err != nil {
			return nil, errs.Wrapf(err, "read snippet")
		}

		out[id(row)] = textOf(value)
	}

	return out, nil
}

// readinessOf is the readiness of a loaded course, empty when it was not loaded.
func readinessOf(c *ent.Course) string {
	if c == nil {
		return ""
	}

	return lo.FromPtr(c.Readiness)
}
