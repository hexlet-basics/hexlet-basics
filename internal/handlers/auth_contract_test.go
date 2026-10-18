package handlers_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"hexletbasics/internal/auth"
)

type openAPISecurityScheme struct {
	Type string `yaml:"type"`
	In   string `yaml:"in"`
	Name string `yaml:"name"`
}

type openAPIOperation struct {
	Security  []map[string][]string `yaml:"security"`
	Responses map[string]any        `yaml:"responses"`
}

type openAPIPath struct {
	Get    openAPIOperation `yaml:"get"`
	Post   openAPIOperation `yaml:"post"`
	Delete openAPIOperation `yaml:"delete"`
}

func TestAuthenticationRequirementsAreDeclaredInOpenAPI(t *testing.T) {
	specData, err := os.ReadFile("../../api-spec/dist/openapi.yaml")
	require.NoError(t, err)

	var spec struct {
		Components struct {
			SecuritySchemes map[string]openAPISecurityScheme `yaml:"securitySchemes"`
		} `yaml:"components"`
		Paths map[string]openAPIPath `yaml:"paths"`
	}
	require.NoError(t, yaml.Unmarshal(specData, &spec))

	assert.Equal(t, openAPISecurityScheme{Type: "apiKey", In: "cookie", Name: auth.CookieName},
		spec.Components.SecuritySchemes["UserSession"])
	assert.Equal(t, openAPISecurityScheme{Type: "apiKey", In: "cookie", Name: auth.CookieName},
		spec.Components.SecuritySchemes["AdminSession"])
	assert.NotContains(t, spec.Components.SecuritySchemes, "XsrfToken",
		"cross-origin writes are refused from browser headers, not a token (ADR-0017)")

	assert.Empty(t, spec.Paths["/api/languages"].Get.Security)
	assert.Empty(t, spec.Paths["/api/lessons/{id}/check"].Post.Security)
	assert.Equal(t,
		[]map[string][]string{{"UserSession": {}}},
		spec.Paths["/api/blog_posts/{id}/likes"].Post.Security,
	)
	assert.Equal(t,
		[]map[string][]string{{"UserSession": {}}},
		spec.Paths["/api/my"].Get.Security,
	)
	assert.Equal(t,
		[]map[string][]string{{"AdminSession": {}}},
		spec.Paths["/api/admin/course_categories"].Get.Security,
	)
	assert.Equal(t,
		[]map[string][]string{{"AdminSession": {}}},
		spec.Paths["/api/admin/course_categories"].Post.Security,
	)
	assert.Equal(t,
		[]map[string][]string{{"AdminSession": {}}},
		spec.Paths["/api/admin/attachments"].Post.Security,
	)

	// The OAuth consent screen (ADR-0022): reading it and deciding both need the
	// session; a cross-site decision is refused by the pipeline, not a token.
	consent := spec.Paths["/api/oauth/consents/{requestId}"]
	assert.Equal(t, []map[string][]string{{"UserSession": {}}}, consent.Get.Security)
	assert.Equal(t,
		[]map[string][]string{{"UserSession": {}}},
		spec.Paths["/api/oauth/consents/{requestId}/approve"].Post.Security,
	)
	assert.Equal(t,
		[]map[string][]string{{"UserSession": {}}},
		spec.Paths["/api/oauth/consents/{requestId}/deny"].Post.Security,
	)

	// A refusal is never an operation's own response: the shared default one is
	// Problem Details for every status (ADR-0025).
	assert.Contains(t, spec.Paths["/api/my"].Get.Responses, "default")
	assert.Contains(t, spec.Paths["/api/admin/course_categories"].Get.Responses, "default")
	assert.NotContains(t, spec.Paths["/api/admin/course_categories"].Get.Responses, "403")
}

// TestBlogPostOperationsDeclareAPermissionAndAKind pins ADR-0020's rule for the
// operations migrated so far: what an admin operation needs is written in the
// contract, as a Permission and whether it reads or writes, not in a Go table.
func TestBlogPostOperationsDeclareAPermissionAndAKind(t *testing.T) {
	specData, err := os.ReadFile("../../api-spec/dist/openapi.yaml")
	require.NoError(t, err)

	type declared struct {
		Name string `yaml:"name"`
		Kind string `yaml:"kind"`
	}
	var spec struct {
		Paths map[string]map[string]struct {
			OperationID string    `yaml:"operationId"`
			Permission  *declared `yaml:"x-permission"`
		} `yaml:"paths"`
	}
	require.NoError(t, yaml.Unmarshal(specData, &spec))

	got := map[string]declared{}
	for path, methods := range spec.Paths {
		if !strings.HasPrefix(path, "/api/admin/blog_posts") {
			continue
		}
		for _, op := range methods {
			require.NotNil(t, op.Permission, "%s declares no x-permission", op.OperationID)
			got[op.OperationID] = *op.Permission
		}
	}

	assert.Equal(t, map[string]declared{
		"adminListBlogPosts":                 {"blog_posts.read", "read"},
		"adminGetBlogPost":                   {"blog_posts.read", "read"},
		"adminCreateBlogPost":                {"blog_posts.create", "write"},
		"adminUpdateBlogPost":                {"blog_posts.update", "write"},
		"adminDeleteBlogPost":                {"blog_posts.delete", "write"},
		"adminSetBlogPostRelatedCourses":     {"blog_posts.update", "write"},
		"adminSuggestBlogPostRelatedCourses": {"blog_posts.update", "write"},
		"adminPublishBlogPost":               {"blog_posts.publish", "write"},
		"adminUnpublishBlogPost":             {"blog_posts.unpublish", "write"},
		"adminArchiveBlogPost":               {"blog_posts.archive", "write"},
	}, got)
}
