package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/internal/api"
)

// Admin course endpoints (legacy `/admin/languages`). The API embeds the
// current version, so reads load it. The cover image (coverAttachmentId) and
// the derived repositoryUrl are not written here — the cover is deferred until
// the Attachments uploader lands, and repositoryUrl is computed from the slug on
// read. There is no delete (parity with legacy). The review (AI re-review) and
// createVersion (exercise build) actions are left to the job subsystem and stay
// unimplemented for now.

func withCourseVersion(q *ent.CourseQuery) *ent.CourseQuery {
	return q.WithCurrentVersion()
}

func (s *Server) AdminListCourses(ctx context.Context, params api.AdminListCoursesParams) (*api.CoursePage, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.CourseQuery { return withCourseVersion(s.db.Course.Query()).Order(ent.Desc(course.FieldID)) },
		s.conv.ToCourses,
		func(items []api.Course, total, page, perPage int32) *api.CoursePage {
			return &api.CoursePage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

func (s *Server) AdminGetCourse(ctx context.Context, params api.AdminGetCourseParams) (*api.Course, error) {
	return getOne(ctx, int(params.ID),
		func(ctx context.Context, id int) (*ent.Course, error) {
			return withCourseVersion(s.db.Course.Query().Where(course.ID(id))).Only(ctx)
		},
		s.conv.ToCourse,
	)
}

func (s *Server) AdminCreateCourse(ctx context.Context, req *api.CourseInput) (*api.Course, error) {
	row, err := s.db.Course.Create().
		SetNillableSlug(nilStringPtr(req.Slug)).
		SetNillableLearnAs(nilLearnAsPtr(req.LearnAs)).
		SetNillableProgress(nilProgressPtr(req.Progress)).
		SetNillableHexletProgramLandingPage(nilStringPtr(req.HexletProgramLandingPage)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.AdminGetCourse(ctx, api.AdminGetCourseParams{ID: int32(row.ID)})
}

func (s *Server) AdminUpdateCourse(ctx context.Context, req *api.CourseInput, params api.AdminUpdateCourseParams) (*api.Course, error) {
	upd := s.db.Course.UpdateOneID(int(params.ID))
	applyNil(req.Slug.Null, req.Slug.Value, upd.SetSlug, upd.ClearSlug)
	applyNil(req.LearnAs.Null, string(req.LearnAs.Value), upd.SetLearnAs, upd.ClearLearnAs)
	applyNil(req.Progress.Null, string(req.Progress.Value), upd.SetProgress, upd.ClearProgress)
	applyNil(req.HexletProgramLandingPage.Null, req.HexletProgramLandingPage.Value, upd.SetHexletProgramLandingPage, upd.ClearHexletProgramLandingPage)

	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.AdminGetCourse(ctx, api.AdminGetCourseParams{ID: int32(row.ID)})
}

// nilLearnAsPtr / nilProgressPtr resolve ogen's nullable enum wrappers to a
// *string for ent's SetNillable* on the plain string columns.
func nilLearnAsPtr(v api.NilCourseLearnAs) *string {
	if v.Null {
		return nil
	}
	s := string(v.Value)
	return &s
}

func nilProgressPtr(v api.NilCourseProgress) *string {
	if v.Null {
		return nil
	}
	s := string(v.Value)
	return &s
}
