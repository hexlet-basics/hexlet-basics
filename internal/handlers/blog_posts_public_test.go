package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/blogpostlike"
	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// blogSlugs projects posts onto their slugs, the business fact the order
// assertions compare (fixture ids are opaque).
func blogSlugs(posts []api.BlogPost) []string {
	slugs := make([]string, len(posts))
	for i, p := range posts {
		slugs[i] = p.Slug.Value
	}
	return slugs
}

// blogPostID resolves a fixture post's id by slug, for the operations the
// contract keys by id.
func blogPostID(t *testing.T, h *testsupport.Harness, slug string) int32 {
	t.Helper()
	row, err := h.DB.BlogPost.Query().Where(blogpost.SlugEQ(slug)).Only(t.Context())
	require.NoError(t, err)
	return int32(row.ID)
}

// The list is the request locale's published posts, newest first: the
// archived ru post and the en posts are not on the ru blog.
func TestListBlogPostsShowsPublishedPostsOfTheRequestLocale(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	page, err := h.Client.ListBlogPosts(t.Context(), api.ListBlogPostsParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, []string{"learning-ruby", "hello-world", "oldest-post"}, blogSlugs(page.Items))
	assert.Equal(t, int32(3), page.Total)
	assert.Equal(t, int32(20), page.PerPage, "legacy pagy default")
}

func TestListBlogPostsInEnglishListsOnlyEnglishPosts(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "en")

	page, err := h.Client.ListBlogPosts(t.Context(), api.ListBlogPostsParams{})
	require.NoError(t, err)

	// second-post is en but a draft.
	assert.Equal(t, []string{"english-post"}, blogSlugs(page.Items))
}

// The post page carries the stored body, the promoted courses' main ru landing
// pages in the post's order, and two other posts to read.
func TestGetBlogPostReturnsThePostPage(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	res, err := h.Client.GetBlogPost(t.Context(), api.GetBlogPostParams{Slug: "hello-world"})
	require.NoError(t, err)
	view, ok := res.(*api.BlogPostView)
	require.True(t, ok, "got %T", res)

	assert.Equal(t, "hello-world", view.Post.Slug.Value)
	assert.Equal(t, "<p>Hello <strong>world</strong> from the blog</p>", view.Post.RichBodyHtml)
	assert.Equal(t, "https://code-basics.com/ru/blog_posts/hello-world", view.Post.URL)
	assert.Equal(t, int32(2), view.Post.LikesCount)

	related := make([]string, len(view.RelatedLandingPages))
	for i, lp := range view.RelatedLandingPages {
		related[i] = lp.Course.Slug
		assert.Equal(t, "ru", lp.Locale.Value, "the request locale's landing page")
	}
	assert.Equal(t, []string{"javascript", "ruby", "python"}, related, "in the post's display order")

	assert.Equal(t, []string{"learning-ruby", "oldest-post"}, blogSlugs(view.RecommendedPosts),
		"newest first, never the post being read")
}

func TestGetBlogPostIsNotFoundOutsidePublishedPostsOfTheLocale(t *testing.T) {
	cases := map[string]struct{ locale, slug string }{
		"draft":          {"en", "second-post"},
		"archived":       {"ru", "archived-notes"},
		"another locale": {"en", "hello-world"},
		"missing":        {"ru", "no-such-post"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			h := testsupport.NewAnonymousHarness(t)
			testsupport.SpeakTo(h, c.locale)

			_, err := h.Client.GetBlogPost(t.Context(), api.GetBlogPostParams{Slug: c.slug})
			require.Error(t, err)
			assert.Equal(t, http.StatusNotFound, h.LastStatus())
		})
	}
}

// Next walks the ru blog from newer to older and ends with a 404 after the
// oldest post.
func TestGetNextBlogPostWalksToOlderPosts(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	walked := []string{}
	id := blogPostID(t, h, "learning-ruby")
	for {
		res, err := h.Client.GetNextBlogPost(t.Context(), api.GetNextBlogPostParams{ID: id})
		if err != nil {
			assert.Equal(t, http.StatusNotFound, h.LastStatus(), "the chain ends in a 404")
			break
		}
		post, ok := res.(*api.BlogPost)
		require.True(t, ok, "got %T", res)
		walked = append(walked, post.Slug.Value)
		id = post.ID
	}

	// archived-notes sits between learning-ruby and hello-world by id and is
	// skipped; the en posts never appear.
	assert.Equal(t, []string{"hello-world", "oldest-post"}, walked)
}

// A like counts once per user: the second one answers with the post unchanged.
func TestLikeBlogPostCountsOncePerUser(t *testing.T) {
	h := testsupport.NewHarness(t)
	id := blogPostID(t, h, "learning-ruby")

	first, err := h.Client.LikeBlogPost(t.Context(), api.LikeBlogPostParams{ID: id})
	require.NoError(t, err)
	liked, ok := first.(*api.BlogPost)
	require.True(t, ok, "got %T", first)
	assert.Equal(t, http.StatusCreated, h.LastStatus())
	assert.Equal(t, int32(1), liked.LikesCount)

	second, err := h.Client.LikeBlogPost(t.Context(), api.LikeBlogPostParams{ID: id})
	require.NoError(t, err)
	again, ok := second.(*api.BlogPost)
	require.True(t, ok, "got %T", second)
	assert.Equal(t, int32(1), again.LikesCount, "a repeat like changes nothing")

	userID, _ := testsupport.HarnessUser(h)
	likes, err := h.DB.BlogPostLike.Query().
		Where(blogpostlike.BlogPostID(int(id)), blogpostlike.UserID(userID)).
		Count(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, likes)
}

func TestLikeBlogPostRequiresSignIn(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	id := blogPostID(t, h, "learning-ruby")

	_, _ = h.Client.LikeBlogPost(t.Context(), api.LikeBlogPostParams{ID: id})
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus())

	likes, err := h.DB.BlogPostLike.Query().
		Where(blogpostlike.BlogPostID(int(id))).
		Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, likes)
}

func TestLikeBlogPostIsNotFoundForAMissingPost(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.LikeBlogPost(t.Context(), api.LikeBlogPostParams{ID: 1})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}
