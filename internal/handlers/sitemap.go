package handlers

import (
	"context"

	"hexletbasics/internal/api"
)

// GetSitemap is the data behind the HTML sitemap page (legacy
// `HomeController#sitemap`); the sitemap module reads the rows.
func (s *Server) GetSitemap(ctx context.Context) (*api.Sitemap, error) {
	page, err := s.sitemap.Read(ctx)
	if err != nil {
		return nil, err
	}

	return &api.Sitemap{
		LandingPages: s.conv.ToSitemapCourseLandingPages(page.LandingPages),
		BlogPosts:    s.conv.ToSitemapBlogPosts(page.BlogPosts),
		Categories:   s.conv.ToCourseCategories(page.Categories),
	}, nil
}
