package handlers

import (
	"context"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/inputconv"
	"hexletbasics/internal/leadcapture"
)

// CreateLead stores a contact request from the site's lead form (legacy
// `leads#create`) and raises LeadCreated for amoCRM. The first-visit
// attribution comes from the form and the IP from the request, since ahoy is
// not ported (ADR-0015); the recorder owns the write and the event together.
func (s *Server) CreateLead(ctx context.Context, req *api.LeadInput) (*api.Lead, error) {
	u := MustAuthenticatedUser(ctx)

	ip, hasIP := ClientIP(ctx)

	created, err := s.leads.Create(ctx, leadcapture.Submission{
		UserID:        u.ID,
		ContactMethod: string(req.ContactMethod),
		ContactValue:  req.ContactValue,
		YMClientID:    inputconv.Ptr(req.YmClientId),
		FirstVisit:    firstVisitOf(req.FirstVisit),
		IP:            lo.Ternary(hasIP, &ip, nil),
	})
	if err != nil {
		return nil, err
	}

	out := s.conv.ToLead(created)

	return &out, nil
}

// firstVisitOf maps the form's attribution. A browser that kept no cookie
// sends none, and then there is no attribution at all — not empty strings.
func firstVisitOf(input api.OptNilFirstVisit) leadcapture.FirstVisit {
	visit, ok := input.Get()
	if !ok {
		return leadcapture.FirstVisit{}
	}

	return leadcapture.FirstVisit{
		UTMSource:   inputconv.Ptr(visit.UtmSource),
		UTMMedium:   inputconv.Ptr(visit.UtmMedium),
		UTMCampaign: inputconv.Ptr(visit.UtmCampaign),
		UTMContent:  inputconv.Ptr(visit.UtmContent),
		UTMTerm:     inputconv.Ptr(visit.UtmTerm),
		LandingPage: inputconv.Ptr(visit.LandingPage),
		Referrer:    inputconv.Ptr(visit.Referrer),
	}
}

// AdminListLeads returns a page of sales leads, newest first. The admin surface
// is read-only (legacy `admin/leads#index`), so there is no get/update path.
func (s *Server) AdminListLeads(ctx context.Context, params api.AdminListLeadsParams) (*api.LeadPage, error) {
	return listFrom(ctx, params.Page, params.PerPage,
		func(ctx context.Context, page domain.Page) ([]*ent.Lead, int, error) {
			return s.leadDirectory.List(ctx, params.Q.Value, page)
		},
		s.conv.ToLeads,
		func(items []api.Lead, total, page, perPage int32) *api.LeadPage {
			return &api.LeadPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}
