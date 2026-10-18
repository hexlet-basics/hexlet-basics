package search

import (
	"errors"
	"strings"
)

// Type is the kind of record a hit is.
type Type string

const (
	// Course is a course's landing page in a locale.
	Course Type = "course"
	// Lesson is a lesson of a course's current version, in a locale.
	Lesson Type = "lesson"
	// Post is a published blog post.
	Post Type = "post"
)

// Types are every kind of record the search covers, in the order a tie between
// two hits is broken: a course before its lessons, and both before a post.
func Types() []Type {
	return []Type{Course, Lesson, Post}
}

// ErrBadID is a reference that is not an id this package makes.
var ErrBadID = errors.New("not a catalogue id")

// A hit's id is typed and made of what the website's own URL carries, so it is
// stable for as long as the page is: `course:<landing page slug>` (the slug the
// course page takes, which is per locale), `lesson:<course slug>/<lesson slug>`
// (the course's own slug, as the lesson page's URL has it) and `post:<slug>`. An
// assistant cannot build one by guessing and nothing in it is a database id.

// CourseID is the id of the course page with the landing page slug.
func CourseID(landingSlug string) string { return string(Course) + ":" + landingSlug }

// LessonID is the id of a lesson of the course with courseSlug.
func LessonID(courseSlug, lessonSlug string) string {
	return string(Lesson) + ":" + courseSlug + "/" + lessonSlug
}

// PostID is the id of a blog post.
func PostID(slug string) string { return string(Post) + ":" + slug }

// Ref is a parsed id: what it names, without a way to say whether a visitor may
// read it (the reader of the record decides that, as for any other request).
type Ref struct {
	Type Type
	// Slug is the landing page slug of a course, the lesson slug of a lesson and
	// the slug of a post.
	Slug string
	// CourseSlug is the course a lesson belongs to; empty for the other types.
	CourseSlug string
}

// ParseID reads an id made by CourseID, LessonID or PostID; anything else,
// including an empty part, is ErrBadID.
func ParseID(id string) (Ref, error) {
	kind, rest, ok := strings.Cut(id, ":")
	if !ok || rest == "" {
		return Ref{}, ErrBadID
	}

	switch Type(kind) {
	case Course, Post:
		if strings.Contains(rest, "/") {
			return Ref{}, ErrBadID
		}

		return Ref{
			Type:       Type(kind),
			Slug:       rest,
			CourseSlug: "",
		}, nil
	case Lesson:
		course, lesson, ok := strings.Cut(rest, "/")
		if !ok || course == "" || lesson == "" || strings.Contains(lesson, "/") {
			return Ref{}, ErrBadID
		}

		return Ref{Type: Lesson, CourseSlug: course, Slug: lesson}, nil
	default:
		return Ref{}, ErrBadID
	}
}
