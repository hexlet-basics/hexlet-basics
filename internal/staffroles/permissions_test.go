package staffroles_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"hexletbasics/internal/domain"
	"hexletbasics/internal/staffroles"
	"hexletbasics/internal/testsupport"
)

func TestApplyToAMissingRoleIsNotFoundAndWritesNothing(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	before := db.StaffRoleGrant.Query().CountX(t.Context())

	err := staffroles.NewPermissions(db, transactor).Apply(t.Context(), 0, []string{"banners.read"})

	require.ErrorIs(t, err, domain.ErrNotFound)
	require.Equal(t, before, db.StaffRoleGrant.Query().CountX(t.Context()))
}
