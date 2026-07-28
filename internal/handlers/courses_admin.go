package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/jobs"
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

// AdminCreateCourseVersion starts a build of a new course version: it creates the
// version row in `created` state and enqueues the exercise-loader job, which
// fetches the course repo, parses it, and (on success) promotes the new version
// live. Maps 1:1 to the legacy VersionsController#create → ExerciseLoaderJob. The
// 201 body is the freshly-created (not yet built) version, so the admin UI can
// poll its state.
func (s *Server) AdminCreateCourseVersion(ctx context.Context, params api.AdminCreateCourseVersionParams) (api.AdminCreateCourseVersionRes, error) {
	c, err := s.db.Course.Get(ctx, int(params.ID))
	if ent.IsNotFound(err) {
		return &api.NotFoundError{Message: "course not found"}, nil
	}
	if err != nil {
		return nil, err
	}

	version, err := s.db.CourseVersion.Create().
		SetLanguageID(c.ID).
		SetState(courseVersionStateCreated).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := s.enqueuer.Insert(ctx, jobs.ExerciseLoaderArgs{VersionID: version.ID}, nil); err != nil {
		return nil, err
	}

	body := apiconv.CourseVersionFromEnt(version)
	return &body, nil
}

// courseVersionStateCreated is the initial version state the loader requires
// before it will build (mirrors the legacy AASM initial state).
const courseVersionStateCreated = "created"

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
