package landingpage

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/transition"
)

// Editor moves a landing page's Publication State. The page's content is edited
// by the generated admin CRUD, which has no rule of its own; the state is not
// part of it (a state is not an update), so a new page is a draft by the column
// default and the state changes only here.
//
// The state column is all a transition changes, and every public read follows
// from it through Published (the catalog, the category pages, the Yandex feed
// and the sitemap), so publishing a page also lists it under its Course
// Categories through the category items it already has: nothing is copied or
// cached. Legacy additionally validated a description on a published page; the
// Go update never enforced that or any other presence rule of the page, so it is
// not carried over by the transition either.
type Editor struct {
	db *ent.Client
}

// NewEditor builds the editor over an ent client.
func NewEditor(db *ent.Client) *Editor {
	return &Editor{db: db}
}

// Publish makes a draft or archived landing page public.
func (e *Editor) Publish(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Publish)
}

// Unpublish takes a published landing page back to draft.
func (e *Editor) Unpublish(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Unpublish)
}

// Archive archives a draft or published landing page.
func (e *Editor) Archive(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Archive)
}

// move checks the move against the table and writes the new state through
// UpdateOne, so the audit hook records the old and the new state. A page with no
// state (nullable in the legacy schema) counts as a draft. A missing id is
// domain.ErrNotFound and an illegal move a *transition.IllegalError, both before any
// write.
func (e *Editor) move(ctx context.Context, id int, move transition.Transition[transition.PublicationState]) error {
	row, err := e.db.CourseLandingPage.Get(ctx, id)
	if err != nil {
		return domain.FromStore(err)
	}

	if err := move.Check(transition.CurrentOf(row.State)); err != nil {
		return err
	}

	return domain.FromStore(e.db.CourseLandingPage.UpdateOneID(id).SetState(string(move.To)).Exec(ctx))
}
