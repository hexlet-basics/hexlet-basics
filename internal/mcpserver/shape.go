package mcpserver

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// The tool list is what an agent loads into its context on every connection, and
// most of its weight is schema: the contract's records inlined once per tool.
// These functions cut it down to what a client needs (ADR-0021 addendum, #878).
//
// They only ever loosen a schema or shorten its prose. A result that satisfied
// the contract's schema satisfies the shape advertised for it, and an argument
// the contract refuses is still refused (by the generated server, which keeps
// the full schema: the advertised input is a description, not the validator).

const (
	// maxOutputDepth is how deep an output shape keeps its structure: the result
	// object (0) lists its fields (1), each as its type. A client needs to know
	// that a list is under `items` and a record has an `id`, not every nested
	// field of a Course; the result itself shows the rest.
	maxOutputDepth = 1

	// maxDescription is the most a property description keeps, cut at a sentence
	// when one ends earlier. Clients cut descriptions at fixed lengths measured
	// from the start, so the first sentence is the one that has to carry the point.
	maxDescription = 100

	// maxToolDescription is the same for a tool's own description, which is read
	// when choosing a tool and so may be a little longer.
	maxToolDescription = 140
)

// firstSentence is text up to the end of its first sentence, and at most limit
// runes: the rest is detail an agent reads from a refusal when it needs it.
func firstSentence(text string, limit int) string {
	text = strings.Join(strings.Fields(text), " ")
	if end := strings.Index(text, ". "); end >= 0 {
		text = text[:end+1]
	}
	if utf8.RuneCountInString(text) <= limit {
		return text
	}
	cut := []rune(text)[:limit-1]
	return strings.TrimSpace(string(cut)) + "…"
}

// trimInput is an input schema without what only decorates it: titles, examples,
// vendor extensions, and the long part of every description. Constraints stay,
// because an agent that sees `minLength` or an enum writes a valid call at the
// first try.
func trimInput(node any) any {
	switch n := node.(type) {
	case map[string]any:
		out := make(map[string]any, len(n))
		for k, v := range n {
			switch {
			case k == "title" || k == "example" || k == "examples" || strings.HasPrefix(k, "x-"):
				continue
			case k == "format" && slices.Contains([]any{"int32", "int64", "float", "double"}, v):
				// The width of a number says nothing an agent can use; date-time
				// and email do, and stay.
				continue
			case k == "description":
				if text, ok := v.(string); ok {
					out[k] = firstSentence(text, maxDescription)
				}
			case k == "properties":
				properties, _ := v.(map[string]any)
				trimmed := make(map[string]any, len(properties))
				for name, property := range properties {
					trimmed[name] = trimInput(property)
				}
				out[k] = trimmed
			case k == "enum" || k == "required" || k == "const" || k == "default":
				out[k] = v
			default:
				out[k] = trimInput(v)
			}
		}
		return out
	case []any:
		out := make([]any, len(n))
		for i, v := range n {
			out[i] = trimInput(v)
		}
		return out
	default:
		return node
	}
}

// outputShape is an output schema reduced to the structure of the result: types,
// property names and list items, to maxOutputDepth. Prose,
// formats, enums and bounds are dropped, since a client reads values, it does not
// validate them against ours. Every reduction accepts more than the contract's
// schema, so a real result still validates against it.
//
// The root is always a plain object: MCP asks for one, and a root `anyOf` or
// `oneOf` is flattened by Claude Code, which then loses the enforcement it
// would have given.
func outputShape(schema map[string]any) map[string]any {
	shape := shapeOf(schema, 0)
	if shape["type"] == nil {
		shape["type"] = "object"
	}
	delete(shape, "anyOf")
	delete(shape, "oneOf")
	delete(shape, "allOf")
	return shape
}

func shapeOf(node any, depth int) map[string]any {
	n, _ := node.(map[string]any)
	out := map[string]any{}
	if t, ok := n["type"]; ok {
		out["type"] = t
	}
	if depth >= maxOutputDepth {
		return out
	}
	if properties, ok := n["properties"].(map[string]any); ok {
		shaped := make(map[string]any, len(properties))
		for name, property := range properties {
			shaped[name] = shapeOf(property, depth+1)
		}
		out["properties"] = shaped
	}
	if items, ok := n["items"]; ok {
		out["items"] = shapeOf(items, depth+1)
	}
	for _, combined := range []string{"anyOf", "oneOf", "allOf"} {
		if parts, ok := n[combined].([]any); ok {
			shaped := make([]any, len(parts))
			for i, part := range parts {
				shaped[i] = shapeOf(part, depth+1)
			}
			out[combined] = shaped
		}
	}
	return out
}
