package handlers

import (
	"context"

	"hexletbasics/internal/api"
	"hexletbasics/internal/permission"
)

// The reads are generated (ent/admincrud.go, AdminCRUD ReadOnly, with the
// eager-loaded user and role); the writes go through permission.StaffEditor,
// which holds the self-protection rules. user_id is unique (one staff record
// per user), so a duplicate grant surfaces as 409.

func (s *Server) AdminListStaffMembers(ctx context.Context, params api.AdminListStaffMembersParams) (*api.StaffMemberPage, error) {
	result, err := s.db.StaffMember.AdminList(ctx, params.Page, params.PerPage, s.conv.ToStaffMembers)
	if err != nil {
		return result, errs.Wrapf(err, "admin list staff member")
	}

	return result, nil
}

func (s *Server) AdminGetStaffMember(ctx context.Context, params api.AdminGetStaffMemberParams) (*api.StaffMember, error) {
	result, err := s.db.StaffMember.AdminGet(ctx, params.ID, s.conv.ToStaffMember)
	if err != nil {
		return result, errs.Wrapf(err, "admin get staff member")
	}

	return result, nil
}

func (s *Server) AdminCreateStaffMember(ctx context.Context, req *api.StaffMemberInput) (*api.StaffMember, error) {
	row, err := s.staffEditor.Create(ctx, actorOfRequest(ctx), staffMemberFieldsOf(req))
	if err != nil {
		return nil, err
	}

	return s.AdminGetStaffMember(ctx, api.AdminGetStaffMemberParams{ID: int32(row.ID)}) //nolint:gosec // G115: a row id fits the contract's int32
}

func (s *Server) AdminUpdateStaffMember(ctx context.Context, req *api.StaffMemberInput, params api.AdminUpdateStaffMemberParams) (*api.StaffMember, error) {
	_, err := s.staffEditor.Update(ctx, actorOfRequest(ctx), int(params.ID), staffMemberFieldsOf(req))
	if err != nil {
		return nil, err
	}

	return s.AdminGetStaffMember(ctx, api.AdminGetStaffMemberParams(params))
}

func (s *Server) AdminDeleteStaffMember(ctx context.Context, params api.AdminDeleteStaffMemberParams) error {
	return s.staffEditor.Delete(ctx, actorOfRequest(ctx), int(params.ID))
}

// staffMemberFieldsOf is the contract's StaffMemberInput as the editor's command.
func staffMemberFieldsOf(in *api.StaffMemberInput) permission.StaffMemberFields {
	return permission.StaffMemberFields{
		UserID:         int(in.UserId),
		RoleID:         int(in.RoleId),
		AllowedLocales: in.AllowedLocales,
	}
}
