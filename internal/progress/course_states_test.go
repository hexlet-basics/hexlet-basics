package progress_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/fixture"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

func TestCourseStatesOfNoCoursesIsEmpty(t *testing.T) {
	db := testsupport.NewClient(t)

	states, err := progress.New(db, nil, nil, nil).CourseStates(t.Context(), fixture.ID(fixture.UserAdminLearner), nil)
	require.NoError(t, err)
	assert.Empty(t, states)
}
