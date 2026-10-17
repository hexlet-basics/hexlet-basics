package blogposts

import (
	"context"

	"hexletbasics/ent/blogpost"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/transitions"
)

// The Publication States of a post (GLOSSARY). The ent field is a plain string,
// so the values the rules compare and write are named here, apart from the
// contract's enum. StatePublished is in public.go, beside the reader's query.
const (
	StateDraft    = "draft"
	StateArchived = "archived"
)

// The legal moves of a post's Publication State. Each operation names the
// states it may start from; repeating the current state is not among them, so it
// is refused like any illegal move (transitions.Illegal).
var (
	publishPost = transitions.Transition[string]{
		To:   StatePublished,
		From: []string{StateDraft, StateArchived},
	}
	unpublishPost = transitions.Transition[string]{
		To:   StateDraft,
		From: []string{StatePublished},
	}
	archivePost = transitions.Transition[string]{
		To:   StateArchived,
		From: []string{StateDraft, StatePublished},
	}
)

// Publish makes a draft or archived post public. The state column is all there
// is to it: legacy kept no published-at stamp and no feed or cache to refresh
// (the public list, the feeds and the sitemap query the state), so there is no
// side effect to carry.
func (e *Editor) Publish(ctx context.Context, id int) error {
	return e.move(ctx, id, publishPost)
}

// Unpublish takes a published post back to draft.
func (e *Editor) Unpublish(ctx context.Context, id int) error {
	return e.move(ctx, id, unpublishPost)
}

// Archive archives a draft or published post.
func (e *Editor) Archive(ctx context.Context, id int) error {
	return e.move(ctx, id, archivePost)
}

// move checks the move against the table and writes the new state through
// UpdateOne, so the audit hook records the old and the new state. A row with no
// state (nullable in the legacy schema) counts as a draft, which is what the
// public queries treat it as: not published. A missing id is domain.ErrNotFound
// and an illegal move a *transitions.Illegal, both before any write.
func (e *Editor) move(ctx context.Context, id int, move transitions.Transition[string]) error {
	post, err := e.db.BlogPost.Query().Where(blogpost.IDEQ(id)).Only(ctx)
	if err != nil {
		return domain.FromStore(err)
	}

	current := StateDraft
	if post.State != nil {
		current = *post.State
	}
	if err := move.Check(current); err != nil {
		return err
	}

	return domain.FromStore(e.db.BlogPost.UpdateOneID(id).SetState(move.To).Exec(ctx))
}
