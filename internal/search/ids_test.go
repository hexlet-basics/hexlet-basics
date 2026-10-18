package search_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/search"
)

func TestAnIdParsesBackToWhatMadeIt(t *testing.T) {
	course, err := search.ParseID(search.CourseID("javascript-ru"))
	require.NoError(t, err)
	lesson, err := search.ParseID(search.LessonID("javascript", "hello-world"))
	require.NoError(t, err)
	post, err := search.ParseID(search.PostID("hello-world"))
	require.NoError(t, err)

	assert.Equal(t, search.Ref{
		Type:       search.Course,
		Slug:       "javascript-ru",
		CourseSlug: "",
	}, course)
	assert.Equal(t, search.Ref{Type: search.Lesson, CourseSlug: "javascript", Slug: "hello-world"}, lesson)
	assert.Equal(t, search.Ref{
		Type:       search.Post,
		Slug:       "hello-world",
		CourseSlug: "",
	}, post)
}

func TestWhatIsNotAnIdIsRefused(t *testing.T) {
	for _, id := range []string{
		"", "javascript-ru", "course:", ":javascript-ru", "course:a/b", "lesson:javascript", "lesson:/hello",
		"lesson:javascript/", "lesson:a/b/c", "post:a/b", "page:home",
	} {
		_, err := search.ParseID(id)
		assert.ErrorIs(t, err, search.ErrBadID, "%q", id)
	}
}
