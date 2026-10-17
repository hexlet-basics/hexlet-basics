// Package blogposts holds the blog-post operations that span several tables
// and so need a transaction of their own.
package blogposts

import (
	"context"
	"database/sql"

	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpostlike"
	"hexletbasics/ent/blogpostrelatedcourseitem"
	"hexletbasics/internal/store"
)

// Remover deletes posts with their dependents.
type Remover struct {
	store store.Transactor
}

// NewRemover builds the production post remover.
func NewRemover(txStore store.Transactor) *Remover {
	return &Remover{store: txStore}
}

// Remove deletes a post with its dependents, mirroring the legacy model
// (`related_language_items dependent: :delete_all`, `likes dependent:
// :destroy`) — the FKs have no ON DELETE CASCADE, so the children go first.
// One transaction: a failure after the children went would otherwise leave a
// post that has lost its likes and related courses.
func (r *Remover) Remove(ctx context.Context, postID int) error {
	return r.store.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		if _, err := db.BlogPostRelatedCourseItem.Delete().
			Where(blogpostrelatedcourseitem.BlogPostID(postID)).Exec(ctx); err != nil {
			return oops.Wrapf(err, "delete related courses of post %d", postID)
		}
		if _, err := db.BlogPostLike.Delete().
			Where(blogpostlike.BlogPostID(postID)).Exec(ctx); err != nil {
			return oops.Wrapf(err, "delete likes of post %d", postID)
		}
		// Left unwrapped: the central handler maps ent's constraint errors.
		return db.BlogPost.DeleteOneID(postID).Exec(ctx)
	})
}
