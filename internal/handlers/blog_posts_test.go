package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// blog_posts.yml seeds these two.
const totalBlogPosts = 2

func TestAdminListBlogPosts(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListBlogPosts(ctx, api.AdminListBlogPostsParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(totalBlogPosts), page.Total)
	require.Len(t, page.Items, totalBlogPosts)
	// Newest first: ids strictly descending.
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}

	bySlug := map[string]api.BlogPost{}
	for _, it := range page.Items {
		bySlug[it.Slug.Value] = it
	}

	// The "full" post: rich body passed through, cover served from its blob key,
	// two likes, embedded creator, ru locale prefixed URL.
	full, ok := bySlug["hello-world"]
	require.True(t, ok, "post hello-world not found")
	assert.Equal(t, "Hello world", full.Name.Value)
	assert.Equal(t, api.BlogPostStatePublished, full.State.Value)
	assert.Equal(t, "<p>Hello <strong>world</strong> from the blog</p>", full.RichBodyHtml)
	// Legacy reading time floors (Ruby integer division); a 5-word post -> 0.
	assert.Equal(t, int32(0), full.ReadingTime)
	assert.Equal(t, int32(2), full.LikesCount)
	assert.Equal(t, int32(3), full.RelatedCourseItemsCount)
	assert.Equal(t, "https://code-basics.com/ru/blog_posts/hello-world", full.URL)
	assert.Equal(t, "alice@example.com", full.Creator.Email.Value)
	require.False(t, full.CoverThumbVariant.Null)
	assert.Equal(t, "http://localhost:3001/storage/blogcoverkey001", full.CoverThumbVariant.Value)
	// All three variants serve the same URL until image variants land (ADR-0005).
	assert.Equal(t, full.CoverThumbVariant.Value, full.CoverListVariant.Value)
	assert.Equal(t, full.CoverThumbVariant.Value, full.CoverMainVariant.Value)

	// The "empty" post: no rich body, no cover, no likes, en locale (no prefix).
	empty, ok := bySlug["second-post"]
	require.True(t, ok, "post second-post not found")
	assert.Equal(t, api.BlogPostStateDraft, empty.State.Value)
	assert.Equal(t, "", empty.RichBodyHtml)
	assert.Equal(t, int32(0), empty.ReadingTime)
	assert.Equal(t, int32(0), empty.LikesCount)
	assert.Equal(t, "https://code-basics.com/blog_posts/second-post", empty.URL)
	assert.True(t, empty.CoverThumbVariant.Null)
	assert.True(t, empty.CoverListVariant.Null)
	assert.True(t, empty.CoverMainVariant.Null)
}

func TestAdminGetBlogPost(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	post, err := h.Client.AdminGetBlogPost(ctx, api.AdminGetBlogPostParams{ID: 6001})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(6001), post.ID)
	assert.Equal(t, "hello-world", post.Slug.Value)
	assert.Equal(t, int32(2), post.LikesCount)
	assert.Equal(t, "http://localhost:3001/storage/blogcoverkey001", post.CoverMainVariant.Value)
}

func TestAdminGetBlogPostReturnsTrustedHTMLUnchanged(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	trustedHTML := `<script>alert("trusted")</script><iframe src="https://example.com/embed" onload="ready()"></iframe><p style="background-image:url(javascript:alert(1))">Body</p>`
	err := h.DB.BlogPost.UpdateOneID(6001).
		SetRichBody(trustedHTML).
		Exec(ctx)
	require.NoError(t, err)

	post, err := h.Client.AdminGetBlogPost(ctx, api.AdminGetBlogPostParams{ID: 6001})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, trustedHTML, post.RichBodyHtml)
}

func TestAdminGetBlogPostNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminGetBlogPost(ctx, api.AdminGetBlogPostParams{ID: 999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}
