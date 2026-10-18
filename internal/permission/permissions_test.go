package permission_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/permission"
	"hexletbasics/internal/testsupport"
)

func TestApplyToAMissingRoleIsNotFoundAndWritesNothing(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	before := db.StaffRoleGrant.Query().CountX(t.Context())

	err := permission.NewPermissions(db, transactor).Apply(t.Context(), 0, []string{"banners.read"})

	require.ErrorIs(t, err, domain.ErrNotFound)
	require.Equal(t, before, db.StaffRoleGrant.Query().CountX(t.Context()))
}

func TestCheckRefusesAGrantTheContractDoesNotDeclare(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	role := fixture.ID(fixture.StaffRoleEditor)

	err := permission.NewPermissions(db, transactor).Check(t.Context(), role, []string{"banners.read", "banners.typo"})

	require.ErrorIs(t, err, permission.ErrUnknownPermission)
	require.ErrorContains(t, err, "banners.typo")
}

func TestCheckOfAMissingRoleIsNotFound(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)

	err := permission.NewPermissions(db, transactor).Check(t.Context(), 0, []string{"banners.typo"})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestApplyKeepsARetiredGrantTheRoleStillHolds(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	role := fixture.ID(fixture.StaffRoleEditor)
	permissions := permission.NewPermissions(db, transactor)
	held, err := db.StaffRoleGrant.Query().Where(staffrolegrant.RoleID(role)).Select(staffrolegrant.FieldPermission).Strings(t.Context())
	require.NoError(t, err)

	require.NoError(t, permissions.Apply(t.Context(), role, append(held, "course_lessons.update")))

	grants, err := db.StaffRoleGrant.Query().Where(staffrolegrant.RoleID(role)).Select(staffrolegrant.FieldPermission).Strings(t.Context())
	require.NoError(t, err)
	require.Contains(t, grants, "course_lessons.update")
}
