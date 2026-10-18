package blog_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/blog"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

const origin = "https://blog.example.test/"

func newReader(t *testing.T) (*blog.Reader, *ent.Client) {
	t.Helper()
	db := testsupport.NewClient(t)

	return blog.NewReader(db, origin), db
}

func TestGetAssemblesCoverLikesAndRelatedCourses(t *testing.T) {
	reader, _ := newReader(t)

	post, err := reader.Get(t.Context(), fixture.ID(fixture.BlogPostHelloWorld))
	require.NoError(t, err)

	assert.Equal(t, "https://blog.example.test/api/storage/blogcoverkey001", post.CoverURL)
	assert.Equal(t, 2, post.LikesCount)
	assert.Equal(t, []int{
		fixture.ID(fixture.CourseJavascript),
		fixture.ID(fixture.CourseRuby),
		fixture.ID(fixture.CoursePython),
	}, post.RelatedCourseIDs, "in display order")
	assert.NotNil(t, post.Creator)
}

func TestGetOfAPostWithNothingAttachedHasNoCoverLikesOrCourses(t *testing.T) {
	reader, _ := newReader(t)

	post, err := reader.Get(t.Context(), fixture.ID(fixture.BlogPostSecondPost))
	require.NoError(t, err)

	assert.Empty(t, post.CoverURL)
	assert.Zero(t, post.LikesCount)
	assert.NotNil(t, post.RelatedCourseIDs, "an empty set, never nil")
	assert.Empty(t, post.RelatedCourseIDs)
	assert.Zero(t, post.ReadingTime, "an empty body reads in no time")
}

func TestGetOfAMissingPostIsNotFound(t *testing.T) {
	reader, _ := newReader(t)

	_, err := reader.Get(t.Context(), 1)
	assert.True(t, ent.IsNotFound(err), "got %v", err)
}

func TestAShortPostReadsInAMinute(t *testing.T) {
	reader, _ := newReader(t)

	post, err := reader.Get(t.Context(), fixture.ID(fixture.BlogPostHelloWorld))
	require.NoError(t, err)

	assert.Equal(t, 1, post.ReadingTime)
}

func TestReadingTimeRoundsUp(t *testing.T) {
	words := func(n int) string { return "<p>" + strings.Repeat("word ", n) + "</p>" }

	assert.Equal(t, 0, blog.ReadingTime(""))
	assert.Equal(t, 1, blog.ReadingTime(words(1)))
	assert.Equal(t, 1, blog.ReadingTime(words(260)))
	assert.Equal(t, 2, blog.ReadingTime(words(261)))
}

func TestAssembleKeepsTheOrderOfTheRows(t *testing.T) {
	reader, db := newReader(t)
	rows := db.BlogPost.Query().WithCreator().
		Where(blog.Published("ru")).AllX(t.Context())
	require.GreaterOrEqual(t, len(rows), 2)

	posts, err := reader.Assemble(t.Context(), rows)
	require.NoError(t, err)

	require.Len(t, posts, len(rows))

	for i, row := range rows {
		assert.Equal(t, row.ID, posts[i].ID)
	}
}

func TestLikeCountsAPersonOnce(t *testing.T) {
	reader, db := newReader(t)
	likes := blog.NewLiker(db)
	postID := fixture.ID(fixture.BlogPostSecondPost)
	userID := fixture.ID(fixture.UserAdminLearner)

	require.NoError(t, likes.Like(t.Context(), postID, userID))
	require.NoError(t, likes.Like(t.Context(), postID, userID), "a repeat like is not an error")

	post, err := reader.Get(t.Context(), postID)
	require.NoError(t, err)
	assert.Equal(t, 1, post.LikesCount)
}

func TestLikeOfAMissingPostIsNotFound(t *testing.T) {
	_, db := newReader(t)

	err := blog.NewLiker(db).Like(t.Context(), 1, fixture.ID(fixture.UserAdminLearner))
	assert.True(t, ent.IsNotFound(err), "got %v", err)
}
