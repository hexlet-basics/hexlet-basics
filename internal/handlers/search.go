package handlers

import (
	"context"

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
