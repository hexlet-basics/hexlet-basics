package coursepage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/progress"
)

// The course-order rules of the read model, on values: what a module lists,
// where the visitor resumes and which lessons sit on either side of one.

func TestModuleLessonSlugsKeepOrderAndSkipLessonsWithoutAUrl(t *testing.T) {
	t.Parallel()

	version := func(slug *string) *ent.CourseLessonVersion {
		return &ent.CourseLessonVersion{Edges: ent.CourseLessonVersionEdges{
			Lesson: &ent.CourseLesson{Slug: slug},
		}}
	}
	module := &ent.CourseModuleTranslation{Edges: ent.CourseModuleTranslationEdges{
		Version: &ent.CourseModuleVersion{Edges: ent.CourseModuleVersionEdges{
			LessonVersions: []*ent.CourseLessonVersion{
				version(new("hello-world")),
				version(nil),
				{}, // a version whose lesson was not loaded
				version(new("variables")),
			},
		}},
	}}

	assert.Equal(t, []string{"hello-world", "variables"}, moduleLessonSlugs(module))
}

func TestModuleLessonSlugsOfAModuleWithoutAVersionAreEmptyNotNull(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{}, moduleLessonSlugs(nil))
	assert.Equal(t, []string{}, moduleLessonSlugs(&ent.CourseModuleTranslation{}))
}

func lessonsNamed(slugs ...string) []LearnerLesson {
	lessons := make([]LearnerLesson, len(slugs))
	for i, slug := range slugs {
		lessons[i] = LearnerLesson{ID: i + 1, Slug: slug, Name: new(slug), Description: nil, Position: i + 1, Finished: false}
	}

	return lessons
}

func TestAdjacentLessonsAreTheNeighboursInCourseOrder(t *testing.T) {
	t.Parallel()

	lessons := lessonsNamed("a", "b", "c")

	middle := adjacentOf("b", lessons)
	require.NotNil(t, middle.Previous)
	require.NotNil(t, middle.Next)
	assert.Equal(t, "a", middle.Previous.Slug)
	assert.Equal(t, "c", *middle.Next.Name, "named, not just addressed")

	first := adjacentOf("a", lessons)
	assert.Nil(t, first.Previous)
	require.NotNil(t, first.Next)
	assert.Equal(t, "b", first.Next.Slug)

	last := adjacentOf("c", lessons)
	require.NotNil(t, last.Previous)
	assert.Equal(t, "b", last.Previous.Slug)
	assert.Nil(t, last.Next)
}

// A lesson that is not in the list, or the only one in it, has no neighbours:
// both ends are nil rather than a guess.
func TestAdjacentLessonsOfALoneOrUnlistedLessonAreNone(t *testing.T) {
	t.Parallel()

	lone := adjacentOf("a", lessonsNamed("a"))
	assert.Nil(t, lone.Previous)
	assert.Nil(t, lone.Next)

	unlisted := adjacentOf("z", lessonsNamed("a", "b"))
	assert.Nil(t, unlisted.Previous)
	assert.Nil(t, unlisted.Next)
}

// The Next Lesson is named from the lessons of the locale; with none left, or a
// lesson the locale has no page for, there is nothing to point at. A course
// with no lessons has no standing at all.
func TestTheNextLessonIsNamedFromTheLessonsOfTheLocale(t *testing.T) {
	t.Parallel()

	lessons := lessonsNamed("a", "b")
	state := func(next *string) *progress.CourseState {
		return &progress.CourseState{State: "started", Completion: 50, NextLessonSlug: next, Lessons: []progress.LessonState{
			{Slug: "a", Position: 1, Finished: true}, {Slug: "b", Position: 2, Finished: false},
		}}
	}

	standing := courseProgressOf(state(new("b")), lessons)
	require.NotNil(t, standing.NextLesson)
	assert.Equal(t, 2, standing.NextLesson.ID)
	assert.Equal(t, 50, standing.Completion)

	assert.Nil(t, courseProgressOf(state(nil), lessons).NextLesson, "a finished course has none")
	assert.Nil(t, courseProgressOf(state(new("untranslated")), lessons).NextLesson)
	assert.Nil(t, courseProgressOf(&progress.CourseState{State: "", Completion: 0, NextLessonSlug: nil, Lessons: nil}, nil))
}
