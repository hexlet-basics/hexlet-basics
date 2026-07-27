package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/lead"
	"hexletbasics/internal/api"
)

// AdminListLeads returns a page of sales leads, newest first. The admin surface
// is read-only (legacy `admin/leads#index`), so there is no get/create path.
func (s *Server) AdminListLeads(ctx context.Context, params api.AdminListLeadsParams) (*api.LeadPage, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.LeadQuery { return s.db.Lead.Query().Order(ent.Desc(lead.FieldID)) },
		s.conv.ToLeads,
		func(items []api.Lead, total, page, perPage int32) *api.LeadPage {
			return &api.LeadPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}
