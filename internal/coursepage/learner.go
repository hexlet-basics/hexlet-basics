package coursepage

import (
	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/internal/progress"
)

// LearnerLesson is a Lesson as the course page and the player list it
// (GLOSSARY, Learner Lesson): its page in the visitor's locale, their Position
// in course order and whether they finished it.
type LearnerLesson struct {
	ID          int
	Slug        string
	Name        *string
	Description *string
	// Position is zero for a lesson the progress state does not place yet.
	Position int
	Finished bool
}

// Module is one of the current Version's modules in the locale, with the
// Learner Lessons it lists in course order.
type Module struct {
	ID          int
	Name        *string
	Description *string
	Lessons     []LearnerLesson
}

// CourseProgress is where a visitor stands in one Course (GLOSSARY, Course
// Progress): its state, Completion and Next Lesson. State is empty for a
// visitor who has not started the course (Not started is the absence of an
// Enrollment, not a state).
type CourseProgress struct {
	State      string
	Completion int
	// NextLesson is where the visitor resumes; nil when they finished every
	// Current Lesson, and when the Next Lesson has no page in the locale.
	NextLesson *LearnerLesson
}

// AdjacentLessons are the Lessons on either side of one in course order
// (GLOSSARY, Adjacent Lessons), whatever the visitor finished; either is nil at
// the edge of the course.
type AdjacentLessons struct {
	Previous *LearnerLesson
	Next     *LearnerLesson
}

// learnerLessons joins the lessons of the current Version in the locale with
// the visitor's standing in each. infos are the translations in course order,
// each with its Lesson loaded (lessonInfos).
//
// A lesson the state does not know has no Position and is not finished: the
// state and the infos read the same Version, so that is a lesson added since
// the state was computed, not a case to fail on.
func learnerLessons(infos []*ent.CourseLessonTranslation, state *progress.CourseState) []LearnerLesson {
	standing := lo.KeyBy(state.Lessons, func(l progress.LessonState) string { return l.Slug })

	return lo.Map(infos, func(info *ent.CourseLessonTranslation, _ int) LearnerLesson {
		slug := ""
		if info.Edges.Lesson != nil {
			slug = lo.FromPtr(info.Edges.Lesson.Slug)
		}

		return LearnerLesson{
			ID:          info.CourseLessonID,
			Slug:        slug,
			Name:        info.Name,
			Description: info.Description,
			Position:    standing[slug].Position,
			Finished:    standing[slug].Finished,
		}
	})
}

// modulesOf assigns each module its lessons. rows are the module translations
// with their lesson versions and lessons loaded (curriculum.Modules), lessons
// already in course order. A module's lesson with no page in the locale is not
// listed, as in the flat list.
func modulesOf(rows []*ent.CourseModuleTranslation, lessons []LearnerLesson) []Module {
	bySlug := lo.KeyBy(lessons, func(l LearnerLesson) string { return l.Slug })

	return lo.Map(rows, func(row *ent.CourseModuleTranslation, _ int) Module {
		id := 0
		if row.Edges.Version != nil {
			id = row.Edges.Version.ModuleID
		}

		return Module{
			ID:          id,
			Name:        row.Name,
			Description: row.Description,
			Lessons: lo.FilterMap(moduleLessonSlugs(row), func(slug string, _ int) (LearnerLesson, bool) {
				lesson, ok := bySlug[slug]

				return lesson, ok
			}),
		}
	})
}

// moduleLessonSlugs are the slugs of a module's lessons in the order its
// version lists them.
func moduleLessonSlugs(row *ent.CourseModuleTranslation) []string {
	if row == nil || row.Edges.Version == nil {
		return []string{}
	}

	return lo.FilterMap(row.Edges.Version.Edges.LessonVersions, func(v *ent.CourseLessonVersion, _ int) (string, bool) {
		if v.Edges.Lesson == nil || v.Edges.Lesson.Slug == nil {
			return "", false
		}

		return *v.Edges.Lesson.Slug, true
	})
}

// courseProgressOf is the visitor's standing, with the Next Lesson named from
// the lessons of the locale. Nil for a course with no lessons, where there is
// nothing to stand in.
func courseProgressOf(state *progress.CourseState, lessons []LearnerLesson) *CourseProgress {
	if len(state.Lessons) == 0 {
		return nil
	}

	standing := &CourseProgress{State: state.State, Completion: state.Completion, NextLesson: nil}
	if state.NextLessonSlug == nil {
		return standing
	}

	if next, ok := lo.Find(lessons, func(l LearnerLesson) bool { return l.Slug == *state.NextLessonSlug }); ok {
		standing.NextLesson = &next
	}

	return standing
}

// adjacentOf are the lessons on either side of the one with slug, in the order
// lessons are given (course order).
func adjacentOf(slug string, lessons []LearnerLesson) AdjacentLessons {
	adjacent := AdjacentLessons{Previous: nil, Next: nil}

	index := lo.IndexOf(lo.Map(lessons, func(l LearnerLesson, _ int) string { return l.Slug }), slug)
	if index < 0 {
		return adjacent
	}

	if index > 0 {
		adjacent.Previous = &lessons[index-1]
	}

	if index < len(lessons)-1 {
		adjacent.Next = &lessons[index+1]
	}

	return adjacent
}
