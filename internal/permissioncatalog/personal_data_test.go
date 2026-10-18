package permissioncatalog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/permissioncatalog"
)

// Personal data is declared in the contract with x-personal-data, on read
// operations only, and names the ent type the operation reads.
func TestPersonalDataReadsAreReadFromTheContract(t *testing.T) {
	t.Parallel()

	entity, ok := permissioncatalog.PersonalDataEntity("AdminListUsers")
	require.True(t, ok)
	assert.Equal(t, "User", entity)

	for _, operation := range []string{"AdminListLeads", "AdminListMessages", "AdminListLessonProgress", "AdminGetUser", "AdminSearchUsers", "AdminListManagementUsers", "AdminGetManagementUser"} {
		assert.Contains(t, permissioncatalog.PersonalDataReads(), operation)
	}

	for _, operation := range []string{"AdminListCourses", "AdminGetCourse", "AdminListBlogPosts", "AdminListBanners", "AdminListCourseLessons", "AdminCreateBlogPost"} {
		_, ok := permissioncatalog.PersonalDataEntity(operation)
		assert.False(t, ok, "%s reads content, not personal data", operation)
	}

	for _, operation := range permissioncatalog.PersonalDataReads() {
		requirement, ok := permissioncatalog.For(operation)
		require.True(t, ok)
		assert.Equal(t, permissioncatalog.Read, requirement.Kind, "%s: personal data is only read", operation)
	}
}

func TestPersonalDataOnAWriteOrUnclassifiedOperationIsRejected(t *testing.T) {
	t.Parallel()

	for name, doc := range map[string]string{
		"on a write": `
paths:
  /api/admin/x:
    post:
      operationId: adminCreateX
      x-permission: {name: x.create, kind: write}
      x-personal-data: User`,
		"without a permission": `
paths:
  /api/admin/x:
    get:
      operationId: adminListX
      x-personal-data: User`,
		"without an operation id": `
paths:
  /api/admin/x:
    get:
      x-permission: {name: x.read, kind: read}
      x-personal-data: User`,
	} {
		_, err := permissioncatalog.Parse([]byte(doc))
		assert.Error(t, err, name)
	}
}
