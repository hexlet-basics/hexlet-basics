package mcpserver

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	apispec "hexletbasics/api-spec"
)

// responseSchema is the JSON Schema of an operation's 200 response body, read
// from the embedded OpenAPI document (the same one internal/permissions reads,
// ADR-0019) with every $ref inlined.
//
// A tool's output schema comes from the contract rather than from reflection on
// the Go types: ogen wraps optional fields in Opt* structs that reflection would
// render as objects with Value and Set properties, which no client could match
// to the JSON the tool really returns. Reading the document also means a tool
// cannot drift from the operation it mirrors: regenerating the contract is what
// updates the schema (ADR-0020).
func responseSchema(operationID string) (map[string]any, error) {
	var doc struct {
		Paths      map[string]map[string]yaml.Node `yaml:"paths"`
		Components struct {
			Schemas map[string]any `yaml:"schemas"`
		} `yaml:"components"`
	}
	if err := yaml.Unmarshal(apispec.OpenAPI, &doc); err != nil {
		return nil, fmt.Errorf("parse OpenAPI: %w", err)
	}

	for _, methods := range doc.Paths {
		for _, node := range methods {
			var op struct {
				OperationID string `yaml:"operationId"`
				Responses   map[string]struct {
					Content map[string]struct {
						Schema any `yaml:"schema"`
					} `yaml:"content"`
				} `yaml:"responses"`
			}
			if err := node.Decode(&op); err != nil || op.OperationID != operationID {
				continue
			}
			body, ok := op.Responses["200"].Content["application/json"]
			if !ok || body.Schema == nil {
				return nil, fmt.Errorf("operation %s declares no JSON 200 response", operationID)
			}
			resolved, err := inline(body.Schema, doc.Components.Schemas, nil)
			if err != nil {
				return nil, fmt.Errorf("operation %s: %w", operationID, err)
			}
			schema, _ := resolved.(map[string]any)
			return schema, nil
		}
	}
	return nil, fmt.Errorf("operation %s is not in the contract", operationID)
}

// inline replaces every `$ref` under node with the schema it names, and turns
// OpenAPI 3.0's `nullable: true` into the JSON Schema type list that MCP clients
// validate against. A cyclic reference is an error: a tool schema must be a
// finite document, and nothing in the catalogue is recursive.
func inline(node any, components map[string]any, stack []string) (any, error) {
	switch n := node.(type) {
	case map[string]any:
		if ref, ok := n["$ref"].(string); ok {
			name, found := strings.CutPrefix(ref, "#/components/schemas/")
			if !found {
				return nil, fmt.Errorf("unsupported reference %q", ref)
			}
			for _, seen := range stack {
				if seen == name {
					return nil, fmt.Errorf("reference cycle through %s", name)
				}
			}
			target, ok := components[name]
			if !ok {
				return nil, fmt.Errorf("unknown schema %q", name)
			}
			return inline(target, components, append(stack, name))
		}
		out := make(map[string]any, len(n))
		for k, v := range n {
			resolved, err := inline(v, components, stack)
			if err != nil {
				return nil, err
			}
			out[k] = resolved
		}
		if nullable, _ := out["nullable"].(bool); nullable {
			delete(out, "nullable")
			return allowNull(out), nil
		}
		delete(out, "nullable")
		return out, nil
	case []any:
		out := make([]any, len(n))
		for i, v := range n {
			resolved, err := inline(v, components, stack)
			if err != nil {
				return nil, err
			}
			out[i] = resolved
		}
		return out, nil
	default:
		return node, nil
	}
}

// allowNull makes a schema accept null, which is what OpenAPI 3.0's
// `nullable: true` means. A plain type becomes a type list, and an enum gains
// null. A schema built from other schemas (the contract writes a nullable
// reference as `allOf` with `nullable`) cannot take null by editing it, since
// every branch would still refuse it, so null becomes an alternative to it.
func allowNull(schema map[string]any) map[string]any {
	for _, combined := range []string{"allOf", "anyOf", "oneOf"} {
		if _, ok := schema[combined]; ok {
			return map[string]any{"anyOf": []any{schema, map[string]any{"type": "null"}}}
		}
	}
	if typ, ok := schema["type"].(string); ok {
		schema["type"] = []any{typ, "null"}
	}
	if values, ok := schema["enum"].([]any); ok {
		schema["enum"] = append(values, nil)
	}
	return schema
}
