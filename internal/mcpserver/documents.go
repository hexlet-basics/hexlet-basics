package mcpserver

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"hexletbasics/internal/api"
	"hexletbasics/internal/htmltext"
	"hexletbasics/internal/search"
)

// document is a catalogue record as the fetch tool answers it: one text a model
// can read, made from what the website's own read of the record returns.
type document struct {
	id    string
	title string
	text  string
	// path is the record's page without the locale prefix and the origin.
	path     string
	metadata map[string]any
}

// nilString is the value of a nullable string, empty for null.
func nilString(s api.NilString) string {
	if s.Null {
		return ""
	}
	return strings.TrimSpace(s.Value)
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

// code writes a fenced block, unless the code is empty.
func (d *body) code(heading, language, source string) {
	if strings.TrimSpace(source) == "" {
		return
	}
	d.b.WriteString("## " + heading + "\n\n```" + language + "\n" + strings.TrimRight(source, "\n") + "\n```\n\n")
}

func (d *body) String() string { return strings.TrimSpace(d.b.String()) }

// lessonLine is one lesson of a course, with the id that fetches it.
func lessonLine(courseSlug string, l api.LearnerLesson) string {
	line := fmt.Sprintf("%d. %s (%s)", l.Position, firstOf(nilString(l.Name), l.Slug), search.LessonID(courseSlug, l.Slug))
	if description := nilString(l.Description); description != "" {
		line += ": " + description
	}
	return line
}

// courseDocument is the course page: the landing copy, the lessons in order and
// the questions and answers.
func courseDocument(ref search.Ref, view *api.CourseView) *document {
	course := view.Course
	landing := view.LandingPage.Value
	hasLanding := !view.LandingPage.Null
	title := firstOf(landing.Name, nilString(course.Name), ref.Slug)

	var d body
	d.title(title)
	if hasLanding {
		if header := strings.TrimSpace(landing.Header); header != title {
			d.section("", header)
		}
		d.section("", landing.Description)
		d.section(firstOf(nilString(landing.OutcomesHeader), "What you will learn"), nilString(landing.OutcomesDescription))
		d.section(nilString(landing.UsedInHeader), nilString(landing.UsedInDescription))
	}
	var lessons strings.Builder
	for _, l := range view.Lessons {
		lessons.WriteString(lessonLine(course.Slug, l) + "\n")
	}
	d.section("Lessons", lessons.String())
	var modules strings.Builder
	for _, m := range view.Modules {
		modules.WriteString("### " + firstOf(nilString(m.Name), "Module") + "\n\n")
		if description := nilString(m.Description); description != "" {
			modules.WriteString(description + "\n\n")
		}
		for _, l := range m.Lessons {
			modules.WriteString(lessonLine(course.Slug, l) + "\n")
		}
		modules.WriteString("\n")
	}
	d.section("Modules", modules.String())
	var questions strings.Builder
	for _, q := range view.QnaItems {
		questions.WriteString("**" + strings.TrimSpace(q.Question) + "**\n\n" + strings.TrimSpace(q.Answer) + "\n\n")
	}
	d.section("Questions and answers", questions.String())

	metadata := map[string]any{
		"type":         string(search.Course),
		"slug":         ref.Slug,
		"courseSlug":   course.Slug,
		"lessonsCount": len(view.Lessons),
	}
	if hasLanding {
		metadata["durationHours"] = landing.Duration
	}
	return &document{
		id: search.CourseID(ref.Slug), title: title, text: d.String(), path: "/languages/" + ref.Slug, metadata: metadata,
	}
}

// lessonDocument is the lesson page: the theory and the exercise, and the code
// the page gives every visitor (the starter code, the reference solution and the
// tests: #851 decided the tools add no rule of their own to the HTTP read).
func lessonDocument(ref search.Ref, view *api.CourseLessonView) *document {
	lesson := view.Lesson
	title := firstOf(nilString(lesson.Name), ref.Slug)
	language := nilString(lesson.ExerciseExtension)
	testLanguage := nilString(lesson.TestExtension)

	var d body
	d.title(title)
	d.section("", nilString(lesson.Description))
	d.section("Theory", nilString(lesson.Theory))
	d.section("Exercise", nilString(lesson.Instructions))
	var definitions strings.Builder
	for _, def := range lesson.Definitions {
		definitions.WriteString("- " + strings.TrimSpace(def.Name) + ": " + strings.TrimSpace(def.Description) + "\n")
	}
	d.section("Definitions", definitions.String())
	var tips strings.Builder
	for _, tip := range lesson.Tips {
		tips.WriteString("- " + strings.TrimSpace(tip) + "\n")
	}
	d.section("Tips", tips.String())
	d.code("Starter code", language, nilString(lesson.PreparedCode))
	d.code("Reference solution", language, nilString(lesson.OriginalCode))
	d.code("Tests", testLanguage, nilString(lesson.TestCode))

	metadata := map[string]any{
		"type":       string(search.Lesson),
		"slug":       ref.Slug,
		"courseSlug": ref.CourseSlug,
		"course":     nilString(lesson.Course.Name),
		"position":   lesson.NaturalOrder,
	}
	if source := nilString(lesson.SourceCodeUrl); source != "" {
		metadata["sourceUrl"] = source
	}
	return &document{
		id: search.LessonID(ref.CourseSlug, ref.Slug), title: title, text: d.String(),
		path: "/languages/" + ref.CourseSlug + "/lessons/" + ref.Slug, metadata: metadata,
	}
}

// postDocument is the blog post as plain text, with the courses it promotes.
func postDocument(ref search.Ref, view *api.BlogPostView) *document {
	post := view.Post
	title := firstOf(nilString(post.Name), ref.Slug)

	var d body
	d.title(title)
	d.section("", nilString(post.Description))
	d.section("", readable(post.RichBodyHtml))

	metadata := map[string]any{
		"type":               string(search.Post),
		"slug":               ref.Slug,
		"readingTimeMinutes": post.ReadingTime,
		"publishedAt":        post.CreatedAt.UTC().Format(time.RFC3339),
	}
	if author := firstOf(nilString(post.Creator.Name)); author != "" {
		metadata["author"] = author
	}
	if len(view.RelatedLandingPages) > 0 {
		related := make([]string, len(view.RelatedLandingPages))
		for i, page := range view.RelatedLandingPages {
			related[i] = search.CourseID(page.Slug)
		}
		metadata["relatedCourses"] = related
	}
	return &document{id: search.PostID(ref.Slug), title: title, text: d.String(), path: "/blog_posts/" + ref.Slug, metadata: metadata}
}

var (
	blockEnd   = strings.NewReplacer("</p>", "</p>\n\n", "</h1>", "</h1>\n\n", "</h2>", "</h2>\n\n", "</h3>", "</h3>\n\n", "</h4>", "</h4>\n\n", "</li>", "</li>\n", "</pre>", "</pre>\n\n", "<br>", "\n", "<br/>", "\n", "<br />", "\n")
	spaces     = regexp.MustCompile(`[ \t]+`)
	blankLines = regexp.MustCompile(`\n[ \t]*(\n[ \t]*)+`)
)

// readable is the text of a post's HTML with its paragraphs kept: the markup is
// dropped and each block ends a line, which is enough for a model and does not
// pretend to be Markdown.
func readable(fragment string) string {
	text := htmltext.PlainText(blockEnd.Replace(fragment))
	text = spaces.ReplaceAllString(text, " ")
	text = strings.ReplaceAll(text, " \n", "\n")
	text = strings.ReplaceAll(text, "\n ", "\n")
	return strings.TrimSpace(blankLines.ReplaceAllString(text, "\n\n"))
}
