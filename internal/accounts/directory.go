package accounts

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/domain"
)

// searchLimit mirrors the legacy typeahead cap (`.limit(20)`).
const searchLimit = 20

// Directory reads user records for an administrator: the list, the typeahead and
// a single record. Reads are kept apart from the Editor's writes, so a surface
// that may only look cannot reach the rules that change things.
type Directory struct {
	db *ent.Client
}

// NewDirectory builds the directory over an ent client.
func NewDirectory(db *ent.Client) *Directory {
	return &Directory{db: db}
}

// List returns one page of users, newest first, and how many there are in all.
// Newest first is id descending: ids only grow, so the order is total and the
// pages never overlap.
func (d *Directory) List(ctx context.Context, page domain.Page) ([]*ent.User, int, error) {
	return domain.Paged(ctx, d.db.User.Query().Order(ent.Desc(user.FieldID)), page)
}

// Search is the admin typeahead: a case-insensitive substring match on first
// name, last name or email (legacy `first_name_or_last_name_or_email_cont`),
// capped at 20 rows. A blank term returns nothing rather than every user.
func (d *Directory) Search(ctx context.Context, term string) ([]*ent.User, error) {
	if term == "" {
		return nil, nil
	}
	rows, err := d.db.User.Query().
		Where(user.Or(
			user.FirstNameContainsFold(term),
			user.LastNameContainsFold(term),
			user.EmailContainsFold(term),
		)).
		Order(ent.Desc(user.FieldID)).
		Limit(searchLimit).
		All(ctx)
	return rows, domain.FromStore(err)
}

// Get returns one user. A missing id is domain.ErrNotFound.
func (d *Directory) Get(ctx context.Context, id int) (*ent.User, error) {
	u, err := d.db.User.Get(ctx, id)
	return u, domain.FromStore(err)
}
