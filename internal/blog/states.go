package blog

import (
	"context"

	"hexletbasics/ent/blogpost"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/transition"
)

// StatePublished is the `state` of a post open to readers (public.go); the
// legal moves between the Publication States are the shared transition.Publish,
// Unpublish and Archive.

// Publish makes a draft or archived post public. The state column is all there
// is to it: legacy kept no published-at stamp and no feed or cache to refresh
// (the public list, the feeds and the sitemap query the state), so there is no
// side effect to carry.
func (e *Editor) Publish(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Publish)
}

// Unpublish takes a published post back to draft.
func (e *Editor) Unpublish(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Unpublish)
}

// Archive archives a draft or published post.
func (e *Editor) Archive(ctx context.Context, id int) error {
	return e.move(ctx, id, transition.Archive)
}

// move checks the move against the table and writes the new state through
// UpdateOne, so the audit hook records the old and the new state. A row with no
// state (nullable in the legacy schema) counts as a draft, which is what the
// public queries treat it as: not published. A missing id is domain.ErrNotFound
// and an illegal move a *transition.IllegalError, both before any write.
func (e *Editor) move(ctx context.Context, id int, move transition.Transition[transition.PublicationState]) error {
	post, err := e.db.BlogPost.Query().Where(blogpost.IDEQ(id)).Only(ctx)
	if err != nil {
		return domain.FromStore(err)
	}
	if err := move.Check(transition.CurrentOf(post.State)); err != nil {
		return err
	}
	return domain.FromStore(e.db.BlogPost.UpdateOneID(id).SetState(string(move.To)).Exec(ctx))
}
