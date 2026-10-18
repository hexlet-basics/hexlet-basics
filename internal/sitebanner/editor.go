package sitebanner

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/transition"
)

// Editor moves a banner's Publication State. The banner's content is edited by
// the generated admin CRUD, which has no rule of its own; the state is not part
// of it (a state is not an update), so a new banner is a draft by the column
// default and the state changes only here.
//
// The state column is all a transition changes. What the public sees follows
// from it: Directory.Active picks among the published banners of the locale
// whose window is open, so publishing a banner outside its window shows nothing
// until the window opens, and the window itself is edited with the banner, not
// with its state. Legacy kept no cache or stamp behind the banner's state.
type Editor struct {
	db *ent.Client
}

// NewEditor builds the editor over an ent client.
func NewEditor(db *ent.Client) *Editor {
	return &Editor{db: db}
}

// Publish makes a draft or archived banner eligible for the pages.
func (e *Editor) Publish(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Publish())
}

// Unpublish takes a published banner back to draft.
func (e *Editor) Unpublish(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Unpublish())
}

// Archive archives a draft or published banner.
func (e *Editor) Archive(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Archive())
}

// move checks the move against the table and writes the new state through
// UpdateOne, so the audit hook records the old and the new state. A missing id
// is domain.ErrNotFound and an illegal move a *transition.IllegalError, both before
// any write.
func (e *Editor) move(ctx context.Context, id int, move transition.Transition[transition.PublicationState]) error {
	row, err := e.db.Banner.Get(ctx, id)
	if err != nil {
		return domain.FromStore(err)
	}

	if err := move.Check(transition.PublicationState(row.State)); err != nil {
		return err
	}

	return domain.FromStore(e.db.Banner.UpdateOneID(id).SetState(string(move.To)).Exec(ctx))
}
