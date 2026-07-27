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
	page := newPagination(params.Page, params.PerPage)

	total, err := s.db.Lead.Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Lead.Query().
		Order(ent.Desc(lead.FieldID)).
		Offset(page.Offset()).
		Limit(page.Limit()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return &api.LeadPage{
		Items:   s.conv.ToLeads(rows),
		Total:   int32(total),
		Page:    page.Page,
		PerPage: page.PerPage,
	}, nil
}
