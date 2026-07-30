package handlers

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/net/html"

	"hexletbasics/ent"
	"hexletbasics/ent/activestorageattachment"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/blogpostlike"
	"hexletbasics/ent/blogpostrelatedlanguageitem"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
)

// Blog posts (legacy `/admin/blog_posts`): full CRUD plus the related-courses
// set action. `rich_body` is trusted editor HTML stored and returned exactly as
// given (no ActionText compatibility layer). The cover remains the single
// ActiveStorage blob served through `/storage/{key}` on read; the input's
// coverAttachmentId is deferred until blob covers land (same deferral as the
// course cover).

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
	return s.getAdminBlogPost(ctx, int(params.ID))
}

// getAdminBlogPost is the shared read model every write handler echoes back.
func (s *Server) getAdminBlogPost(ctx context.Context, id int) (*api.BlogPost, error) {
	post, err := s.db.BlogPost.Query().
		Where(blogpost.IDEQ(id)).
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

// AdminCreateBlogPost creates a post. The input-to-builder mapping is the
// generated SetInput; the creator is the authenticated admin (legacy
// `blog_post.creator = current_user`). Locale: legacy stamps the admin UI's
// request locale, which does not reach ogen handlers yet (the known
// admin-locale design gap), so the default locale is pinned like the lesson
// lists do — thread the request locale through when it lands.
func (s *Server) AdminCreateBlogPost(ctx context.Context, req *api.BlogPostInput) (api.AdminCreateBlogPostRes, error) {
	creator, ok := AuthenticatedUser(ctx)
	if !ok {
		return nil, errUnauthenticated
	}

	row, err := s.db.BlogPost.Create().
		SetInput(req).
		SetCreatorID(creator.ID).
		SetLocale(defaultAdminLocale).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.getAdminBlogPost(ctx, row.ID)
}

// AdminUpdateBlogPost updates a post. The generated SetInput keeps the legacy
// assign_attributes semantics: a null nullable field clears the column. The
// stored locale is left untouched (divergence from legacy, which re-stamped
// the request locale on every update — an accident of its service signature,
// not behavior worth keeping).
func (s *Server) AdminUpdateBlogPost(ctx context.Context, req *api.BlogPostInput, params api.AdminUpdateBlogPostParams) (api.AdminUpdateBlogPostRes, error) {
	row, err := s.db.BlogPost.UpdateOneID(int(params.ID)).SetInput(req).Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.getAdminBlogPost(ctx, row.ID)
}

// AdminDeleteBlogPost removes a post with its dependents, mirroring the legacy
// model (`related_language_items dependent: :delete_all`, `likes dependent:
// :destroy`) — the FKs have no ON DELETE CASCADE, so the children go first.
func (s *Server) AdminDeleteBlogPost(ctx context.Context, params api.AdminDeleteBlogPostParams) (api.AdminDeleteBlogPostRes, error) {
	id := int(params.ID)

	// Ensure the post exists first (404 for a missing id, before any write).
	if _, err := s.db.BlogPost.Query().Where(blogpost.IDEQ(id)).Only(ctx); err != nil {
		return nil, err
	}

	if _, err := s.db.BlogPostRelatedLanguageItem.Delete().
		Where(blogpostrelatedlanguageitem.BlogPostID(id)).Exec(ctx); err != nil {
		return nil, err
	}
	if _, err := s.db.BlogPostLike.Delete().
		Where(blogpostlike.BlogPostID(id)).Exec(ctx); err != nil {
		return nil, err
	}
	if err := s.db.BlogPost.DeleteOneID(id).Exec(ctx); err != nil {
		return nil, err
	}
	return &api.AdminDeleteBlogPostNoContent{}, nil
}

// AdminSetBlogPostRelatedCourses replaces the post's promoted-courses set with
// the submitted ids, keeping their order as the display order and the counter
// column in sync. Deliberate divergence from legacy (which enqueued an AI
// "find related courses" job): the contract makes the selection explicit. An
// unknown course id fails the FK constraint, surfaced as 409 centrally.
func (s *Server) AdminSetBlogPostRelatedCourses(ctx context.Context, req *api.BlogPostRelatedCoursesInput, params api.AdminSetBlogPostRelatedCoursesParams) (api.AdminSetBlogPostRelatedCoursesRes, error) {
	id := int(params.ID)

	// Ensure the post exists first (404 for a missing id, before any write).
	if _, err := s.db.BlogPost.Query().Where(blogpost.IDEQ(id)).Only(ctx); err != nil {
		return nil, err
	}

	if _, err := s.db.BlogPostRelatedLanguageItem.Delete().
		Where(blogpostrelatedlanguageitem.BlogPostID(id)).Exec(ctx); err != nil {
		return nil, err
	}

	courseIDs := lo.Uniq(req.CourseIds)
	builders := make([]*ent.BlogPostRelatedLanguageItemCreate, len(courseIDs))
	for i, courseID := range courseIDs {
		builders[i] = s.db.BlogPostRelatedLanguageItem.Create().
			SetBlogPostID(id).
			SetLanguageID(int(courseID)).
			SetOrder(i)
	}
	if _, err := s.db.BlogPostRelatedLanguageItem.CreateBulk(builders...).Save(ctx); err != nil {
		return nil, err
	}

	if err := s.db.BlogPost.UpdateOneID(id).
		SetRelatedLanguageItemsCount(len(courseIDs)).
		Exec(ctx); err != nil {
		return nil, err
	}

	return s.getAdminBlogPost(ctx, id)
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
