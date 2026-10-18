package handlers_test

import (
	"context"
	"hexletbasics/internal/jobs"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/courseversion"
	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

func TestAdminCreateCourseVersion(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	res, err := h.Client.AdminCreateCourseVersion(ctx, api.AdminCreateCourseVersionParams{ID: int32(courseRubyIDA())})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())

	version := res
	assert.Equal(t, "created", version.State.Value)

	// The row exists in `created` state, owned by the ruby course.
	row := h.DB.CourseVersion.Query().Where(courseversion.ID(int(version.ID))).OnlyX(ctx)
	require.NotNil(t, row.State)
	assert.Equal(t, "created", *row.State)
	assert.Equal(t, courseRubyIDA(), row.CourseID)

	// The loader job was enqueued for exactly this version.
	queued := testsupport.RequireEnqueued(t, h.Services, jobs.ExerciseLoaderArgs{})
	assert.Equal(t, int(version.ID), queued.VersionID)
}

func TestAdminCreateCourseVersionNotFound(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminCreateCourseVersion(ctx, api.AdminCreateCourseVersionParams{ID: 999999999})
	requireProblem(t, err, http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())

	// A missing course enqueues nothing.
	testsupport.RequireNotEnqueued(t, h.Services, jobs.ExerciseLoaderArgs{})
}
