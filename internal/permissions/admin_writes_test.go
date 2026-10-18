package permissions_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/permissions"
)

// The audit trail must reach an admin write whatever marks it: a Permission of
// kind write, an administrators-only or any-staff marker, or nothing but the
// admin session scheme. Reads and public operations are never admin writes.
func TestAnAdminWriteIsDecidedFromTheContractNotFromAList(t *testing.T) {
	t.Parallel()
	assert.True(t, permissions.IsAdminWrite("AdminCreateBlogPost"), "a Permission of kind write")
	assert.True(t, permissions.IsAdminWrite("AdminUploadAttachment"), "an any-staff marker")
	assert.True(t, permissions.IsAdminWrite("AdminCreateUser"), "a write Permission on what was administrators-only")
	assert.False(t, permissions.IsAdminWrite("AdminListBlogPosts"), "a Permission of kind read")
	assert.False(t, permissions.IsAdminWrite("CreateSession"), "a public operation")

	assert.Contains(t, permissions.AdminWrites(), "AdminCreateBlogPost")
	assert.True(t, slices.IsSorted(permissions.AdminWrites()))
}

func TestAnOperationOnlyAnAccessMarkerClassifiesIsJudgedByItsMethod(t *testing.T) {
	t.Parallel()
	d, err := permissions.Parse([]byte(`
paths:
  /a:
    get:
      operationId: adminReadThing
      security: [{AdminSession: []}]
      x-access: any-staff
    post:
      operationId: adminWriteThing
      security: [{AdminSession: []}]
      x-access: any-staff
    delete:
      operationId: adminBareDelete
      security: [{AdminSession: [], XsrfToken: []}]
    put:
      operationId: publicPut
`))
	require.NoError(t, err)
	assert.Equal(t, []string{"AdminBareDelete", "AdminWriteThing"}, d.AdminWrites())
}

func TestAStandaloneOperationIsReadFromTheContract(t *testing.T) {
	t.Parallel()
	assert.True(t, permissions.Standalone("AdminUploadAttachment"))
	assert.False(t, permissions.Standalone("AdminCreateBlogPost"))

	_, err := permissions.Parse([]byte(`
paths:
  /a:
    post:
      operationId: adminThing
      x-unit-of-work: sometimes
`))
	require.Error(t, err)
}
