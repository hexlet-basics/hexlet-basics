package mcpserver_test

import (
	"encoding/json"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/mcpserver"
)

// Every admin operation of the contract is a tool or is excluded with a reason,
// so one added later cannot be forgotten.
func TestEveryAdminOperationIsAToolOrExcludedWithAReason(t *testing.T) {
	allowed := mcpserver.AllowlistedAdminOperations()
	excluded := mcpserver.ExcludedAdminOperations()
	contractIDs := mcpserver.ContractAdminOperations(t)

	for _, id := range contractIDs {
		_, isExcluded := excluded[id]
		assert.NotEqual(t, slices.Contains(allowed, id), isExcluded,
			"%s must be in exactly one of the allowlist and the exclusions", id)
	}

	for _, id := range allowed {
		assert.Contains(t, contractIDs, id, "an allowlisted name that is not an admin operation")
	}

	for id, reason := range excluded {
		assert.Contains(t, contractIDs, id, "an exclusion that is not an admin operation")
		assert.NotEmpty(t, reason)
	}

	sorted := slices.Sorted(slices.Values(allowed))
	assert.Len(t, slices.Compact(sorted), len(allowed), "no operation is listed twice")
}

// Every tool is built from the contract, and a representative argument set (the
// required properties, filled from the schema itself) validates against its own
// input schema, so the derivation handles every shape the contract uses.
func TestEveryAdminToolSchemaBuildsAndAcceptsAnExample(t *testing.T) {
	for _, tool := range mcpserver.AdminToolsFromTheContract(t) {
		t.Run(tool.Name, func(t *testing.T) {
			raw, err := json.Marshal(tool.InputSchema)
			require.NoError(t, err)

			var schema jsonschema.Schema
			require.NoError(t, json.Unmarshal(raw, &schema))
			resolved, err := schema.Resolve(nil)
			require.NoError(t, err, "the input schema is a valid JSON Schema")

			var input map[string]any
			require.NoError(t, json.Unmarshal(raw, &input))
			example := exampleOf(input)
			require.NoError(t, resolved.Validate(example), "%v", example)

			if tool.OutputSchema != nil {
				raw, err := json.Marshal(tool.OutputSchema)
				require.NoError(t, err)

				var out jsonschema.Schema
				require.NoError(t, json.Unmarshal(raw, &out))
				_, err = out.Resolve(nil)
				assert.NoError(t, err, "the output schema is a valid JSON Schema")
			}
		})
	}
}

// The tool list is what an agent loads into its context on every connection, so
// its size is a budget: 79 KB for the 87 admin tools when this limit was set (it
// was 229 KB before #878 trimmed the schemas and descriptions). Growing past it is
// a decision to make on purpose, not a side effect of a contract change.
func TestTheAdminToolListStaysWithinItsSizeBudget(t *testing.T) {
	const budget = 90 * 1024

	total := 0

	for _, tool := range mcpserver.AdminToolsFromTheContract(t) {
		raw, err := json.Marshal(tool)
		require.NoError(t, err)

		total += len(raw)

		assert.NotEmpty(t, tool.Description, tool.Name)
	}

	assert.Less(t, total, budget)
}

// exampleOf builds a minimal value of a decoded schema: required properties
// only, the first non-null alternative of a choice, the first enum member.
func exampleOf(schema map[string]any) any {
	if values, ok := schema["enum"].([]any); ok && len(values) > 0 {
		return values[0]
	}

	if parts, ok := schema["allOf"].([]any); ok {
		return allOfExample(parts)
	}

	for _, choice := range []string{"anyOf", "oneOf"} {
		if parts, ok := schema[choice].([]any); ok {
			for _, p := range parts {
				if pm, _ := p.(map[string]any); pm["type"] != "null" {
					return exampleOf(pm)
				}
			}
		}
	}

	typ := schema["type"]
	if list, ok := typ.([]any); ok {
		typ = list[0]
	}

	switch typ {
	case "object":
		return objectExample(schema)
	case "array":
		return []any{}
	case "integer", "number":
		if m, ok := schema["minimum"].(float64); ok {
			return m
		}

		return 1
	case "boolean":
		return true
	default:
		n := 1
		if m, ok := schema["minLength"].(float64); ok && m > 1 {
			n = int(m)
		}

		return strings.Repeat("x", n)
	}
}

// allOfExample merges the examples of an allOf's parts; a scalar part (a named
// enum) is the value itself.
func allOfExample(parts []any) any {
	merged := map[string]any{}

	for _, p := range parts {
		part, _ := p.(map[string]any)
		example := exampleOf(part)

		combined, isObject := example.(map[string]any)
		if !isObject {
			return example
		}

		maps.Copy(merged, combined)
	}

	return merged
}

// objectExample is an object with its required properties only.
func objectExample(schema map[string]any) map[string]any {
	out := map[string]any{}
	props, _ := schema["properties"].(map[string]any)

	required, _ := schema["required"].([]any)
	for _, raw := range required {
		name, _ := raw.(string)
		if p, ok := props[name].(map[string]any); ok {
			out[name] = exampleOf(p)
		}
	}

	return out
}
