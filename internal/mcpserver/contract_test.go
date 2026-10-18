package mcpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseSchemaOfAnOperationIsFullyInlined(t *testing.T) {
	schema, err := responseSchema("listCourses")

	require.NoError(t, err)
	assert.Equal(t, "array", schema["type"])
	assert.NotContains(t, schema, "$ref")
	item := as[map[string]any](t, schema["items"])
	assert.Equal(t, "object", item["type"])
	// A nullable property becomes a type list, and a referenced schema its body.
	assert.Equal(t, []any{"string", "null"}, as[map[string]any](t, as[map[string]any](t, item["properties"])["name"])["type"])
	assert.Equal(t, "object", as[map[string]any](t, as[map[string]any](t, item["properties"])["course"])["type"])
}

func TestResponseSchemaRefusesWhatTheContractDoesNotDeclare(t *testing.T) {
	_, err := responseSchema("noSuchOperation")
	require.ErrorContains(t, err, "not in the contract")

	_, err = responseSchema("adminDeleteBanner")
	require.Error(t, err, "an operation without a JSON 200 body has no output schema")
}

func TestInlineRefusesReferencesItCannotResolve(t *testing.T) {
	components := map[string]any{
		"Loop": map[string]any{"properties": map[string]any{"next": map[string]any{"$ref": "#/components/schemas/Loop"}}},
	}
	ref := func(to string) map[string]any { return map[string]any{"$ref": to} }

	_, err := inline(ref("#/components/schemas/Loop"), components, nil)
	require.ErrorContains(t, err, "cycle")

	_, err = inline(ref("#/components/schemas/Missing"), components, nil)
	require.ErrorContains(t, err, "unknown schema")

	_, err = inline([]any{ref("#/definitions/Elsewhere")}, components, nil)
	require.ErrorContains(t, err, "unsupported reference")
}
