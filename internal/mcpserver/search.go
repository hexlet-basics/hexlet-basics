package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/search"
)

// The two content tools of the catalogue, in the shapes ChatGPT's connectors and
// deep research require of an MCP server (#877): `search(query)` answers
// `{results: [{id, title, url}]}` and `fetch(id)` answers `{id, title, text, url,
// metadata}`, both as structured content and as the same JSON in a text block.
// They are named as that convention names them, not after an operationId, because
// a client built to that convention looks for these two names; what they read is
// still the HTTP operations' own (searchCatalogue, getCourse, getCourseLesson,
// getBlogPost), through the Catalogue.
const (
	searchTool = "search"
	fetchTool  = "fetch"
)

// maxQueryLength is the longest query, the contract's own limit.
const maxQueryLength = 200

// alwaysLoad is the `_meta` hint that tells a client which defers tools until it
// needs them to load these two up front: they are how an assistant reads the
// catalogue at all, so a search for them would only cost a round trip.
var alwaysLoad = map[string]any{"anthropic/alwaysLoad": true}

// searchArgs are the arguments of the search tool.
type searchArgs struct {
	Query   string `json:"query"`
	Type    string `json:"type"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
}

// fetchArgs are the arguments of the fetch tool.
type fetchArgs struct {
	ID string `json:"id"`
}

// searchHit is one result of the search tool: the three fields the convention
// requires, plus the kind and a snippet, which it allows.
type searchHit struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	URL     string `json:"url"`
	Type    string `json:"type"`
	Snippet string `json:"snippet"`
}

// addSearchTools registers search and fetch. Like every catalogue tool they need
// the catalog:read Scope and spend the catalogue rate budget, and they answer
// what a signed-out visitor may read.
func (s *Server) addSearchTools(server *mcp.Server) error {
	types := make([]any, len(search.Types))
	for i, t := range search.Types {
		types[i] = string(t)
	}
	properties := map[string]any{
		"query": map[string]any{
			"type": "string", "minLength": 1, "maxLength": maxQueryLength,
			"description": "What to look for, in the words a person would use, in the language of the connection's locale. " +
				`A phrase in quotes is kept together, "or" and a leading "-" work as in a web search.`,
		},
		"type": map[string]any{
			"type": "string", "enum": types,
			"description": "Only this kind of record: a course, a lesson or a blog post. All kinds when absent.",
		},
	}
	for name, schema := range s.pageProperties() {
		properties[name] = schema
	}

	hit := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id":      map[string]any{"type": "string", "description": "The id fetch takes."},
			"title":   map[string]any{"type": "string"},
			"url":     map[string]any{"type": "string", "description": "The record's page on the website, absolute."},
			"type":    map[string]any{"type": "string", "enum": types},
			"snippet": map[string]any{"type": "string", "description": "A few words of the record's text around the match."},
		},
		"required":             []string{"id", "title", "url"},
		"additionalProperties": false,
	}
	integer := map[string]any{"type": "integer"}
	searchOutput := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"results": map[string]any{"type": "array", "items": hit},
			"total":   integer,
			"page":    integer,
			"perPage": integer,
		},
		"required":             []string{"results", "total", "page", "perPage"},
		"additionalProperties": false,
	}
	fetchOutput := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id":    map[string]any{"type": "string"},
			"title": map[string]any{"type": "string"},
			"text":  map[string]any{"type": "string", "description": "The record's text, as plain text or Markdown."},
			"url":   map[string]any{"type": "string", "description": "The record's page on the website, absolute."},
			"metadata": map[string]any{
				"type": "object", "description": "What else is known about the record (its kind, slugs, order).",
			},
		},
		"required":             []string{"id", "title", "text", "url", "metadata"},
		"additionalProperties": false,
	}

	s.addTool(server, toolSpec{
		name:  searchTool,
		title: "Search the catalogue",
		description: "Searches the Hexlet Basics catalogue (courses, lessons and blog posts) by what they say and returns the best matches, best first, " +
			"each with an id, a title, a snippet and the page's URL. Only text the site's authors wrote is searched, and only what a visitor can read. " +
			"Pass an id to fetch to read the whole record. Results are in the language of the connection's locale (Accept-Language).",
		input:  object(properties, "query"),
		output: searchOutput,
		scope:  oauthserver.ScopeCatalogRead,
		class:  ClassCatalogue,
		meta:   alwaysLoad,
	}, s.search)
	s.addTool(server, toolSpec{
		name:  fetchTool,
		title: "Fetch a catalogue record",
		description: "Reads one record of the catalogue by the id search returned: a course (`course:SLUG`), a lesson (`lesson:COURSE/SLUG`) or a blog post (`post:SLUG`). " +
			"Returns its title, its text and its URL, and metadata. It is what getCourse, getCourseLesson and getBlogPost answer a visitor, as one document; " +
			"use those tools for the full structure.",
		input:  object(map[string]any{"id": map[string]any{"type": "string", "minLength": 1, "description": "An id as search returns it, for example course:javascript-ru, lesson:javascript/hello-world or post:hello-world."}}, "id"),
		output: fetchOutput,
		scope:  oauthserver.ScopeCatalogRead,
		class:  ClassCatalogue,
		meta:   alwaysLoad,
	}, s.fetch)
	return nil
}

// absolute is the website URL of a path.
func (s *Server) absolute(path string) string {
	return strings.TrimRight(s.deps.SiteURL, "/") + path
}

// localePrefix is the prefix of the website's URLs in the locale of the call.
func (s *Server) localePrefix(ctx context.Context) string {
	if s.deps.Locale == nil {
		return localization.PathPrefix("")
	}
	return localization.PathPrefix(s.deps.Locale(ctx))
}

func (s *Server) search(ctx context.Context, c call) (result, error) {
	var args searchArgs
	if err := c.decode(&args); err != nil {
		return result{}, err
	}
	query := strings.TrimSpace(args.Query)
	switch {
	case query == "":
		return result{}, userError{errors.New("query must not be empty")}
	case len(query) > maxQueryLength:
		return result{}, userError{fmt.Errorf("query is longer than %d characters", maxQueryLength)}
	}
	page, perPage, _, err := s.resolve(pageArgs{Page: args.Page, PerPage: args.PerPage})
	if err != nil {
		return result{}, err
	}
	params := api.SearchCatalogueParams{Q: query, Page: api.NewOptInt32(int32(page)), PerPage: api.NewOptInt32(int32(perPage))}
	if args.Type != "" {
		if _, ok := searchType(args.Type); !ok {
			return result{}, userError{fmt.Errorf("type must be one of course, lesson, post, not %q", args.Type)}
		}
		params.Type = api.NewOptSearchResultType(api.SearchResultType(args.Type))
	}

	found, err := s.deps.Catalogue.SearchCatalogue(ctx, params)
	if err != nil {
		return result{}, fmt.Errorf("search the catalogue: %w", err)
	}
	hits := make([]searchHit, len(found.Items))
	for i, item := range found.Items {
		hits[i] = searchHit{ID: item.ID, Title: item.Title, URL: s.absolute(item.Path), Type: string(item.Type), Snippet: item.Snippet}
	}
	return result{
		data:       map[string]any{"results": hits, "total": found.Total, "page": found.Page, "perPage": found.PerPage},
		list:       &listInfo{key: "results", page: int(found.Page), perPage: int(found.PerPage)},
		textIsData: true,
	}, nil
}

// searchType reads the argument as a kind of record.
func searchType(name string) (search.Type, bool) {
	for _, t := range search.Types {
		if string(t) == name {
			return t, true
		}
	}
	return "", false
}

func (s *Server) fetch(ctx context.Context, c call) (result, error) {
	var args fetchArgs
	if err := c.decode(&args); err != nil {
		return result{}, err
	}
	ref, err := search.ParseID(args.ID)
	if err != nil {
		return result{}, userError{fmt.Errorf("%q is not an id search returns: they look like course:SLUG, lesson:COURSE/SLUG and post:SLUG", args.ID)}
	}

	var doc *document
	switch ref.Type {
	case search.Course:
		doc, err = s.fetchCourse(ctx, ref)
	case search.Lesson:
		doc, err = s.fetchLesson(ctx, ref)
	case search.Post:
		doc, err = s.fetchPost(ctx, ref)
	}
	if err != nil {
		return result{}, err
	}
	if doc == nil {
		// The id is well formed and names nothing a visitor can read: the same
		// answer whether the record is missing, a draft or in another locale.
		return result{}, userError{fmt.Errorf("no record has the id %q (it may be in another locale: ids of courses and lessons depend on the connection's Accept-Language)", args.ID)}
	}
	return result{
		data: map[string]any{
			"id": doc.id, "title": doc.title, "text": doc.text,
			"url": s.absolute(s.localePrefix(ctx) + doc.path), "metadata": doc.metadata,
		},
		textIsData: true,
	}, nil
}

// fetchCourse reads the course page the id names. A landing page that was
// archived answers with where it moved to, which is told rather than followed:
// the id the assistant holds is not the one that is current.
//
// The website serves a landing page by its URL whatever its state, and a lesson
// of a course that is not completed yet, because it never links to them; an id is
// a link an assistant holds, and what a search can never return a fetch must not
// start to hand out, so a landing page that is not published and listed, or a
// course that is not completed, is not found here, as it is for search.
func (s *Server) fetchCourse(ctx context.Context, ref search.Ref) (*document, error) {
	res, err := s.deps.Catalogue.GetCourse(ctx, api.GetCourseParams{Slug: ref.Slug})
	if missing(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	view, ok := res.(*api.CourseView)
	if !ok {
		return nil, notFoundOrUnexpected(res, "course")
	}
	if !view.MovedTo.Null {
		return nil, userError{fmt.Errorf("%s has moved: read %s instead", search.CourseID(ref.Slug), search.CourseID(view.MovedTo.Value))}
	}
	if view.LandingPage.Null || view.LandingPage.Value.State != api.NewNilLandingPageState(api.LandingPageStatePublished) ||
		view.LandingPage.Value.Listed != api.NewNilBool(true) ||
		view.Course.Readiness != api.NewNilCourseReadiness(api.CourseReadinessCompleted) {
		return nil, nil
	}
	return courseDocument(ref, view), nil
}

func (s *Server) fetchLesson(ctx context.Context, ref search.Ref) (*document, error) {
	res, err := s.deps.Catalogue.GetCourseLesson(ctx, api.GetCourseLessonParams{CourseSlug: ref.CourseSlug, Slug: ref.Slug})
	if missing(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	view, ok := res.(*api.CourseLessonView)
	if !ok {
		return nil, notFoundOrUnexpected(res, "lesson")
	}
	if view.Lesson.Course.Readiness != api.NewNilCourseReadiness(api.CourseReadinessCompleted) {
		return nil, nil
	}
	return lessonDocument(ref, view), nil
}

func (s *Server) fetchPost(ctx context.Context, ref search.Ref) (*document, error) {
	res, err := s.deps.Catalogue.GetBlogPost(ctx, api.GetBlogPostParams{Slug: ref.Slug})
	if missing(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	view, ok := res.(*api.BlogPostView)
	if !ok {
		return nil, notFoundOrUnexpected(res, "blog post")
	}
	return postDocument(ref, view), nil
}

// missing is a read that found nothing a visitor may open: the handlers answer
// that with the contract's not-found body or, for the blog, a store not-found.
func missing(err error) bool {
	return ent.IsNotFound(err) || errors.Is(err, domain.ErrNotFound)
}

// notFoundOrUnexpected is a nil document for a not-found body, which fetch
// answers in one way for every id that names nothing, and an error for anything
// else.
func notFoundOrUnexpected(res any, what string) error {
	if _, ok := res.(*api.NotFoundError); ok {
		return nil
	}
	return fmt.Errorf("unexpected %s result %T", what, res)
}
