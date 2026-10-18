package search

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/htmltext"
)

// FetchRecord is the read of one record an editor's search found, by the id the
// hit carried.
type FetchRecord struct {
	Ref Ref
	// Locale is the locale of the record: ids of courses, lessons and posts are
	// made of slugs, which are per locale.
	Locale string
}

// Record is a record of the catalogue as an editor reads it, whatever state it
// is in: one text a model can read, and what is known about how published it is.
type Record struct {
	Type  Type
	ID    string
	Title string
	Text  string
	// Path is the record's page on the website, with the locale prefix and
	// without an origin. A record that is not public has no page there yet.
	Path string
	// Slug is the slug the id carries; CourseSlug the course of a lesson.
	Slug       string
	CourseSlug string
	// State is the publication state of a landing page, lesson or post; Readiness
	// the readiness of the course (empty for a post).
	State     string
	Readiness string
	// Listed says whether a landing page is on the course list; nil for the
	// other kinds.
	Listed *bool
	// Build is the state of the build a lesson was read from; empty for the other
	// kinds.
	Build string
	// Position is the natural order of a lesson in its course; nil for the others.
	Position *int
}

// Fetch reads the record the id names, in any state: a draft, an archived or
// unlisted landing page, a course still in development, a lesson of the newest
// build, an unpublished post. It answers domain.ErrNotFound for an id that names
// nothing in the locale. Like an unrestricted Search it is for a caller who may
// read drafts; the state tells them what they are reading.
func (c *Catalogue) Fetch(ctx context.Context, cmd FetchRecord) (*Record, error) {
	var (
		rec *Record
		err error
	)

	switch cmd.Ref.Type {
	case Course:
		rec, err = c.fetchCourse(ctx, cmd)
	case Lesson:
		rec, err = c.fetchLesson(ctx, cmd)
	case Post:
		rec, err = c.fetchPost(ctx, cmd)
	default:
		return nil, errs.With("type", cmd.Ref.Type).Wrapf(ErrBadID, "fetch a record of unknown type")
	}

	if err != nil {
		return nil, domain.FromStore(err)
	}

	return rec, nil
}

func (c *Catalogue) fetchCourse(ctx context.Context, cmd FetchRecord) (*Record, error) {
	// The oldest page with the slug, so a slug that two pages share always names
	// the same one.
	page, err := c.db.CourseLandingPage.Query().
		Where(courselandingpage.SlugEQ(cmd.Ref.Slug), courselandingpage.LocaleEQ(cmd.Locale)).
		WithCourse().
		Order(ent.Asc(courselandingpage.FieldID)).
		First(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "find course landing page")
	}

	title := titleOf(cmd.Ref.Slug, page.Name, page.Header)

	var doc body
	doc.title(title)

	if header := strings.TrimSpace(lo.FromPtr(page.Header)); header != title {
		doc.section("", header)
	}

	doc.section("", lo.FromPtr(page.Description))
	doc.section(firstOf(lo.FromPtr(page.OutcomesHeader), "What you will learn"), lo.FromPtr(page.OutcomesDescription))
	doc.section(lo.FromPtr(page.UsedInHeader), lo.FromPtr(page.UsedInDescription))

	// The lessons of the newest build, in order, each with the id that fetches it.
	courseSlug := ""
	if page.Edges.Course != nil {
		courseSlug = lo.FromPtr(page.Edges.Course.Slug)
	}

	translations, err := c.db.CourseLessonTranslation.Query().
		Where(courselessontranslation.CourseID(page.CourseID), courselessontranslation.LocaleEQ(cmd.Locale), newestBuild(cmd.Locale)).
		WithLesson().
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list course lesson translations")
	}

	translations = slices.DeleteFunc(translations, func(t *ent.CourseLessonTranslation) bool {
		return t.Edges.Lesson == nil || t.Edges.Lesson.Slug == nil
	})
	position := func(t *ent.CourseLessonTranslation) int { return lo.FromPtr(t.Edges.Lesson.NaturalOrder) }

	slices.SortFunc(translations, func(a, b *ent.CourseLessonTranslation) int {
		if p := position(a) - position(b); p != 0 {
			return p
		}

		return a.ID - b.ID
	})

	var list strings.Builder

	for i, t := range translations {
		slug := *t.Edges.Lesson.Slug
		fmt.Fprintf(&list, "%d. %s (%s)", i+1, titleOf(slug, t.Name), LessonID(courseSlug, slug))

		if description := strings.TrimSpace(lo.FromPtr(t.Description)); description != "" {
			list.WriteString(": " + description)
		}

		list.WriteString("\n")
	}

	doc.section("Lessons", list.String())

	return &Record{
		Type: Course, ID: CourseID(cmd.Ref.Slug), Title: title, Text: doc.String(),
		Path: pagePath(cmd.Locale, "/languages/"+cmd.Ref.Slug), Slug: cmd.Ref.Slug, CourseSlug: courseSlug,
		State: lo.FromPtr(page.State), Readiness: readinessOf(page.Edges.Course), Listed: page.Listed,
		Build:    "",
		Position: nil,
	}, nil
}

func (c *Catalogue) fetchLesson(ctx context.Context, cmd FetchRecord) (*Record, error) {
	parent, err := c.db.Course.Query().Where(course.SlugEQ(cmd.Ref.CourseSlug)).Only(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "load course")
	}

	lesson, err := c.db.CourseLesson.Query().
		Where(courselesson.CourseID(parent.ID), courselesson.SlugEQ(cmd.Ref.Slug)).
		Only(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "load course lesson")
	}
	// The newest build the lesson is translated into, as the search finds it.
	courseLessonTranslation, err := c.db.CourseLessonTranslation.Query().
		Where(courselessontranslation.CourseLessonID(lesson.ID), courselessontranslation.LocaleEQ(cmd.Locale)).
		WithCourseVersion().
		Order(ent.Desc(courselessontranslation.FieldCourseVersionID)).
		First(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "find course lesson translation")
	}

	title := titleOf(cmd.Ref.Slug, courseLessonTranslation.Name)

	var doc body
	doc.title(title)
	doc.section("", lo.FromPtr(courseLessonTranslation.Description))
	doc.section("Theory", lo.FromPtr(courseLessonTranslation.Theory))
	doc.section("Exercise", lo.FromPtr(courseLessonTranslation.Instructions))

	rec := &Record{
		Type: Lesson, ID: LessonID(cmd.Ref.CourseSlug, cmd.Ref.Slug), Title: title, Text: doc.String(),
		Path: pagePath(cmd.Locale, "/languages/"+cmd.Ref.CourseSlug+"/lessons/"+cmd.Ref.Slug),
		Slug: cmd.Ref.Slug, CourseSlug: cmd.Ref.CourseSlug,
		State: lo.FromPtr(lesson.State), Readiness: lo.FromPtr(parent.Readiness), Position: lesson.NaturalOrder,
		Listed: nil,
		Build:  "",
	}
	if courseLessonTranslation.Edges.CourseVersion != nil {
		rec.Build = lo.FromPtr(courseLessonTranslation.Edges.CourseVersion.State)
	}

	return rec, nil
}

func (c *Catalogue) fetchPost(ctx context.Context, cmd FetchRecord) (*Record, error) {
	post, err := c.db.BlogPost.Query().
		Where(blogpost.SlugEQ(cmd.Ref.Slug), blogpost.LocaleEQ(cmd.Locale)).
		Order(ent.Asc(blogpost.FieldID)).
		First(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "find blog post")
	}

	title := titleOf(cmd.Ref.Slug, post.Name)

	var doc body
	doc.title(title)
	doc.section("", lo.FromPtr(post.Description))
	doc.section("", readable(post.RichBody))

	return &Record{
		Type: Post, ID: PostID(cmd.Ref.Slug), Title: title, Text: doc.String(),
		Path: pagePath(cmd.Locale, "/blog_posts/"+cmd.Ref.Slug), Slug: cmd.Ref.Slug, State: lo.FromPtr(post.State),
		CourseSlug: "",
		Readiness:  "",
		Listed:     nil,
		Build:      "",
		Position:   nil,
	}, nil
}

// firstOf is the first of the texts that is not blank.
func firstOf(texts ...string) string {
	for _, t := range texts {
		if strings.TrimSpace(t) != "" {
			return strings.TrimSpace(t)
		}
	}

	return ""
}

// body builds a Markdown text out of sections, leaving out the empty ones.
type body struct{ b strings.Builder }

func (d *body) String() string { return strings.TrimSpace(d.b.String()) }

// title writes the heading of the whole text.
func (d *body) title(text string) { d.b.WriteString("# " + text + "\n\n") }

// section writes a heading and its paragraph, unless the paragraph is empty.
func (d *body) section(heading, text string) {
	if text = strings.TrimSpace(text); text == "" {
		return
	}

	if heading != "" {
		d.b.WriteString("## " + heading + "\n\n")
	}

	d.b.WriteString(text + "\n\n")
}

var (
	spaces     = regexp.MustCompile(`[ \t]+`)
	blankLines = regexp.MustCompile(`\n[ \t]*(\n[ \t]*)+`)
)

// readable is the text of a post's HTML with its paragraphs kept: the markup is
// dropped and each block ends a line, which is enough for a model and does not
// pretend to be Markdown.
func readable(fragment string) string {
	text := htmltext.PlainText(blockEnd().Replace(fragment))
	text = spaces.ReplaceAllString(text, " ")
	text = strings.ReplaceAll(text, " \n", "\n")
	text = strings.ReplaceAll(text, "\n ", "\n")

	return strings.TrimSpace(blankLines.ReplaceAllString(text, "\n\n"))
}

func blockEnd() *strings.Replacer {
	return strings.NewReplacer(
		"</p>", "</p>\n\n",
		"</h1>", "</h1>\n\n",
		"</h2>", "</h2>\n\n",
		"</h3>", "</h3>\n\n",
		"</h4>", "</h4>\n\n",
		"</li>", "</li>\n",
		"</pre>", "</pre>\n\n",
		"<br>", "\n",
		"<br/>", "\n",
		"<br />", "\n",
	)
}
