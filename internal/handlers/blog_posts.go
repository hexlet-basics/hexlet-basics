package handlers

import (
	"context"
	"strings"

	"golang.org/x/net/html"

	"hexletbasics/ent"
	"hexletbasics/ent/activestorageattachment"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/blogpostlike"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
)

// Blog posts are READ-ONLY for now (list/get): create/update/delete/
// relatedCourses stay on the embedded UnimplementedHandler. `rich_body` is
// trusted editor HTML and is returned exactly as stored; the five production
// posts are migrated by hand, with no ActionText compatibility layer. The cover
// remains the single ActiveStorage blob served through `/storage/{key}`.

const (
	// wordsPerMinute mirrors the legacy reading-time divisor (BlogPostResource).
	wordsPerMinute = 260
	// blogPostRecordType/coverName key the polymorphic ActiveStorage rows to a
	// BlogPost (Rails stores the model class name).
	blogPostRecordType = "BlogPost"
	coverName          = "cover"
)

// AdminListBlogPosts returns a page of blog posts, newest first.
func (s *Server) AdminListBlogPosts(ctx context.Context, params api.AdminListBlogPostsParams) (api.AdminListBlogPostsRes, error) {
	page := newPagination(params.Page, params.PerPage)

	total, err := s.db.BlogPost.Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	posts, err := s.db.BlogPost.Query().
		WithCreator().
		Order(ent.Desc(blogpost.FieldID)).
		Offset(page.Offset()).
		Limit(page.Limit()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items, err := s.blogPostsToAPI(ctx, posts)
	if err != nil {
		return nil, err
	}

	return &api.BlogPostPage{Items: items, Total: int32(total), Page: page.Page, PerPage: page.PerPage}, nil
}

// AdminGetBlogPost returns one blog post; a missing id surfaces as ent
// not-found, which the central APIErrorHandler maps to 404.
func (s *Server) AdminGetBlogPost(ctx context.Context, params api.AdminGetBlogPostParams) (api.AdminGetBlogPostRes, error) {
	post, err := s.db.BlogPost.Query().
		Where(blogpost.IDEQ(int(params.ID))).
		WithCreator().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	items, err := s.blogPostsToAPI(ctx, []*ent.BlogPost{post})
	if err != nil {
		return nil, err
	}
	return &items[0], nil
}

// blogPostsToAPI assembles the read model for a set of posts. The cover and like
// count are not columns on blog_posts, so they are fetched in two batched
// queries keyed by post id rather than per-post, keeping the list handler off an
// N+1 path.
func (s *Server) blogPostsToAPI(ctx context.Context, posts []*ent.BlogPost) ([]api.BlogPost, error) {
	if len(posts) == 0 {
		return []api.BlogPost{}, nil
	}

	ids := make([]int, len(posts))
	for i, p := range posts {
		ids[i] = p.ID
	}

	coverKeyByPost, err := s.blogCoverKeys(ctx, ids)
	if err != nil {
		return nil, err
	}
	likesByPost, err := s.blogLikeCounts(ctx, ids)
	if err != nil {
		return nil, err
	}

	items := make([]api.BlogPost, len(posts))
	for i, p := range posts {
		items[i] = api.BlogPost{
			ID:                      int32(p.ID),
			Creator:                 s.conv.ToUser(p.Edges.Creator),
			Name:                    apiconv.NilStringFromPtr(p.Name),
			Slug:                    apiconv.NilStringFromPtr(p.Slug),
			Description:             apiconv.NilStringFromPtr(p.Description),
			State:                   nilBlogPostState(p.State),
			Locale:                  apiconv.NilStringFromPtr(p.Locale),
			URL:                     s.blogPostURL(p.Slug, p.Locale),
			RichBodyHtml:            p.RichBody,
			ReadingTime:             readingTime(p.RichBody),
			LikesCount:              int32(likesByPost[p.ID]),
			RelatedCourseItemsCount: int32(p.RelatedLanguageItemsCount),
			CoverThumbVariant:       s.coverVariant(coverKeyByPost[p.ID]),
			CoverListVariant:        s.coverVariant(coverKeyByPost[p.ID]),
			CoverMainVariant:        s.coverVariant(coverKeyByPost[p.ID]),
			CreatedAt:               p.CreatedAt,
		}
	}
	return items, nil
}

// blogCoverKeys returns each post's cover blob storage key by id (posts without
// a cover are simply absent from the map).
func (s *Server) blogCoverKeys(ctx context.Context, ids []int) (map[int]string, error) {
	rows, err := s.db.ActiveStorageAttachment.Query().
		Where(
			activestorageattachment.RecordType(blogPostRecordType),
			activestorageattachment.Name(coverName),
			activestorageattachment.RecordIDIn(ids...),
		).
		WithBlob().
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[int]string, len(rows))
	for _, r := range rows {
		if r.Edges.Blob != nil {
			out[r.RecordID] = r.Edges.Blob.Key
		}
	}
	return out, nil
}

// blogLikeCounts returns each post's like count by id via a single grouped
// aggregate over blog_post_likes.
func (s *Server) blogLikeCounts(ctx context.Context, ids []int) (map[int]int, error) {
	var rows []struct {
		BlogPostID int `json:"blog_post_id"`
		Count      int `json:"count"`
	}
	err := s.db.BlogPostLike.Query().
		Where(blogpostlike.BlogPostIDIn(ids...)).
		GroupBy(blogpostlike.FieldBlogPostID).
		Aggregate(ent.As(ent.Count(), "count")).
		Scan(ctx, &rows)
	if err != nil {
		return nil, err
	}
	out := make(map[int]int, len(rows))
	for _, r := range rows {
		out[r.BlogPostID] = r.Count
	}
	return out, nil
}

// nilBlogPostState bridges the nullable ent state string to the API enum.
func nilBlogPostState(v *string) api.NilBlogPostState {
	if v == nil {
		return api.NilBlogPostState{Null: true}
	}
	return api.NewNilBlogPostState(api.BlogPostState(*v))
}

// coverVariant builds the (nullable) cover URL from a blob key. The bytes are
// served by this server's own `/storage/{key}` path (PublicURL origin); the same
// URL fills all three variant fields until image variants land (ADR-0005).
func (s *Server) coverVariant(key string) api.NilString {
	if key == "" {
		return api.NilString{Null: true}
	}
	return api.NewNilString(s.cfg.PublicURL + "/storage/" + key)
}

// blogPostURL mirrors legacy blog_post_url(slug, suffix): the canonical site URL
// for the post, with a locale path segment for non-default locales (en, the
// default, has no prefix — see legacy AppHost.locale_for_url).
//
// Divergence, on purpose: legacy uses the REQUEST locale (I18n.locale, the admin
// UI's current language) as the suffix, not the post's own locale. There is no
// request locale at the ogen handler boundary yet (the known admin-locale design
// gap), so we substitute the post's own locale — arguably more correct for a
// canonical link, and revisitable once request locale reaches handlers.
func (s *Server) blogPostURL(slug, locale *string) string {
	var slugPart string
	if slug != nil {
		slugPart = *slug
	}
	prefix := ""
	if locale != nil && *locale != "" && *locale != "en" {
		prefix = "/" + *locale
	}
	return "https://" + s.cfg.AppHost + prefix + "/blog_posts/" + slugPart
}

// readingTime estimates minutes-to-read from the rich body, matching the legacy
// serializer EXACTLY for deploy parity. Legacy is `(split.size / 260).ceil` in
// Ruby, where both operands are Integers — so `/` already floors and `.ceil` is a
// no-op. That means it is floor division, not round-up: a <260-word post reads as
// 0. We reproduce that with integer division (not math.Ceil) on purpose.
func readingTime(richBodyHTML string) int32 {
	tokenizer := html.NewTokenizer(strings.NewReader(richBodyHTML))
	var plain strings.Builder
	for {
		switch tokenizer.Next() {
		case html.TextToken:
			plain.Write(tokenizer.Text())
			plain.WriteByte(' ')
		case html.ErrorToken:
			words := len(strings.Fields(plain.String()))
			return int32(words / wordsPerMinute)
		}
	}
}
