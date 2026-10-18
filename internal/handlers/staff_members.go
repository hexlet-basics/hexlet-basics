package handlers

import (
	"context"

	"hexletbasics/internal/api"
)

// The admin operations are generated (ent/admincrud.go); the schema's AdminCRUD
// annotation carries the eager-loaded user and role (with its permissions).
// user_id is unique (one staff record per user), so a duplicate grant surfaces
// as 409 via the central handler.

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
	if err := s.selfProtection.StaffRecordChange(ctx, actorOfRequest(ctx), 0, int(req.UserId)); err != nil {
		return nil, err
	}

	result, err := s.db.StaffMember.AdminCreate(ctx, req, s.conv.ToStaffMember)
	if err != nil {
		return result, errs.Wrapf(err, "admin create staff member")
	}

	return result, nil
}

func (s *Server) AdminUpdateStaffMember(ctx context.Context, req *api.StaffMemberInput, params api.AdminUpdateStaffMemberParams) (*api.StaffMember, error) {
	if err := s.selfProtection.StaffRecordChange(ctx, actorOfRequest(ctx), int(params.ID), int(req.UserId)); err != nil {
		return nil, err
	}

	result, err := s.db.StaffMember.AdminUpdate(ctx, params.ID, req, s.conv.ToStaffMember)
	if err != nil {
		return result, errs.Wrapf(err, "admin update staff member")
	}

	return result, nil
}

func (s *Server) AdminDeleteStaffMember(ctx context.Context, params api.AdminDeleteStaffMemberParams) error {
	if err := s.selfProtection.StaffRecordChange(ctx, actorOfRequest(ctx), int(params.ID), 0); err != nil {
		return err
	}

	if err := s.db.StaffMember.AdminDelete(ctx, params.ID); err != nil {
		return errs.Wrapf(err, "admin delete staff member")
	}

	return nil
}
