// Package versionbuilds starts durable course-version builds.
package versionbuilds

import (
	"context"
	"database/sql"

	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/internal/domain"
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

// Start creates a version in `created` state and schedules its loader job. A
// missing course is domain.ErrNotFound, before anything is written.
func (s *Starter) Start(ctx context.Context, courseID int) (*ent.CourseVersion, error) {
	return s.start(ctx, func(c *ent.Client) (int, error) {
		exists, err := c.Course.Query().Where(course.ID(courseID)).Exist(ctx)
		if err == nil && !exists {
			return 0, domain.ErrNotFound
		}
		return courseID, domain.FromStore(err)
	})
}

// StartBySlug is Start for the course with the given slug (what a delivery from
// the course's repository names). No such course is domain.ErrNotFound.
func (s *Starter) StartBySlug(ctx context.Context, slug string) (*ent.CourseVersion, error) {
	return s.start(ctx, func(c *ent.Client) (int, error) {
		id, err := c.Course.Query().Where(course.Slug(slug)).OnlyID(ctx)
		return id, domain.FromStore(err)
	})
}

// start resolves the course and creates its version and its loader job in one
// transaction, so neither survives without the other.
func (s *Starter) start(ctx context.Context, resolve func(*ent.Client) (int, error)) (*ent.CourseVersion, error) {
	var version *ent.CourseVersion
	err := s.store.WithinTx(ctx, func(tx *sql.Tx, txClient *ent.Client) error {
		courseID, err := resolve(txClient)
		if err != nil {
			return err
		}

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
