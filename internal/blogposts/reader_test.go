package blogposts_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/blogposts"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

const origin = "https://blog.example.test/"

func newReader(t *testing.T) (*blogposts.Reader, *ent.Client) {
	t.Helper()
	db := testsupport.NewClient(t)
	return blogposts.NewReader(db, origin), db
}

func TestGetAssemblesCoverLikesAndRelatedCourses(t *testing.T) {
	reader, _ := newReader(t)

	post, err := reader.Get(t.Context(), fixtures.ID("blog_post_hello_world"))
	require.NoError(t, err)

	assert.Equal(t, "https://blog.example.test/api/storage/blogcoverkey001", post.CoverURL)
	assert.Equal(t, 2, post.LikesCount)
	assert.Equal(t, []int{
		fixtures.ID("course_javascript"),
		fixtures.ID("course_ruby"),
		fixtures.ID("course_python"),
	}, post.RelatedCourseIDs, "in display order")
	assert.NotNil(t, post.Creator)
}

func TestGetOfAPostWithNothingAttachedHasNoCoverLikesOrCourses(t *testing.T) {
	reader, _ := newReader(t)

	post, err := reader.Get(t.Context(), fixtures.ID("blog_post_second_post"))
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

	post, err := reader.Get(t.Context(), fixtures.ID("blog_post_hello_world"))
	require.NoError(t, err)

	assert.Equal(t, 1, post.ReadingTime)
}

func TestReadingTimeRoundsUp(t *testing.T) {
	words := func(n int) string { return "<p>" + strings.Repeat("word ", n) + "</p>" }

	assert.Equal(t, 0, blogposts.ReadingTime(""))
	assert.Equal(t, 1, blogposts.ReadingTime(words(1)))
	assert.Equal(t, 1, blogposts.ReadingTime(words(260)))
	assert.Equal(t, 2, blogposts.ReadingTime(words(261)))
}

func TestAssembleKeepsTheOrderOfTheRows(t *testing.T) {
	reader, db := newReader(t)
	rows := db.BlogPost.Query().WithCreator().
		Where(blogposts.Published("ru")).AllX(t.Context())
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
	likes := blogposts.NewLiker(db)
	postID := fixtures.ID("blog_post_second_post")
	userID := fixtures.ID("user_admin_learner")

	require.NoError(t, likes.Like(t.Context(), postID, userID))
	require.NoError(t, likes.Like(t.Context(), postID, userID), "a repeat like is not an error")

	post, err := reader.Get(t.Context(), postID)
	require.NoError(t, err)
	assert.Equal(t, 1, post.LikesCount)
}

func TestLikeOfAMissingPostIsNotFound(t *testing.T) {
	_, db := newReader(t)

	err := blogposts.NewLiker(db).Like(t.Context(), 1, fixtures.ID("user_admin_learner"))
	assert.True(t, ent.IsNotFound(err), "got %v", err)
}
