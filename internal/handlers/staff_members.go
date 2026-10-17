package handlers

import (
	"context"

	"hexletbasics/internal/api"
)

// The admin operations are generated (ent/admincrud.go); the schema's AdminCRUD
// annotation carries the eager-loaded user and role (with its permissions).
// user_id is unique (one staff record per user), so a duplicate grant surfaces
// as 409 via the central handler.

func (s *Server) AdminListStaffMembers(ctx context.Context, params api.AdminListStaffMembersParams) (api.AdminListStaffMembersRes, error) {
	return s.db.StaffMember.AdminList(ctx, params.Page, params.PerPage, s.conv.ToStaffMembers)
}

func (s *Server) AdminGetStaffMember(ctx context.Context, params api.AdminGetStaffMemberParams) (api.AdminGetStaffMemberRes, error) {
	return s.db.StaffMember.AdminGet(ctx, params.ID, s.conv.ToStaffMember)
}

func (s *Server) AdminCreateStaffMember(ctx context.Context, req *api.StaffMemberInput) (api.AdminCreateStaffMemberRes, error) {
	if err := s.selfProtection.StaffRecordChange(ctx, actorOfRequest(ctx), 0, int(req.UserId)); err != nil {
		return nil, protection(err)
	}
	return s.db.StaffMember.AdminCreate(ctx, req, s.conv.ToStaffMember)
}

func (s *Server) AdminUpdateStaffMember(ctx context.Context, req *api.StaffMemberInput, params api.AdminUpdateStaffMemberParams) (api.AdminUpdateStaffMemberRes, error) {
	if err := s.selfProtection.StaffRecordChange(ctx, actorOfRequest(ctx), int(params.ID), int(req.UserId)); err != nil {
		return nil, protection(err)
	}
	return s.db.StaffMember.AdminUpdate(ctx, params.ID, req, s.conv.ToStaffMember)
}

func (s *Server) AdminDeleteStaffMember(ctx context.Context, params api.AdminDeleteStaffMemberParams) (api.AdminDeleteStaffMemberRes, error) {
	if err := s.selfProtection.StaffRecordChange(ctx, actorOfRequest(ctx), int(params.ID), 0); err != nil {
		return nil, protection(err)
	}
	if err := s.db.StaffMember.AdminDelete(ctx, params.ID); err != nil {
		return nil, err
	}
	return &api.AdminDeleteStaffMemberNoContent{}, nil
}
