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
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	res, err := h.Client.AdminCreateCourseVersion(ctx, api.AdminCreateCourseVersionParams{ID: int32(courseRubyIDA)})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())

	version, ok := res.(*api.CourseVersion)
	require.True(t, ok, "expected a CourseVersion body, got %T", res)
	assert.Equal(t, "created", version.State.Value)

	// The row exists in `created` state, owned by the ruby course.
	row := h.DB.CourseVersion.Query().Where(courseversion.ID(int(version.ID))).OnlyX(ctx)
	require.NotNil(t, row.State)
	assert.Equal(t, "created", *row.State)
	assert.Equal(t, courseRubyIDA, row.CourseID)

	// The loader job was enqueued for exactly this version.
	queued := testsupport.RequireEnqueued(t, h.Services, jobs.ExerciseLoaderArgs{})
	assert.Equal(t, int(version.ID), queued.VersionID)
}

func TestAdminCreateCourseVersionNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	res, err := h.Client.AdminCreateCourseVersion(ctx, api.AdminCreateCourseVersionParams{ID: 999999999})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
	_, ok := res.(*api.NotFoundError)
	assert.True(t, ok, "expected NotFoundError, got %T", res)

	// A missing course enqueues nothing.
	testsupport.RequireNotEnqueued(t, h.Services, jobs.ExerciseLoaderArgs{})
}
