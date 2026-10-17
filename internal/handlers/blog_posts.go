package handlers

import (
	"context"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/blogposts"
	"hexletbasics/internal/inputconv"
)

// Blog posts (legacy `/admin/blog_posts`): full CRUD plus the related-courses
// actions (a hand-picked set and the AI suggestion job). `rich_body` is trusted editor HTML stored and returned exactly as
// given (no ActionText compatibility layer). The cover remains the single
// ActiveStorage blob served through `/api/storage/{key}` on read; the input's
// coverAttachmentId is deferred until blob covers land (same deferral as the
// course cover).

// AdminListBlogPosts returns a page of blog posts, newest first.
func (s *Server) AdminListBlogPosts(ctx context.Context, params api.AdminListBlogPostsParams) (api.AdminListBlogPostsRes, error) {
	return listFrom(ctx, params.Page, params.PerPage, s.blogReader.List, s.conv.ToBlogPosts,
		func(items []api.BlogPost, total, page, perPage int32) *api.BlogPostPage {
			return &api.BlogPostPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminGetBlogPost returns one blog post; a missing id surfaces as ent
// not-found, which the central APIErrorHandler maps to 404.
func (s *Server) AdminGetBlogPost(ctx context.Context, params api.AdminGetBlogPostParams) (api.AdminGetBlogPostRes, error) {
	return s.getAdminBlogPost(ctx, int(params.ID))
}

// getAdminBlogPost is the shared read model every write handler echoes back.
func (s *Server) getAdminBlogPost(ctx context.Context, id int) (*api.BlogPost, error) {
	post, err := s.blogReader.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	item := s.conv.ToBlogPost(post)
	return &item, nil
}

// AdminCreateBlogPost creates a post by the authenticated admin (legacy
// `blog_post.creator = current_user`), stamped with the request locale.
func (s *Server) AdminCreateBlogPost(ctx context.Context, req *api.BlogPostInput) (api.AdminCreateBlogPostRes, error) {
	row, err := s.postEditor.Create(ctx, MustAuthenticatedUser(ctx).ID, s.i18n.Locale(ctx), postFieldsOf(req))
	if err != nil {
		return nil, err
	}
	return s.getAdminBlogPost(ctx, row.ID)
}

// AdminUpdateBlogPost updates a post. A null nullable field clears its column;
// the stored locale is left untouched.
func (s *Server) AdminUpdateBlogPost(ctx context.Context, req *api.BlogPostInput, params api.AdminUpdateBlogPostParams) (api.AdminUpdateBlogPostRes, error) {
	row, err := s.postEditor.Update(ctx, int(params.ID), postFieldsOf(req))
	if err != nil {
		return nil, err
	}
	return s.getAdminBlogPost(ctx, row.ID)
}

// postFieldsOf is the contract's BlogPostInput as the editor's command.
func postFieldsOf(in *api.BlogPostInput) blogposts.PostFields {
	return blogposts.PostFields{
		Name:        inputconv.StringPtr(in.Name),
		Slug:        inputconv.StringPtr(in.Slug),
		Description: inputconv.StringPtr(in.Description),
		State:       inputconv.StringPtr(in.State),
		RichBody:    in.RichBody,
	}
}

// AdminDeleteBlogPost removes a post with its dependents (the likes and related
// courses go with it, in one transaction, inside the remover, which also
// reports a missing post as not-found).
func (s *Server) AdminDeleteBlogPost(ctx context.Context, params api.AdminDeleteBlogPostParams) (api.AdminDeleteBlogPostRes, error) {
	if err := s.postRemover.Remove(ctx, int(params.ID)); err != nil {
		return nil, err
	}
	return &api.AdminDeleteBlogPostNoContent{}, nil
}

// AdminSetBlogPostRelatedCourses replaces the post's promoted-courses set with
// the submitted ids, keeping their order as the display order. Legacy had no
// hand-picked set (only the AI job below), so this is an addition, sharing the
// job's replace step. An unknown course id fails the FK constraint, surfaced
// as 409 centrally.
func (s *Server) AdminSetBlogPostRelatedCourses(ctx context.Context, req *api.BlogPostRelatedCoursesInput, params api.AdminSetBlogPostRelatedCoursesParams) (api.AdminSetBlogPostRelatedCoursesRes, error) {
	id := int(params.ID)

	courseIDs := lo.Map(req.CourseIds, func(courseID int32, _ int) int { return int(courseID) })
	if err := s.relatedCoursesSet.Replace(ctx, id, courseIDs); err != nil {
		return nil, err
	}

	return s.getAdminBlogPost(ctx, id)
}

// AdminSuggestBlogPostRelatedCourses asks for the AI pick of the post's related
// courses (legacy `related_courses` member action →
// FindRelatedCoursesForBlogPostJob). The worker replaces the set
// asynchronously, so the answer carries no post; a missing id is a 404.
func (s *Server) AdminSuggestBlogPostRelatedCourses(ctx context.Context, params api.AdminSuggestBlogPostRelatedCoursesParams) (api.AdminSuggestBlogPostRelatedCoursesRes, error) {
	if err := s.relatedSuggester.Suggest(ctx, int(params.ID)); err != nil {
		return nil, err
	}
	return &api.AdminSuggestBlogPostRelatedCoursesNoContent{}, nil
}

// blogPostsToAPI is the read model of a set of posts: the blog reader assembles
// them (cover, likes, promoted courses, reading time) and the converter copies
// the result to the contract.
func (s *Server) blogPostsToAPI(ctx context.Context, rows []*ent.BlogPost) ([]api.BlogPost, error) {
	posts, err := s.blogReader.Assemble(ctx, rows)
	if err != nil {
		return nil, err
	}
	return s.conv.ToBlogPosts(posts), nil
}
