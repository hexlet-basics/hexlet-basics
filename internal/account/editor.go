package account

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/domain"
)

// UserFields is what an administrator can set on a user record. A nil name or
// flag clears the column on an update, matching the legacy assign_attributes
// semantics the admin form relies on.
type UserFields struct {
	Email     string
	FirstName *string
	LastName  *string
	Admin     *bool
}

// Editor creates, changes and deletes user records on an administrator's
// behalf. It owns the rules around them, so each surface that edits users
// (the admin pages, the management pages, any other adapter) applies the same
// ones: the email is stored the way registration stores it, or the account
// could not be found at sign-in, and the last administrator can neither lose
// the flag nor be deleted.
type Editor struct {
	db *ent.Client
}

// NewEditor builds the editor over an ent client.
func NewEditor(db *ent.Client) *Editor {
	return &Editor{db: db}
}

// Create adds a user. Email uniqueness is the baseline unique index's job, not
// a pre-check, so there is no race; a duplicate is an invalid input naming the
// email field (domain.ErrInvalid).
func (e *Editor) Create(ctx context.Context, fields UserFields) (*ent.User, error) {
	u, err := e.db.User.Create().
		SetEmail(NormalizeEmail(fields.Email)).
		SetNillableFirstName(fields.FirstName).
		SetNillableLastName(fields.LastName).
		SetNillableAdmin(fields.Admin).
		Save(ctx)

	return u, domain.FromStore(err, userUniques()...)
}

// Update replaces a user's editable fields. A missing id is domain.ErrNotFound;
// taking the flag from the last administrator is ErrLastAdmin.
func (e *Editor) Update(ctx context.Context, id int, fields UserFields) (*ent.User, error) {
	if fields.Admin == nil || !*fields.Admin {
		err := EnsureAnotherAdmin(ctx, e.db, id)
		if err != nil {
			return nil, domain.FromStore(err)
		}
	}

	update := e.db.User.UpdateOneID(id).SetEmail(NormalizeEmail(fields.Email))
	if fields.FirstName == nil {
		update.ClearFirstName()
	} else {
		update.SetFirstName(*fields.FirstName)
	}

	if fields.LastName == nil {
		update.ClearLastName()
	} else {
		update.SetLastName(*fields.LastName)
	}

	if fields.Admin == nil {
		update.ClearAdmin()
	} else {
		update.SetAdmin(*fields.Admin)
	}

	u, err := update.Save(ctx)

	return u, domain.FromStore(err, userUniques()...)
}

// Delete removes a user. A missing id is domain.ErrNotFound; the last
// administrator is refused with ErrLastAdmin.
func (e *Editor) Delete(ctx context.Context, id int) error {
	err := EnsureAnotherAdmin(ctx, e.db, id)
	if err != nil {
		return domain.FromStore(err)
	}

	return domain.FromStore(e.db.User.DeleteOneID(id).Exec(ctx))
}
