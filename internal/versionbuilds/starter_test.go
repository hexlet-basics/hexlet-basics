package versionbuilds_test

import (
	"testing"

	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport"
	"hexletbasics/internal/versionbuilds"
)

// TestStarterCreatesTheVersionAndItsJobTogether exercises ent and River over
// their shared database/sql transaction: starting a build leaves a `created`
// version and exactly one build job for it.
func TestStarterCreatesTheVersionAndItsJobTogether(t *testing.T) {
	ctx := t.Context()
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	riverClient, err := jobs.NewInsertOnlyClient(nil, nil, nil, nil)
	require.NoError(t, err)
	starter := versionbuilds.NewStarter(transactor, jobs.NewEnqueuer(transactor, riverClient))

	version, err := starter.Start(ctx, fixtures.ID(fixtures.CourseRuby))
	require.NoError(t, err)

	saved := db.CourseVersion.GetX(ctx, version.ID)
	require.NotNil(t, saved.State)
	assert.Equal(t, "created", *saved.State)
	job := rivertest.RequireInsertedTx[*riverdatabasesql.Driver](ctx, t, tx, jobs.ExerciseLoaderArgs{}, nil)
	assert.Equal(t, version.ID, job.Args.VersionID)

	// The build is unique per version: asking for it again queues nothing new.
	duplicate, err := riverClient.InsertTx(ctx, tx, jobs.ExerciseLoaderArgs{VersionID: version.ID}, nil)
	require.NoError(t, err)
	assert.True(t, duplicate.UniqueSkippedAsDuplicate)
}
