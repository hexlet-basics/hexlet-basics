package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/blogpostlike"
	"hexletbasics/ent/blogpostrelatedlanguageitem"
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
	// Promoted-course ids surface in display order (fixture order 0,1,2).
	assert.Equal(t, []int32{82481401, 207281424, 617920698}, post.RelatedCourseIds)
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

func TestAdminCreateBlogPost(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	post, err := h.Client.AdminCreateBlogPost(ctx, &api.BlogPostInput{
		Name:        api.NewNilString("Go rewrite notes"),
		Slug:        api.NewNilString("go-rewrite-notes"),
		Description: api.NewNilString("Notes from the migration"),
		State:       api.NewNilBlogPostState(api.BlogPostStatePublished),
		RichBody:    "<p>We moved to Go.</p>",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())

	assert.Equal(t, "go-rewrite-notes", post.Slug.Value)
	assert.Equal(t, api.BlogPostStatePublished, post.State.Value)
	assert.Equal(t, "<p>We moved to Go.</p>", post.RichBodyHtml)
	// The creator is the authenticated admin (the harness signs in as the
	// fixture admin), and the locale is pinned to the default until request
	// locale reaches handlers.
	assert.Equal(t, "alice@example.com", post.Creator.Email.Value)
	assert.Equal(t, "en", post.Locale.Value)
	assert.Equal(t, int32(0), post.LikesCount)
	assert.Equal(t, int32(0), post.RelatedCourseItemsCount)
}

func TestAdminUpdateBlogPost(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	post, err := h.Client.AdminUpdateBlogPost(ctx, &api.BlogPostInput{
		Name:        api.NewNilString("Second post, renamed"),
		Slug:        api.NewNilString("second-post"),
		Description: api.NilString{Null: true},
		State:       api.NewNilBlogPostState(api.BlogPostStatePublished),
		RichBody:    "<p>Now with a body.</p>",
	}, api.AdminUpdateBlogPostParams{ID: 6002})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, "Second post, renamed", post.Name.Value)
	assert.Equal(t, api.BlogPostStatePublished, post.State.Value)
	assert.Equal(t, "<p>Now with a body.</p>", post.RichBodyHtml)
	// Legacy assign_attributes semantics: a null nullable field clears the column.
	assert.True(t, post.Description.Null)
	// The stored locale is not re-stamped on update.
	assert.Equal(t, "en", post.Locale.Value)
}

func TestAdminUpdateBlogPostNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminUpdateBlogPost(ctx, &api.BlogPostInput{
		Name:        api.NilString{Null: true},
		Slug:        api.NilString{Null: true},
		Description: api.NilString{Null: true},
		State:       api.NilBlogPostState{Null: true},
		RichBody:    "<p>ghost</p>",
	}, api.AdminUpdateBlogPostParams{ID: 999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminDeleteBlogPost(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// 6001 carries both kinds of dependents: 2 likes and 3 related-course items.
	err := h.Client.AdminDeleteBlogPost(ctx, api.AdminDeleteBlogPostParams{ID: 6001})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	gone, err := h.DB.BlogPost.Query().Where(blogpost.IDEQ(6001)).Exist(ctx)
	require.NoError(t, err)
	assert.False(t, gone, "post should be deleted")

	likes, err := h.DB.BlogPostLike.Query().Where(blogpostlike.BlogPostIDEQ(6001)).Count(ctx)
	require.NoError(t, err)
	assert.Zero(t, likes, "likes should be deleted with the post")

	items, err := h.DB.BlogPostRelatedLanguageItem.Query().
		Where(blogpostrelatedlanguageitem.BlogPostIDEQ(6001)).Count(ctx)
	require.NoError(t, err)
	assert.Zero(t, items, "related-course items should be deleted with the post")
}

func TestAdminDeleteBlogPostNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	err := h.Client.AdminDeleteBlogPost(ctx, api.AdminDeleteBlogPostParams{ID: 999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminSetBlogPostRelatedCourses(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// Replace 6001's three promoted courses with two, in an explicit order.
	post, err := h.Client.AdminSetBlogPostRelatedCourses(ctx, &api.BlogPostRelatedCoursesInput{
		CourseIds: []int32{207281424, 82481401}, // ruby, javascript
	}, api.AdminSetBlogPostRelatedCoursesParams{ID: 6001})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, int32(2), post.RelatedCourseItemsCount)

	items, err := h.DB.BlogPostRelatedLanguageItem.Query().
		Where(blogpostrelatedlanguageitem.BlogPostIDEQ(6001)).
		Order(ent.Asc(blogpostrelatedlanguageitem.FieldOrder)).
		All(ctx)
	require.NoError(t, err)
	require.Len(t, items, 2)
	// Submission order is the display order.
	assert.Equal(t, 207281424, items[0].LanguageID)
	assert.Equal(t, 82481401, items[1].LanguageID)

	// Clearing the set drops the counter to zero.
	post, err = h.Client.AdminSetBlogPostRelatedCourses(ctx, &api.BlogPostRelatedCoursesInput{
		CourseIds: []int32{},
	}, api.AdminSetBlogPostRelatedCoursesParams{ID: 6001})
	require.NoError(t, err)
	assert.Equal(t, int32(0), post.RelatedCourseItemsCount)
}

func TestAdminSetBlogPostRelatedCoursesNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminSetBlogPostRelatedCourses(ctx, &api.BlogPostRelatedCoursesInput{
		CourseIds: []int32{82481401},
	}, api.AdminSetBlogPostRelatedCoursesParams{ID: 999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}
