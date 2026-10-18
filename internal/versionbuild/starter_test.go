package versionbuild_test

import (
	"testing"

	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport"
	"hexletbasics/internal/versionbuild"
)

// TestStarterCreatesTheVersionAndItsJobTogether exercises ent and River over
// their shared database/sql transaction: starting a build leaves a `created`
// version and exactly one build job for it.
func TestStarterCreatesTheVersionAndItsJobTogether(t *testing.T) {
	ctx := t.Context()
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	riverClient, err := jobs.NewInsertOnlyClient(nil, nil, nil, nil)
	require.NoError(t, err)

	starter := versionbuild.NewStarter(transactor, jobs.NewEnqueuer(transactor, riverClient))

	version, err := starter.Start(ctx, fixture.ID(fixture.CourseRuby))
	require.NoError(t, err)

	saved := db.CourseVersion.GetX(ctx, version.ID)
	require.NotNil(t, saved.State)
	assert.Equal(t, "created", *saved.State)
	job := rivertest.RequireInsertedTx[*riverdatabasesql.Driver](ctx, t, tx, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	}, nil)
	assert.Equal(t, version.ID, job.Args.VersionID)

	// The build is unique per version: asking for it again queues nothing new.
	duplicate, err := riverClient.InsertTx(ctx, tx, jobs.ExerciseLoaderArgs{VersionID: version.ID}, nil)
	require.NoError(t, err)
	assert.True(t, duplicate.UniqueSkippedAsDuplicate)
}

// A delivery from a course's repository names the course by slug: the build
// starts for that course, and an unknown slug starts nothing.
func TestStartBySlugBuildsTheCourseAndAnUnknownSlugNothing(t *testing.T) {
	ctx := t.Context()
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	services := testsupport.NewServices(t, tx, transactor)
	ruby := db.Course.GetX(ctx, fixture.ID(fixture.CourseRuby))
	require.NotNil(t, ruby.Slug)

	version, err := services.Starter.StartBySlug(ctx, *ruby.Slug)
	require.NoError(t, err)
	assert.Equal(t, ruby.ID, version.CourseID)
	assert.Equal(t, version.ID, testsupport.RequireEnqueued(t, services, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	}).VersionID)

	_, err = services.Starter.StartBySlug(ctx, "no-such-course")
	require.ErrorIs(t, err, domain.ErrNotFound)
	_, err = services.Starter.Start(ctx, 0)
	require.ErrorIs(t, err, domain.ErrNotFound)
}
