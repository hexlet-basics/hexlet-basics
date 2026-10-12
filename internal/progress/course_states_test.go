package progress_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/progress"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

func TestCourseStatesOfNoCoursesIsEmpty(t *testing.T) {
	sqlDB, err := store.NewDB(testdb.DatabaseURL())
	require.NoError(t, err)
	db := store.NewClient(sqlDB)
	t.Cleanup(func() { _ = db.Close() })

	states, err := progress.New(db, nil, nil, nil).CourseStates(t.Context(), 1001, nil)
	require.NoError(t, err)
	assert.Empty(t, states)
}
