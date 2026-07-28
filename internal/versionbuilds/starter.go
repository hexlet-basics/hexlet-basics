// Package versionbuilds starts durable course-version builds.
package versionbuilds

import (
	"context"
	"database/sql"

	"github.com/riverqueue/river"
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
	db    *sql.DB
	river *river.Client[*sql.Tx]
}

// NewStarter wires the shared pool and River client.
func NewStarter(db *sql.DB, riverClient *river.Client[*sql.Tx]) *Starter {
	return &Starter{db: db, river: riverClient}
}

// Start creates a version in `created` state and schedules its loader job.
func (s *Starter) Start(ctx context.Context, courseID int) (_ *ent.CourseVersion, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, oops.Wrapf(err, "begin course version tx")
	}
	// Rollback is harmless after a successful commit and also releases the
	// connection if a future edit introduces an early return or panic.
	defer func() { _ = tx.Rollback() }()

	txClient := store.NewTxClient(tx)
	version, err := txClient.CourseVersion.Create().
		SetLanguageID(courseID).
		SetState(stateCreated).
		Save(ctx)
	if err != nil {
		return nil, oops.Wrapf(err, "create course version")
	}

	if _, err = s.river.InsertTx(ctx, tx, jobs.ExerciseLoaderArgs{VersionID: version.ID}, nil); err != nil {
		return nil, oops.Wrapf(err, "enqueue exercise loader")
	}
	if err = tx.Commit(); err != nil {
		return nil, oops.Wrapf(err, "commit course version tx")
	}
	return version, nil
}
