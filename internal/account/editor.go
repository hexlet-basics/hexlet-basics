package account

import (
	"context"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/store"
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

// Editor creates and changes user records on an administrator's behalf. It
// owns the rules around them, so each surface that edits users (the admin
// pages, the MCP admin tools through them) applies the same ones: the email is
// stored the way registration stores it, or the account could not be found at
// sign-in; only an administrator changes the administrator flag; and the last
// administrator cannot lose it. Removing an account is the Remover's.
type Editor struct {
	db *ent.Client
}

// NewEditor builds the editor over an ent client.
func NewEditor(db *ent.Client) *Editor {
	return &Editor{db: db}
}

// Create adds a user on actor's behalf. Email uniqueness is the baseline unique
// index's job, not a pre-check, so there is no race; a duplicate is an invalid
// input naming the email field (domain.ErrInvalid). A non-administrator
// creating an administrator is ErrAdministratorFlagRestricted.
func (e *Editor) Create(ctx context.Context, actor domain.Actor, fields UserFields) (*ent.User, error) {
	err := flagChange(actor, false, lo.FromPtr(fields.Admin))
	if err != nil {
		return nil, err
	}

	u, err := e.db.User.Create().
		SetEmail(NormalizeEmail(fields.Email)).
		SetNillableFirstName(fields.FirstName).
		SetNillableLastName(fields.LastName).
		SetNillableAdmin(fields.Admin).
		Save(ctx)

	return u, store.ToDomain(err, userUniques()...)
}

// Update replaces a user's editable fields on actor's behalf. A missing id is
// domain.ErrNotFound; a non-administrator changing the flag is
// ErrAdministratorFlagRestricted; taking the flag from the last administrator
// is ErrLastAdmin.
func (e *Editor) Update(ctx context.Context, actor domain.Actor, id int, fields UserFields) (*ent.User, error) {
	current, err := e.db.User.Get(ctx, id)
	if err != nil {
		return nil, store.ToDomain(err)
	}

	err = flagChange(actor, lo.FromPtr(current.Admin), lo.FromPtr(fields.Admin))
	if err != nil {
		return nil, err
	}

	if !lo.FromPtr(fields.Admin) {
		err = EnsureAnotherAdmin(ctx, e.db, id)
		if err != nil {
			return nil, store.ToDomain(err)
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

	return u, store.ToDomain(err, userUniques()...)
}

// flagChange lets only an administrator change the administrator flag. was is
// the flag before the change (false for a new user), will the flag after.
func flagChange(actor domain.Actor, was, will bool) error {
	if was == will || actor.Administrator {
		return nil
	}

	return ErrAdministratorFlagRestricted
}
