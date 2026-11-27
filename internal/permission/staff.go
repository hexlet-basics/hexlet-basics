package permission

import (
	"context"

	"github.com/lib/pq"

	"hexletbasics/ent"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/store"
)

// StaffMemberFields is what an administrator writes on a Staff Member record:
// the user, the Staff Role they hold and the locales they work in.
type StaffMemberFields struct {
	UserID         int
	RoleID         int
	AllowedLocales []string
}

// StaffEditor writes the Staff Member records, the one way a user comes to
// hold a Staff Role. No one changes their own record, administrators included
// (another administrator can), and no one moves a record onto themselves: the
// rules are here, with the write, so no adapter writes past them. One record
// per user is the unique index's rule, a duplicate a domain.ErrConflict.
type StaffEditor struct {
	db *ent.Client
}

// NewStaffEditor builds the editor over an ent client.
func NewStaffEditor(db *ent.Client) *StaffEditor {
	return &StaffEditor{db: db}
}

// Create gives a user a Staff Role. Giving it to oneself is ErrOwnStaffRecord.
func (e *StaffEditor) Create(ctx context.Context, actor domain.Actor, fields StaffMemberFields) (*ent.StaffMember, error) {
	if fields.UserID == actor.UserID {
		return nil, ErrOwnStaffRecord
	}

	row, err := e.db.StaffMember.Create().
		SetUserID(fields.UserID).
		SetRoleID(fields.RoleID).
		SetAllowedLocales(pq.StringArray(fields.AllowedLocales)).
		Save(ctx)

	return row, store.ToDomain(err)
}

// Update replaces a record. A missing one is domain.ErrNotFound; the actor's
// own record, or one the change would move onto the actor, is
// ErrOwnStaffRecord.
func (e *StaffEditor) Update(ctx context.Context, actor domain.Actor, id int, fields StaffMemberFields) (*ent.StaffMember, error) {
	err := e.notOwn(ctx, actor, id)
	if err != nil {
		return nil, err
	}

	if fields.UserID == actor.UserID {
		return nil, ErrOwnStaffRecord
	}

	row, err := e.db.StaffMember.UpdateOneID(id).
		SetUserID(fields.UserID).
		SetRoleID(fields.RoleID).
		SetAllowedLocales(pq.StringArray(fields.AllowedLocales)).
		Save(ctx)

	return row, store.ToDomain(err)
}

// Delete takes a Staff Role from its holder. A missing record is
// domain.ErrNotFound; the actor's own is ErrOwnStaffRecord.
func (e *StaffEditor) Delete(ctx context.Context, actor domain.Actor, id int) error {
	err := e.notOwn(ctx, actor, id)
	if err != nil {
		return err
	}

	return store.ToDomain(e.db.StaffMember.DeleteOneID(id).Exec(ctx))
}

// notOwn refuses a change to a record the actor holds.
func (e *StaffEditor) notOwn(ctx context.Context, actor domain.Actor, id int) error {
	record, err := e.db.StaffMember.Get(ctx, id)
	if err != nil {
		return store.ToDomain(err)
	}

	if record.UserID == actor.UserID {
		return ErrOwnStaffRecord
	}

	return nil
}
