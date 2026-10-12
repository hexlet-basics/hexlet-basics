package handlers_test

import (
	"context"
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
	assert.Equal(t, courseRubyIDA, row.CourseID)

	// One build was started, for the ruby course.
	calls := h.Starter.StartCalls()
	require.Len(t, calls, 1)
	assert.Equal(t, courseRubyIDA, calls[0].CourseID)
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
	assert.Empty(t, h.Starter.StartCalls())
}
