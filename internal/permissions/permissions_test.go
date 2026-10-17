package permissions_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/permissions"
)

func TestAnOperationReportsTheDeclaredPermissionAndKind(t *testing.T) {
	read, ok := permissions.For("AdminListBlogPosts")
	require.True(t, ok)
	assert.Equal(t, permissions.Requirement{Permission: "blog_posts.read", Kind: permissions.Read}, read)

	write, ok := permissions.For("AdminSuggestBlogPostRelatedCourses")
	require.True(t, ok)
	assert.Equal(t, permissions.Requirement{Permission: "blog_posts.update", Kind: permissions.Write}, write)
}

func TestAnOperationWithoutADeclarationHasNoRequirement(t *testing.T) {
	_, ok := permissions.For("AdminListBanners")
	assert.False(t, ok)
}

func TestTheCatalogListsEachDeclaredPermissionOnceInOrder(t *testing.T) {
	assert.Equal(t, []permissions.Permission{
		{Name: "blog_posts.create", Kind: permissions.Write},
		{Name: "blog_posts.delete", Kind: permissions.Write},
		{Name: "blog_posts.read", Kind: permissions.Read},
		{Name: "blog_posts.update", Kind: permissions.Write},
	}, permissions.Catalog())
}

func TestKnownTellsCatalogNamesFromInventedOnes(t *testing.T) {
	assert.True(t, permissions.Known("blog_posts.delete"))
	assert.False(t, permissions.Known("blog_posts.publish"))
	assert.False(t, permissions.Known(""))
}

func TestMalformedDeclarationsAreRejected(t *testing.T) {
	for name, doc := range map[string]string{
		"not yaml": "paths: [",
		"unknown kind": `
paths:
  /api/admin/x:
    get:
      operationId: adminListX
      x-permission: {name: x.read, kind: peek}`,
		"permission declared with two kinds": `
paths:
  /api/admin/x:
    get:
      operationId: adminListX
      x-permission: {name: x.use, kind: read}
    post:
      operationId: adminCreateX
      x-permission: {name: x.use, kind: write}`,
		"declaration without a name": `
paths:
  /api/admin/x:
    get:
      operationId: adminListX
      x-permission: {kind: read}`,
	} {
		_, err := permissions.Parse([]byte(doc))
		assert.Error(t, err, name)
	}
}
