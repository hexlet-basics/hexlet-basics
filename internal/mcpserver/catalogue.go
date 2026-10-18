package mcpserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/oauthserver"
)

// The catalogue tools, named after the operationId they mirror. All of them
// read what a signed-out visitor reads on the website, through the handler code
// of the HTTP operation: the handler sees no session (a bearer token is not one),
// so what comes back is what it answers a visitor, and a field that exists only
// for a signed-in user is empty here.
const (
	listCoursesTool     = "listCourses"
	listCategoriesTool  = "listPublicCourseCategories"
	getCategoryTool     = "getPublicCourseCategory"
	getCourseTool       = "getCourse"
	getCourseLessonTool = "getCourseLesson"
	listBlogPostsTool   = "listBlogPosts"
	getBlogPostTool     = "getBlogPost"
)

// pageArgs are the arguments of every list tool.
type pageArgs struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"` //nolint:tagliatelle // an MCP tool argument is named as in the contract, camelCase
}

// pageProperties is the input schema of pageArgs. The maximum of perPage is not
// in the schema because the limit is configurable; the description states the
// default and the tool clamps a larger value (and says so in its text).
func (s *Server) pageProperties() map[string]any {
	return map[string]any{
		"page": map[string]any{"type": "integer", "minimum": 1, "description": "Page number, starting at 1. Default 1."},
		"perPage": map[string]any{
			"type": "integer", "minimum": 1,
			"description": fmt.Sprintf("Records per page. A value above %d is clamped to %d. Default %d.", s.limits.MaxPageSize, s.limits.MaxPageSize, s.limits.MaxPageSize),
		},
	}
}

// errPageArguments is the answer to a negative page or perPage.
var errPageArguments = errors.New("page and perPage must be at least 1")

// resolve checks the arguments and applies the page size limit: a perPage above
// the limit is clamped (not refused, an agent asking for "everything" still gets
// a useful page) and note says so for the text of the result.
func (s *Server) resolve(a pageArgs) (int, int, string, error) {
	if a.Page < 0 || a.PerPage < 0 {
		return 0, 0, "", userError{errPageArguments}
	}

	page := max(a.Page, 1)
	perPage := a.PerPage

	var note string

	if perPage == 0 {
		perPage = s.limits.MaxPageSize
	}

	if perPage > s.limits.MaxPageSize {
		note = fmt.Sprintf(" (perPage %d is above the limit, clamped to %d)", perPage, s.limits.MaxPageSize)
		perPage = s.limits.MaxPageSize
	}

	return page, perPage, note, nil
}

// slicePage is page of the list in pages of perPage; past the end it is empty.
func slicePage[T any](items []T, page, perPage int) []T {
	start := (page - 1) * perPage
	if start >= len(items) {
		return []T{}
	}

	return items[start:min(start+perPage, len(items))]
}

// addCatalogueTools registers the read-only catalogue tools. Each needs the
// catalog:read Scope and counts against the catalogue rate limit.
func (s *Server) addCatalogueTools(server *mcp.Server) error {
	specs := []struct {
		name string
		spec toolSpec
		run  func(ctx context.Context, c call) (result, error)
	}{
		{listCoursesTool, toolSpec{
			title:       "List courses",
			description: "Lists the published courses in catalogue order. An item's slug is the one getCourse takes.",
			input:       object(s.pageProperties()),
			name:        "",
			output:      nil,
			scope:       "",
			class:       0,
		}, s.listCourses},
		{listCategoriesTool, toolSpec{
			title:       "List course categories",
			description: "Lists the course categories (programming languages and topics) of the request locale. The slug of an item is the one getPublicCourseCategory takes.",
			input:       object(s.pageProperties()),
			name:        "",
			output:      nil,
			scope:       "",
			class:       0,
		}, s.listCategories},
		{getCategoryTool, toolSpec{
			title:       "Get a course category",
			description: "A course category by slug, with the courses it groups and its questions and answers.",
			input:       object(map[string]any{"slug": slugArg("The category's slug, as listPublicCourseCategories returns it.")}, "slug"),
			name:        "",
			output:      nil,
			scope:       "",
			class:       0,
		}, s.getCategory},
		{getCourseTool, toolSpec{
			title: "Get a course",
			description: "A course page: the course, its landing copy, lessons in order and modules. " +
				"The slug is the landing page's, as listCourses returns it (javascript-ru).",
			input:  object(map[string]any{"slug": slugArg("The landing page slug, as listCourses returns it.")}, "slug"),
			name:   "",
			output: nil,
			scope:  "",
			class:  0,
		}, s.getCourse},
		{getCourseLessonTool, toolSpec{
			title: "Get a lesson",
			description: "A lesson: instructions, theory, starter code, reference solution and tests. " +
				"courseSlug is the course's own (course.slug of getCourse), not a landing page slug.",
			input: object(map[string]any{
				"courseSlug": slugArg("The course's own slug (course.slug in getCourse)."),
				"slug":       slugArg("The lesson's slug, from the lessons of getCourse."),
			}, "courseSlug", "slug"),
			name:   "",
			output: nil,
			scope:  "",
			class:  0,
		}, s.getCourseLesson},
		{listBlogPostsTool, toolSpec{
			title:       "List blog posts",
			description: "Published blog posts of the request locale, newest first, in pages. The slug of an item is the one getBlogPost takes.",
			input:       object(s.pageProperties()),
			name:        "",
			output:      nil,
			scope:       "",
			class:       0,
		}, s.listBlogPosts},
		{getBlogPostTool, toolSpec{
			title:       "Get a blog post",
			description: "A published blog post by slug, with posts to read next and the courses it promotes.",
			input:       object(map[string]any{"slug": slugArg("The post's slug, as listBlogPosts returns it.")}, "slug"),
			name:        "",
			output:      nil,
			scope:       "",
			class:       0,
		}, s.getBlogPost},
	}

	for _, t := range specs {
		spec := t.spec
		spec.name = t.name
		spec.scope = oauthserver.ScopeCatalogRead
		spec.class = ClassCatalogue

		output, err := s.outputSchema(t.name)
		if err != nil {
			return errs.With("tool", t.name).Wrapf(err, "build output schema")
		}

		spec.output = output
		s.addTool(server, spec, t.run)
	}

	return nil
}

// outputSchema is the contract's schema of the operation the tool mirrors. The
// two operations that return a bare array are wrapped, because MCP structured
// content is an object: the array sits under one key, with the paging fields.
func (s *Server) outputSchema(operationID string) (map[string]any, error) {
	schema, err := responseSchema(operationID)
	if err != nil {
		return nil, err
	}

	key := map[string]string{listCoursesTool: "courses", listCategoriesTool: "categories"}[operationID]
	if key == "" {
		return schema, nil
	}

	integer := map[string]any{"type": "integer"}

	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			key:       schema,
			"total":   integer,
			"page":    integer,
			"perPage": integer,
		},
		"required":             []string{key, "total", "page", "perPage"},
		"additionalProperties": false,
	}, nil
}

func (s *Server) listCourses(ctx context.Context, c call) (result, error) {
	var args pageArgs
	if err := c.decode(&args); err != nil {
		return result{}, err
	}

	page, perPage, note, err := s.resolve(args)
	if err != nil {
		return result{}, err
	}

	all, err := s.deps.Catalogue.ListCourses(ctx)
	if err != nil {
		return result{}, errs.Wrapf(err, "list courses")
	}

	items := slicePage(all, page, perPage)

	slugs := make([]string, len(items))
	for i, course := range items {
		slugs[i] = course.Slug
	}

	return result{
		text: fmt.Sprintf("%s (page %d, %d in total)%s", summarize("courses", slugs), page, len(all), note),
		data: map[string]any{"courses": items, "total": len(all), "page": page, "perPage": perPage},
		list: &listInfo{key: "courses", page: page, perPage: perPage},
	}, nil
}

func (s *Server) listCategories(ctx context.Context, c call) (result, error) {
	var args pageArgs
	if err := c.decode(&args); err != nil {
		return result{}, err
	}

	page, perPage, note, err := s.resolve(args)
	if err != nil {
		return result{}, err
	}

	all, err := s.deps.Catalogue.ListPublicCourseCategories(ctx)
	if err != nil {
		return result{}, errs.Wrapf(err, "list course categories")
	}

	items := slicePage(all, page, perPage)

	slugs := make([]string, len(items))
	for i, category := range items {
		slugs[i] = category.Slug.Value
	}

	return result{
		text: fmt.Sprintf("%s (page %d, %d in total)%s", summarize("categories", slugs), page, len(all), note),
		data: map[string]any{"categories": items, "total": len(all), "page": page, "perPage": perPage},
		list: &listInfo{key: "categories", page: page, perPage: perPage},
	}, nil
}

type slugArgs struct {
	Slug string `json:"slug"`
}

func (s *Server) getCategory(ctx context.Context, c call) (result, error) {
	var args slugArgs
	if err := c.decode(&args); err != nil {
		return result{}, err
	}

	view, err := s.deps.Catalogue.GetPublicCourseCategory(ctx, api.GetPublicCourseCategoryParams{Slug: args.Slug})
	if err != nil {
		return result{}, notFoundOf(err, "category")
	}

	return result{
		text: fmt.Sprintf("Category %s with %d courses", args.Slug, len(view.LandingPages)),
		data: view,
		list: nil,
	}, nil
}

func (s *Server) getCourse(ctx context.Context, c call) (result, error) {
	var args slugArgs
	if err := c.decode(&args); err != nil {
		return result{}, err
	}

	view, err := s.deps.Catalogue.GetCourse(ctx, api.GetCourseParams{Slug: args.Slug})
	if err != nil {
		return result{}, notFoundOf(err, "course")
	}

	return result{
		text: fmt.Sprintf("Course %s with %d lessons", args.Slug, len(view.Lessons)),
		data: view,
		list: nil,
	}, nil
}

type lessonArgs struct {
	CourseSlug string `json:"courseSlug"` //nolint:tagliatelle // an MCP tool argument is named as in the contract, camelCase
	Slug       string `json:"slug"`
}

func (s *Server) getCourseLesson(ctx context.Context, c call) (result, error) {
	var args lessonArgs
	if err := c.decode(&args); err != nil {
		return result{}, err
	}

	view, err := s.deps.Catalogue.GetCourseLesson(ctx, api.GetCourseLessonParams{CourseSlug: args.CourseSlug, Slug: args.Slug})
	if err != nil {
		return result{}, notFoundOf(err, "lesson")
	}

	return result{
		text: fmt.Sprintf("Lesson %s of course %s", args.Slug, args.CourseSlug),
		data: view,
		list: nil,
	}, nil
}

func (s *Server) listBlogPosts(ctx context.Context, c call) (result, error) {
	var args pageArgs
	if err := c.decode(&args); err != nil {
		return result{}, err
	}

	page, perPage, note, err := s.resolve(args)
	if err != nil {
		return result{}, err
	}
	// The handler pages in SQL; the clamped page size goes to it, not the
	// agent's.
	params := api.ListBlogPostsParams{Page: api.NewOptInt32(apiconv.Int32FromInt(page)), PerPage: api.NewOptInt32(apiconv.Int32FromInt(perPage))}

	posts, err := s.deps.Catalogue.ListBlogPosts(ctx, params)
	if err != nil {
		return result{}, errs.Wrapf(err, "list blog posts")
	}

	slugs := make([]string, len(posts.Items))
	for i, post := range posts.Items {
		slugs[i] = post.Slug.Value
	}

	return result{
		text: fmt.Sprintf("%s (page %d, %d in total)%s", summarize("posts", slugs), posts.Page, posts.Total, note),
		data: posts,
		list: &listInfo{key: "items", page: int(posts.Page), perPage: int(posts.PerPage)},
	}, nil
}

func (s *Server) getBlogPost(ctx context.Context, c call) (result, error) {
	var args slugArgs
	if err := c.decode(&args); err != nil {
		return result{}, err
	}

	view, err := s.deps.Catalogue.GetBlogPost(ctx, api.GetBlogPostParams{Slug: args.Slug})
	if err != nil {
		return result{}, notFoundOf(err, "blog post")
	}

	return result{
		text: "Blog post " + args.Slug,
		data: view,
		list: nil,
	}, nil
}

// notFoundOf is the answer for a failed read of one record: a missing one is
// the tool's own not-found, anything else is the failure itself.
func notFoundOf(err error, what string) error {
	if errors.Is(domain.FromStore(err), domain.ErrNotFound) {
		return errNotFound(what)
	}

	return err
}
