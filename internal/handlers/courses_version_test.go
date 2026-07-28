package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/courseversion"
	"hexletbasics/internal/api"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport"
)

func TestAdminCreateCourseVersion(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	res, err := h.Client.AdminCreateCourseVersion(ctx, api.AdminCreateCourseVersionParams{ID: courseRubyIDA})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())

	version, ok := res.(*api.CourseVersion)
	require.True(t, ok, "expected a CourseVersion body, got %T", res)
	assert.Equal(t, "created", version.State.Value)

	// The row exists in `created` state, owned by the ruby course.
	row := h.DB.CourseVersion.Query().Where(courseversion.ID(int(version.ID))).OnlyX(ctx)
	require.NotNil(t, row.State)
	assert.Equal(t, "created", *row.State)
	assert.Equal(t, courseRubyIDA, row.LanguageID)

	// The loader job was enqueued for exactly this version.
	require.Len(t, h.Enqueuer.Inserted, 1)
	args, ok := h.Enqueuer.Inserted[0].(jobs.ExerciseLoaderArgs)
	require.True(t, ok, "expected ExerciseLoaderArgs, got %T", h.Enqueuer.Inserted[0])
	assert.Equal(t, int(version.ID), args.VersionID)
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
	assert.Empty(t, h.Enqueuer.Inserted)
}
