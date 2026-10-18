package mcpserver

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/permissioncatalog"
)

func TestEveryAllowlistedOperationIsAnAdminOperationThatDeclaresItsPermission(t *testing.T) {
	require.NotEmpty(t, adminOperations)

	for _, id := range adminOperations {
		op, err := readOperation(id)
		require.NoError(t, err, id)
		assert.True(t, op.Admin, "%s requires the AdminSession scheme, so the checkpoint runs", id)
		name := strings.ToUpper(id[:1]) + id[1:]
		_, declares := permissioncatalog.For(name)
		// An any-staff operation (the upload) declares its access, not a Permission.
		declares = declares || permissioncatalog.AccessOf(name) == permissioncatalog.AnyStaff
		assert.True(t, declares, "%s declares its Permission or access in the contract", id)

		_, err = adminTool(op)
		assert.NoError(t, err, id)
	}
}

func TestAnOperationWithoutTheAdminSessionIsNotOffered(t *testing.T) {
	op, err := readOperation("listCourses")
	require.NoError(t, err)

	_, err = adminTool(op)

	require.ErrorContains(t, err, "not an admin operation")
}

func TestNewRefusesAnAllowlistEntryThatIsNotInTheContract(t *testing.T) {
	was := adminOperations

	t.Cleanup(func() { adminOperations = was })

	adminOperations = []string{"adminNoSuchOperation"}

	err := (&Server{}).addAdminTools(nil)

	require.ErrorContains(t, err, "not in the contract")
}

func TestAToolRefusesAParameterThatCollidesWithAReservedArgument(t *testing.T) {
	op := &operation{ID: "x", Admin: true, Params: []parameter{{Name: "locale", In: "query"}}}

	_, err := adminTool(op)

	require.ErrorContains(t, err, "collides")
}

func TestOutputIsWrappedOnlyWhenItIsNotAnObject(t *testing.T) {
	object := map[string]any{"type": "object"}
	assert.Equal(t, object, wrapObject(object))

	wrapped := wrapObject(map[string]any{"type": "array"})
	assert.Equal(t, "object", wrapped["type"])
	assert.Contains(t, wrapped["properties"], "result")
}

func TestPathParametersAreFoundInTheTemplate(t *testing.T) {
	assert.Equal(t, []string{"id", "lessonId"}, pathParams("/api/a/{id}/b/{lessonId}"))
	assert.Empty(t, pathParams("/api/a"))
	assert.Empty(t, pathParams("/api/a/{broken"))
}

func TestScalarsKeepNumbersAsWritten(t *testing.T) {
	for raw, want := range map[string][]string{
		`"a b"`:   {"a b"},
		`3`:       {"3"},
		`true`:    {"true"},
		`[1,"x"]`: {"1", "x"},
	} {
		got, err := scalars(json.RawMessage(raw))
		require.NoError(t, err, raw)
		assert.Equal(t, want, got, raw)
	}

	for _, raw := range []string{`{}`, `[]`, `[[1]]`, `null`, `nope`} {
		_, err := scalars(json.RawMessage(raw))
		assert.Error(t, err, raw)
	}
}

func TestParseArgumentsAndLocaleAreStrict(t *testing.T) {
	args, problem := parseArguments(nil)
	assert.Empty(t, problem)
	assert.Empty(t, args)

	_, problem = parseArguments(json.RawMessage(`[1]`))
	assert.NotEmpty(t, problem)
}

func TestSourceOfReadsTheVerifiersIdsAndIgnoresWhatIsNotThere(t *testing.T) {
	// The ids of a real call are asserted through its Audit Record; this is the
	// case of a token info that carries none.
	source := sourceOf(&auth.TokenInfo{Extra: map[string]any{}})
	assert.Equal(t, "mcp", string(source.Channel))
	assert.Nil(t, source.TokenID)
	assert.Nil(t, source.ConnectedAppID)
}
