package permissioncatalog_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/permissioncatalog"
)

// The audit trail must reach an admin write whatever marks it: a Permission of
// kind write, an administrators-only or any-staff marker, or nothing but the
// admin session scheme. Reads and public operations are never admin writes.
func TestAnAdminWriteIsDecidedFromTheContractNotFromAList(t *testing.T) {
	t.Parallel()
	assert.True(t, permissioncatalog.IsAdminWrite("AdminCreateBlogPost"), "a Permission of kind write")
	assert.True(t, permissioncatalog.IsAdminWrite("AdminUploadAttachment"), "an any-staff marker")
	assert.True(t, permissioncatalog.IsAdminWrite("AdminCreateUser"), "a write Permission on what was administrators-only")
	assert.False(t, permissioncatalog.IsAdminWrite("AdminListBlogPosts"), "a Permission of kind read")
	assert.False(t, permissioncatalog.IsAdminWrite("CreateSession"), "a public operation")

	assert.Contains(t, permissioncatalog.AdminWrites(), "AdminCreateBlogPost")
	assert.True(t, slices.IsSorted(permissioncatalog.AdminWrites()))
}

func TestAnOperationOnlyAnAccessMarkerClassifiesIsJudgedByItsMethod(t *testing.T) {
	t.Parallel()
	d, err := permissioncatalog.Parse([]byte(`
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
	assert.True(t, permissioncatalog.Standalone("AdminUploadAttachment"))
	assert.False(t, permissioncatalog.Standalone("AdminCreateBlogPost"))

	_, err := permissioncatalog.Parse([]byte(`
paths:
  /a:
    post:
      operationId: adminThing
      x-unit-of-work: sometimes
`))
	require.Error(t, err)
}
