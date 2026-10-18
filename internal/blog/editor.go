package blog

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/transition"
)

// PostFields is what an administrator writes on a post. A nil field clears the
// column on an update, matching the legacy assign_attributes semantics. The
// cover is not here: it is deferred until blob covers land.
type PostFields struct {
	Name        *string
	Slug        *string
	Description *string
	// RichBody is trusted editor HTML, stored and returned exactly as given.
	RichBody string
}

// Editor creates and changes posts on an administrator's behalf.
type Editor struct {
	db *ent.Client
}

// NewEditor builds the editor over an ent client.
func NewEditor(db *ent.Client) *Editor {
	return &Editor{db: db}
}

// Create adds a post, always as a draft: the state is not part of the fields, it
// moves only through Publish, Unpublish and Archive. The creator is the author
// (legacy `blog_post.creator = current_user`), and the post is stamped with the
// locale of the request: legacy stamped the admin UI's.
func (e *Editor) Create(ctx context.Context, creatorID int, locale string, fields PostFields) (*ent.BlogPost, error) {
	row, err := e.db.BlogPost.Create().
		SetNillableName(fields.Name).
		SetNillableSlug(fields.Slug).
		SetNillableDescription(fields.Description).
		SetState(string(transition.Draft)).
		SetRichBody(fields.RichBody).
		SetCreatorID(creatorID).
		SetLocale(locale).
		Save(ctx)
	return row, domain.FromStore(err)
}

// Update replaces a post's fields; a nil one clears its column. The stored
// locale is left untouched (a divergence from legacy, which re-stamped the
// request locale on every update: an accident of its service signature, not
// behavior worth keeping). A missing id is domain.ErrNotFound.
func (e *Editor) Update(ctx context.Context, id int, fields PostFields) (*ent.BlogPost, error) {
	update := e.db.BlogPost.UpdateOneID(id).SetRichBody(fields.RichBody)
	if fields.Name == nil {
		update.ClearName()
	} else {
		update.SetName(*fields.Name)
	}
	if fields.Slug == nil {
		update.ClearSlug()
	} else {
		update.SetSlug(*fields.Slug)
	}
	if fields.Description == nil {
		update.ClearDescription()
	} else {
		update.SetDescription(*fields.Description)
	}
	row, err := update.Save(ctx)
	return row, domain.FromStore(err)
}
