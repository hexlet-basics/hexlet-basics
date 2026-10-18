// Package catalog reads the courses and course categories: the administrators'
// course picker, the visitor's category pages and the readiness of a course.
package catalog

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/transition"
)

// Readiness changes a Course's Readiness (GLOSSARY). A state is not an update:
// the admin CRUD never writes it, so a new Course is a draft by the column
// default and the value moves only through these operations.
//
// Any change between two distinct values is legal, repeating the current one is
// a *transition.IllegalError. The readiness column is all a move changes. What it
// controls is read from it by the queries, with nothing copied or cached: the
// public catalogue, the course category pages, the sitemap and the lessons list
// ask for completed courses (landingpages.Completed), so marking a course
// completed lists it and marking it back withdraws it. The course loader writes
// the readiness of a course version (its spec's `progress`), never the
// course's.
type Readiness struct {
	db *ent.Client
}

// NewReadiness builds the readiness editor over an ent client.
func NewReadiness(db *ent.Client) *Readiness {
	return &Readiness{db: db}
}

// MarkDraft takes a course back to a draft.
func (r *Readiness) MarkDraft(ctx context.Context, id int) error {
	return r.move(ctx, id, transition.MarkDraft)
}

// MarkInDevelopment marks a course as being developed.
func (r *Readiness) MarkInDevelopment(ctx context.Context, id int) error {
	return r.move(ctx, id, transition.MarkInDevelopment)
}

// MarkCompleted marks a course as finished, which lists it publicly.
func (r *Readiness) MarkCompleted(ctx context.Context, id int) error {
	return r.move(ctx, id, transition.MarkCompleted)
}

// move checks the move against the table and writes the new readiness through
// UpdateOne, so the audit hook records the old and the new value. A course with
// no readiness (nullable in the legacy schema) counts as a draft. A missing id
// is domain.ErrNotFound and an illegal move a *transition.IllegalError, both before
// any write.
func (r *Readiness) move(ctx context.Context, id int, move transition.Transition[transition.Readiness]) error {
	row, err := r.db.Course.Get(ctx, id)
	if err != nil {
		return domain.FromStore(err)
	}
	if err := move.Check(transition.ReadinessOf(row.Readiness)); err != nil {
		return err
	}
	return domain.FromStore(r.db.Course.UpdateOneID(id).SetReadiness(string(move.To)).Exec(ctx))
}
