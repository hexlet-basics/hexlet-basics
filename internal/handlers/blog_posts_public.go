package handlers

import (
	"context"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/blogpostlike"
	"hexletbasics/ent/landingpage"
	"hexletbasics/ent/predicate"
	"hexletbasics/internal/api"
)

// The public blog (legacy `Web::BlogPostsController`, `Web::BlogPosts::
// LikesController`, `Api::BlogPostsController#next`). Every read is scoped the
// way legacy scoped it — published posts in the request locale, newest (highest
// id) first — so a draft, an archived post or a post in another language is
// simply not found. The like is the exception: legacy looked the post up by slug
// alone, so any existing post can be liked.

// recommendedPostsLimit is how many other posts the post page offers (legacy
// `.limit(2)`).
const recommendedPostsLimit = 2

// publishedInLocale is legacy `BlogPost.published_state.with_locale`.
func (s *Server) publishedInLocale(ctx context.Context) predicate.BlogPost {
	return blogpost.And(
		blogpost.StateEQ(string(api.BlogPostStatePublished)),
		blogpost.LocaleEQ(s.i18n.Locale(ctx)),
	)
}

// ListBlogPosts returns a page of published posts in the request locale,
// newest first (legacy `order(id: :desc)`, pagy's 20 per page).
func (s *Server) ListBlogPosts(ctx context.Context, params api.ListBlogPostsParams) (*api.BlogPostPage, error) {
	return listPageWith(ctx, params.Page, params.PerPage,
		func() *ent.BlogPostQuery {
			return s.db.BlogPost.Query().
				Where(s.publishedInLocale(ctx)).
				WithCreator().
				Order(ent.Desc(blogpost.FieldID))
		},
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
	row, err := s.db.BlogPost.Query().
		Where(s.publishedInLocale(ctx), blogpost.SlugEQ(params.Slug)).
		WithCreator().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Divergence, on purpose: legacy wrote `.except(blog_post)`, which is
	// ActiveRecord's query-part remover, not a record exclusion — a no-op, so its
	// two posts were unordered and could include the post being read. Excluding
	// the current post and ordering newest first is what that line meant.
	recommendedRows, err := s.db.BlogPost.Query().
		Where(s.publishedInLocale(ctx), blogpost.IDNEQ(row.ID)).
		WithCreator().
		Order(ent.Desc(blogpost.FieldID)).
		Limit(recommendedPostsLimit).
		All(ctx)
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
// one is skipped. The course ids arrive already in the post's display order
// (blogRelatedCourseIDs orders them in SQL), so the pages are put back in that
// order here rather than joining the items table again for a handful of rows.
// A course with several main pages in one locale (legacy `has_one` picked one
// arbitrarily) keeps its oldest.
func (s *Server) relatedLandingPages(ctx context.Context, courseIDs []int32) ([]api.CourseCatalogItem, error) {
	if len(courseIDs) == 0 {
		return []api.CourseCatalogItem{}, nil
	}

	ids := lo.Map(courseIDs, func(id int32, _ int) int { return int(id) })

	pages, err := s.db.LandingPage.Query().
		Where(
			landingpage.CourseIDIn(ids...),
			landingpage.MainEQ(true),
			landingpage.LocaleEQ(s.i18n.Locale(ctx)),
		).
		WithCourse(func(q *ent.CourseQuery) { q.WithCurrentVersion() }).
		Order(landingpage.ByID()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	byCourse := make(map[int]*ent.LandingPage, len(pages))
	for _, p := range pages {
		if _, seen := byCourse[p.CourseID]; !seen {
			byCourse[p.CourseID] = p
		}
	}

	ordered := make([]*ent.LandingPage, 0, len(byCourse))
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
	row, err := s.db.BlogPost.Query().
		Where(s.publishedInLocale(ctx), blogpost.IDLT(int(params.ID))).
		WithCreator().
		Order(ent.Desc(blogpost.FieldID)).
		First(ctx)
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
