package handlers

import (
	"context"
	"errors"

	"hexletbasics/internal/api"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/search"
)

// SearchCatalogue is the full-text search of the public catalogue (#877): the
// courses, lessons and blog posts a visitor can read, in the request locale,
// best match first. Which records those are, and how they are ranked, is the
// search module's; the handler turns the request into its command and the hits
// into the contract's compact results, so the MCP `search` tool, which calls
// this method, answers what the website's API answers.
func (s *Server) SearchCatalogue(ctx context.Context, params api.SearchCatalogueParams) (*api.SearchResultPage, error) {
	var types []search.Type
	if params.Type.IsSet() {
		types = []search.Type{search.Type(params.Type.Value)}
	}
	locale := s.i18n.Locale(ctx)

	return listFrom(ctx, params.Page, params.PerPage,
		func(ctx context.Context, page domain.Page) ([]search.Hit, int, error) {
			return s.searcher.Search(ctx, search.SearchCatalogue{Text: params.Q, Locale: locale, Types: types, Page: page})
		},
		func(hits []search.Hit) []api.SearchResult {
			results := make([]api.SearchResult, len(hits))
			for i, hit := range hits {
				results[i] = api.SearchResult{
					Type: api.SearchResultType(hit.Type), ID: hit.ID, Title: hit.Title, Snippet: hit.Snippet, Path: hit.Path,
				}
			}
			return results
		},
		func(items []api.SearchResult, total, page, perPage int32) *api.SearchResultPage {
			return &api.SearchResultPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminSearch is the catalogue search of the people who edit it (and of the MCP
// tool of the same name): the same module and ranking as SearchCatalogue, with
// the storefront's filters off, so a draft, archived or unlisted landing page, a
// course not completed yet, a lesson of the newest build and an unpublished post
// are found, each with its state. Who may call it is the contract's (any staff
// member, in a locale their role allows), judged before this runs.
func (s *Server) AdminSearch(ctx context.Context, params api.AdminSearchParams) (*api.AdminSearchResultPage, error) {
	var types []search.Type
	if params.Type.IsSet() {
		types = []search.Type{search.Type(params.Type.Value)}
	}
	locale := s.i18n.Locale(ctx)

	page, err := listFrom(ctx, params.Page, params.PerPage,
		func(ctx context.Context, page domain.Page) ([]search.Hit, int, error) {
			return s.searcher.Search(ctx, search.SearchCatalogue{Text: params.Q, Locale: locale, Types: types, Page: page, Unrestricted: true})
		},
		func(hits []search.Hit) []api.AdminSearchResult {
			results := make([]api.AdminSearchResult, len(hits))
			for i, hit := range hits {
				results[i] = api.AdminSearchResult{
					Type: api.SearchResultType(hit.Type), ID: hit.ID, Title: hit.Title, Snippet: hit.Snippet, Path: hit.Path,
					State: hit.State, Readiness: hit.Readiness,
				}
			}
			return results
		},
		func(items []api.AdminSearchResult, total, page, perPage int32) *api.AdminSearchResultPage {
			return &api.AdminSearchResultPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
	if err != nil {
		return nil, err
	}
	return page, nil
}

// AdminFetch reads the record an AdminSearch result names, in the request locale
// and whatever state it is in. An id that is not one search returns, and one that
// names nothing, are the same not found.
func (s *Server) AdminFetch(ctx context.Context, params api.AdminFetchParams) (*api.AdminRecord, error) {
	ref, err := search.ParseID(params.ID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	rec, err := s.searcher.Fetch(ctx, search.FetchRecord{Ref: ref, Locale: s.i18n.Locale(ctx)})
	if errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	meta := api.AdminRecordMetadata{
		Type: api.SearchResultType(rec.Type), Slug: rec.Slug, State: rec.State, Readiness: rec.Readiness,
		CourseSlug: optString(rec.CourseSlug), BuildState: optString(rec.Build),
	}
	if rec.Listed != nil {
		meta.Listed = api.NewOptBool(*rec.Listed)
	}
	if rec.Position != nil {
		meta.Position = api.NewOptInt32(int32(*rec.Position))
	}
	return &api.AdminRecord{ID: rec.ID, Title: rec.Title, Text: rec.Text, Path: rec.Path, Metadata: meta}, nil
}

// optString is a contract optional string that is absent when empty.
func optString(value string) api.OptString {
	if value == "" {
		return api.OptString{}
	}
	return api.NewOptString(value)
}
