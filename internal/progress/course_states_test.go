package progress_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

// Asking for many Courses at once must answer exactly what asking for each one
// does: the batch is a cheaper route to the same read model, not a second
// opinion about availability or completion. A Course the learner never joined
// still gets its shape, with nothing finished.
func TestCourseStatesMatchesCourseStatePerCourse(t *testing.T) {
	sqlDB, err := store.NewDB(testdb.DatabaseURL())
	require.NoError(t, err)
	db := store.NewClient(sqlDB)
	t.Cleanup(func() { _ = db.Close() })
	tracker := progress.New(db, nil, nil, nil)

	courseIDs := db.Course.Query().Where(course.CurrentVersionIDNotNil()).IDsX(t.Context())
	require.Greater(t, len(courseIDs), 1, "fixtures hold several published courses")

	const learnerID = 1001
	states, err := tracker.CourseStates(t.Context(), learnerID, courseIDs)
	require.NoError(t, err)
	require.Len(t, states, len(courseIDs))

	for _, courseID := range courseIDs {
		single, err := tracker.CourseState(t.Context(), progress.Learner{UserID: learnerID}, courseID)
		require.NoError(t, err)
		assert.Equal(t, single, states[courseID], "course %d", courseID)
	}
}

func TestCourseStatesOfNoCoursesIsEmpty(t *testing.T) {
	sqlDB, err := store.NewDB(testdb.DatabaseURL())
	require.NoError(t, err)
	db := store.NewClient(sqlDB)
	t.Cleanup(func() { _ = db.Close() })

	states, err := progress.New(db, nil, nil, nil).CourseStates(t.Context(), 1001, nil)
	require.NoError(t, err)
	assert.Empty(t, states)
}
