// Package versionbuilds starts durable course-version builds.
package versionbuilds

import (
	"context"
	"database/sql"

	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/store"
)

const stateCreated = "created"

// Starter atomically creates a course version and its River job. Both rows are
// written through the same pgx transaction, so neither can survive without the
// other.
type Starter struct {
	store store.Transactor
	jobs  *jobs.Enqueuer
}

// NewStarter wires the shared pool and the job enqueuer.
func NewStarter(txStore store.Transactor, enqueuer *jobs.Enqueuer) *Starter {
	return &Starter{store: txStore, jobs: enqueuer}
}

// Start creates a version in `created` state and schedules its loader job.
func (s *Starter) Start(ctx context.Context, courseID int) (_ *ent.CourseVersion, err error) {
	var version *ent.CourseVersion
	err = s.store.WithinTx(ctx, func(tx *sql.Tx, txClient *ent.Client) error {
		version, err = txClient.CourseVersion.Create().
			SetCourseID(courseID).
			SetState(stateCreated).
			Save(ctx)
		if err != nil {
			return oops.Wrapf(err, "create course version")
		}

		return s.jobs.EnqueueTx(ctx, tx, jobs.ExerciseLoaderArgs{VersionID: version.ID})
	})
	if err != nil {
		return nil, err
	}
	return version, nil
}
