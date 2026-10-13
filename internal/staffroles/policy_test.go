package staffroles_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"hexletbasics/internal/staffroles"
	"hexletbasics/internal/testsupport"
)

// The fixtures seed bob (1002) as an editor reading blog posts and reading and
// creating banners, in ru and en; alice (1001) is an administrator.
const (
	adminID = 1001
	bobID   = 1002
)

func TestPolicyFollowsTheRoleMatrix(t *testing.T) {
	db := testsupport.NewClient(t)
	policy := staffroles.NewPolicy(db)

	for name, c := range map[string]struct {
		user      int
		locale    string
		operation string
		allowed   bool
	}{
		"administrator may do anything":        {adminID, "es", "AdminDeleteUser", true},
		"granted bit":                          {bobID, "ru", "AdminListBlogPosts", true},
		"another granted bit on a resource":    {bobID, "en", "AdminCreateBanner", true},
		"bit the role lacks":                   {bobID, "ru", "AdminCreateBlogPost", false},
		"resource the role has no row for":     {bobID, "ru", "AdminListReviews", false},
		"locale outside the staff record":      {bobID, "es", "AdminListBlogPosts", false},
		"operation only administrators may do": {bobID, "ru", "AdminListUsers", false},
		"upload needs only being staff":        {bobID, "ru", "AdminUploadAttachment", true},
	} {
		u := db.User.GetX(t.Context(), c.user)
		err := policy.Allow(t.Context(), u, c.locale, c.operation)
		if c.allowed {
			assert.NoError(t, err, name)
		} else {
			assert.ErrorIs(t, err, staffroles.ErrForbidden, name)
		}
	}
}

func TestPolicyRefusesAUserWhoIsNotStaff(t *testing.T) {
	db := testsupport.NewClient(t)
	u := db.User.Create().SetEmail("plain@example.test").SaveX(t.Context())

	err := staffroles.NewPolicy(db).Allow(t.Context(), u, "ru", "AdminListBanners")

	assert.ErrorIs(t, err, staffroles.ErrForbidden)
}

// TestEveryAdminOperationIsClassified fails when a new admin operation is added
// to the contract without saying who may use it: it would be closed to staff
// by default, which is safe but surprising.
func TestEveryAdminOperationIsClassified(t *testing.T) {
	raw, err := os.ReadFile("../../api-spec/dist/openapi.yaml")
	require.NoError(t, err)
	var spec struct {
		Paths map[string]map[string]struct {
			OperationID string                `yaml:"operationId"`
			Security    []map[string][]string `yaml:"security"`
		} `yaml:"paths"`
	}
	require.NoError(t, yaml.Unmarshal(raw, &spec))

	var unclassified []string
	for path, methods := range spec.Paths {
		if !strings.HasPrefix(path, "/api/admin/") {
			continue
		}
		for _, op := range methods {
			if op.OperationID != "" && !staffroles.IsClassified(staffroles.OperationName(op.OperationID)) {
				unclassified = append(unclassified, op.OperationID)
			}
		}
	}
	assert.Empty(t, unclassified, "classify these in staffroles/rules.go")
}
