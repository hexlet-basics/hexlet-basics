package coursereview

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/transition"
)

// Editor moves a review's Publication State. The review's content is edited by
// the generated admin CRUD, which has no rule of its own; the state is not part
// of it (a state is not an update), so a new review is a draft by the column
// default and the state changes only here.
//
// The state column is all a transition changes: the public reviews page
// (Public) lists the published reviews of the locale, so publishing is what puts
// a review there. `pinned` is a separate flag of the review, edited with it, and
// legacy kept no cache or stamp behind the state.
type Editor struct {
	db *ent.Client
}

// NewEditor builds the editor over an ent client.
func NewEditor(db *ent.Client) *Editor {
	return &Editor{db: db}
}

// Publish makes a draft or archived review public.
func (e *Editor) Publish(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Publish)
}

// Unpublish takes a published review back to draft.
func (e *Editor) Unpublish(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Unpublish)
}

// Archive archives a draft or published review.
func (e *Editor) Archive(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Archive)
}

// move checks the move against the table and writes the new state through
// UpdateOne, so the audit hook records the old and the new state. A review with
// no state (nullable in the legacy schema) counts as a draft. A missing id is
// domain.ErrNotFound and an illegal move a *transition.IllegalError, both before any
// write.
func (e *Editor) move(ctx context.Context, id int, move transition.Transition[transition.PublicationState]) error {
	row, err := e.db.Review.Get(ctx, id)
	if err != nil {
		return domain.FromStore(err)
	}
	if err := move.Check(transition.CurrentOf(row.State)); err != nil {
		return err
	}
	return domain.FromStore(e.db.Review.UpdateOneID(id).SetState(string(move.To)).Exec(ctx))
}
