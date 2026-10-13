package handlers

import (
	"context"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/blogpostlike"
	"hexletbasics/internal/api"
)

// The public blog (legacy `Web::BlogPostsController`, `Web::BlogPosts::
// LikesController`, `Api::BlogPostsController#next`). Every read is scoped the
// way legacy scoped it — published posts in the request locale, newest (highest
// id) first — so a draft, an archived post or a post in another language is
// simply not found. The like is the exception: legacy looked the post up by slug
// alone, so any existing post can be liked.

// ListBlogPosts returns a page of published posts in the request locale,
// newest first (legacy `order(id: :desc)`, pagy's 20 per page).
func (s *Server) ListBlogPosts(ctx context.Context, params api.ListBlogPostsParams) (*api.BlogPostPage, error) {
	return listPageWith(ctx, params.Page, params.PerPage,
		func() *ent.BlogPostQuery { return s.blogPublic.Query(s.i18n.Locale(ctx)) },
		s.blogPostsToAPI,
		func(items []api.BlogPost, total, page, perPage int32) *api.BlogPostPage {
			return &api.BlogPostPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// GetBlogPost returns the post page: the post, two other posts to read next and
// the courses it promotes. Keyed by slug, so the id-keyed getOne helper does not
// fit; a miss is ent not-found, which the central handler maps to 404.
func (s *Server) GetBlogPost(ctx context.Context, params api.GetBlogPostParams) (api.GetBlogPostRes, error) {
	locale := s.i18n.Locale(ctx)
	row, err := s.blogPublic.BySlug(ctx, locale, params.Slug)
	if err != nil {
		return nil, err
	}
	recommendedRows, err := s.blogPublic.Recommended(ctx, locale, row)
	if err != nil {
		return nil, err
	}

	posts, err := s.blogPostsToAPI(ctx, append([]*ent.BlogPost{row}, recommendedRows...))
	if err != nil {
		return nil, err
	}
	post := posts[0]

	related, err := s.relatedLandingPages(ctx, post.RelatedCourseIds)
	if err != nil {
		return nil, err
	}

	return &api.BlogPostView{
		Post:                post,
		RecommendedPosts:    posts[1:],
		RelatedLandingPages: related,
	}, nil
}

// relatedLandingPages is legacy `related_main_language_landing_pages`: for each
// promoted course, its main landing page in the request locale; a course without
// one is skipped. Legacy's `has_one :main_landing_page` filtered only on main and
// locale, so a draft page could stand in; the public Web scope is applied here
// instead of linking readers to a page they cannot open. The course ids arrive already in the post's display order
// (blogRelatedCourseIDs orders them in SQL), so the pages are put back in that
// order here rather than joining the items table again for a handful of rows.
// A course with several main pages in one locale (legacy `has_one` picked one
// arbitrarily) keeps its oldest.
func (s *Server) relatedLandingPages(ctx context.Context, courseIDs []int32) ([]api.CourseCatalogItem, error) {
	if len(courseIDs) == 0 {
		return []api.CourseCatalogItem{}, nil
	}

	ids := lo.Map(courseIDs, func(id int32, _ int) int { return int(id) })

	byCourse, err := s.catalog.MainPages(ctx, s.i18n.Locale(ctx), ids...)
	if err != nil {
		return nil, err
	}

	ordered := make([]*ent.CourseLandingPage, 0, len(byCourse))
	for _, id := range lo.Uniq(ids) {
		if p, ok := byCourse[id]; ok {
			ordered = append(ordered, p)
		}
	}
	return s.conv.ToCatalogItems(ordered), nil
}

// GetNextBlogPost is legacy `Api::BlogPostsController#next`: the next OLDER
// published post in the request locale (id strictly below the given one). The
// oldest post has no next one — 404, which ends the page's infinite scroll.
func (s *Server) GetNextBlogPost(ctx context.Context, params api.GetNextBlogPostParams) (api.GetNextBlogPostRes, error) {
	row, err := s.blogPublic.NextOlder(ctx, s.i18n.Locale(ctx), int(params.ID))
	if err != nil {
		return nil, err
	}
	return s.blogPostToAPI(ctx, row)
}

// LikeBlogPost records the signed-in user's like and returns the post.
//
// Legacy counted a like once per browser session; the Go stack has no server
// session, so the rule is one like per user and post. Like legacy, it is a
// check before the insert, not a constraint: a unique index would make legacy's
// own repeat likes fail during the rollback window (ADR-0015: additive
// migrations only), and the table may already hold duplicate pairs. A repeat
// like is ignored rather than failing — the answer is the post, unchanged. Two
// truly concurrent clicks can still both insert, as they could on legacy.
func (s *Server) LikeBlogPost(ctx context.Context, params api.LikeBlogPostParams) (api.LikeBlogPostRes, error) {
	liker, ok := AuthenticatedUser(ctx)
	if !ok {
		return nil, errUnauthenticated
	}

	row, err := s.db.BlogPost.Query().
		Where(blogpost.IDEQ(int(params.ID))).
		WithCreator().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	liked, err := s.db.BlogPostLike.Query().
		Where(blogpostlike.BlogPostID(row.ID), blogpostlike.UserID(liker.ID)).
		Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !liked {
		if err := s.db.BlogPostLike.Create().
			SetBlogPostID(row.ID).
			SetUserID(liker.ID).
			Exec(ctx); err != nil {
			return nil, err
		}
	}

	return s.blogPostToAPI(ctx, row)
}

// blogPostToAPI is the single-post read model shared by next and like.
func (s *Server) blogPostToAPI(ctx context.Context, row *ent.BlogPost) (*api.BlogPost, error) {
	items, err := s.blogPostsToAPI(ctx, []*ent.BlogPost{row})
	if err != nil {
		return nil, err
	}
	return &items[0], nil
}
