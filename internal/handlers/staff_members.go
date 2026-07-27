package handlers

import (
	"context"

	"github.com/lib/pq"

	"hexletbasics/ent"
	"hexletbasics/ent/staffmember"
	"hexletbasics/internal/api"
)

// withStaffMemberEdges loads the user and the role (with its permission count),
// both embedded in full by the StaffMember API model.
func withStaffMemberEdges(q *ent.StaffMemberQuery) *ent.StaffMemberQuery {
	return q.
		WithUser().
		WithRole(func(rq *ent.StaffRoleQuery) { rq.WithPermissions() })
}

func (s *Server) AdminListStaffMembers(ctx context.Context, params api.AdminListStaffMembersParams) (*api.StaffMemberPage, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.StaffMemberQuery {
			return withStaffMemberEdges(s.db.StaffMember.Query()).Order(ent.Desc(staffmember.FieldID))
		},
		s.conv.ToStaffMembers,
		func(items []api.StaffMember, total, page, perPage int32) *api.StaffMemberPage {
			return &api.StaffMemberPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

func (s *Server) AdminGetStaffMember(ctx context.Context, params api.AdminGetStaffMemberParams) (*api.StaffMember, error) {
	return getOne(ctx, int(params.ID),
		func(ctx context.Context, id int) (*ent.StaffMember, error) {
			return withStaffMemberEdges(s.db.StaffMember.Query().Where(staffmember.ID(id))).Only(ctx)
		},
		s.conv.ToStaffMember,
	)
}

// AdminCreateStaffMember grants a user a role. user_id is unique (one staff
// record per user), so a duplicate surfaces as 409 via the central handler.
func (s *Server) AdminCreateStaffMember(ctx context.Context, req *api.StaffMemberInput) (*api.StaffMember, error) {
	row, err := s.db.StaffMember.Create().
		SetUserID(int(req.UserId)).
		SetRoleID(int(req.RoleId)).
		SetAllowedLocales(pq.StringArray(req.AllowedLocales)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.AdminGetStaffMember(ctx, api.AdminGetStaffMemberParams{ID: int32(row.ID)})
}

func (s *Server) AdminUpdateStaffMember(ctx context.Context, req *api.StaffMemberInput, params api.AdminUpdateStaffMemberParams) (*api.StaffMember, error) {
	row, err := s.db.StaffMember.UpdateOneID(int(params.ID)).
		SetUserID(int(req.UserId)).
		SetRoleID(int(req.RoleId)).
		SetAllowedLocales(pq.StringArray(req.AllowedLocales)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.AdminGetStaffMember(ctx, api.AdminGetStaffMemberParams{ID: int32(row.ID)})
}

func (s *Server) AdminDeleteStaffMember(ctx context.Context, params api.AdminDeleteStaffMemberParams) error {
	return s.db.StaffMember.DeleteOneID(int(params.ID)).Exec(ctx)
}
