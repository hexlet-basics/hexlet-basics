package staffroles_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/staffroles"
	"hexletbasics/internal/testsupport"
)

func TestApplyToAMissingRoleIsNotFoundAndWritesNothing(t *testing.T) {
	db := testsupport.NewClient(t)
	before := db.StaffRolePermission.Query().CountX(t.Context())

	err := staffroles.NewPermissions(db).Apply(t.Context(), 0, []api.RolePermissionInput{{Resource: api.PermissionResourceBanners}})

	require.ErrorIs(t, err, domain.ErrNotFound)
	require.Equal(t, before, db.StaffRolePermission.Query().CountX(t.Context()))
}
