package handlers_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
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

	assert.Equal(t, openAPISecurityScheme{Type: "apiKey", In: "cookie", Name: "JWT"},
		spec.Components.SecuritySchemes["UserSession"])
	assert.Equal(t, openAPISecurityScheme{Type: "apiKey", In: "cookie", Name: "JWT"},
		spec.Components.SecuritySchemes["AdminSession"])
	assert.Equal(t, openAPISecurityScheme{Type: "apiKey", In: "header", Name: "X-XSRF-TOKEN"},
		spec.Components.SecuritySchemes["XsrfToken"])

	assert.Empty(t, spec.Paths["/languages"].Get.Security)
	assert.Empty(t, spec.Paths["/lessons/{id}/check"].Post.Security)
	assert.Equal(t,
		[]map[string][]string{{"UserSession": {}, "XsrfToken": {}}},
		spec.Paths["/blog_posts/{id}/likes"].Post.Security,
	)
	assert.Equal(t,
		[]map[string][]string{{"UserSession": {}}},
		spec.Paths["/my"].Get.Security,
	)
	assert.Equal(t,
		[]map[string][]string{{"AdminSession": {}}},
		spec.Paths["/admin/course_categories"].Get.Security,
	)
	assert.Equal(t,
		[]map[string][]string{{"AdminSession": {}, "XsrfToken": {}}},
		spec.Paths["/admin/course_categories"].Post.Security,
	)
	assert.Equal(t,
		[]map[string][]string{{"AdminSession": {}, "XsrfToken": {}}},
		spec.Paths["/admin/attachments"].Post.Security,
	)

	assert.Contains(t, spec.Paths["/my"].Get.Responses, "401")
	assert.Contains(t, spec.Paths["/admin/course_categories"].Get.Responses, "401")
	assert.Contains(t, spec.Paths["/admin/course_categories"].Get.Responses, "403")
}
