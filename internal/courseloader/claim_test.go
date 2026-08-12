package courseloader

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

func TestClaimAllowsOnlyOneConcurrentOwner(t *testing.T) {
	ctx := context.Background()
	sqlDB, err := sql.Open("pgx", testdb.DatabaseURL())
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	db := store.NewClient(sqlDB)
	ruby, err := db.Course.Query().Where(course.Slug("ruby")).Only(ctx)
	require.NoError(t, err)
	version, err := db.CourseVersion.Create().
		SetCourseID(ruby.ID).
		SetState(stateCreated).
		Save(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = sqlDB.ExecContext(
			context.Background(),
			"DELETE FROM language_versions WHERE id = $1",
			version.ID,
		)
	})

	loader := NewLoader(db, nil, nil, nil, nil)
	const contenders = 8
	start := make(chan struct{})
	results := make(chan bool, contenders)
	errors := make(chan error, contenders)
	var wg sync.WaitGroup
	for range contenders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			claimed, err := loader.claim(ctx, version.ID)
			results <- claimed
			errors <- err
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	close(errors)

	for err := range errors {
		require.NoError(t, err)
	}
	claimedCount := 0
	for claimed := range results {
		if claimed {
			claimedCount++
		}
	}
	assert.Equal(t, 1, claimedCount)

	saved, err := db.CourseVersion.Get(ctx, version.ID)
	require.NoError(t, err)
	require.NotNil(t, saved.State)
	assert.Equal(t, stateBuilding, *saved.State)
}
