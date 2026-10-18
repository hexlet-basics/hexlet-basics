package mcpserver

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	apispec "hexletbasics/api-spec"
)

// operation is one contract operation as an admin tool needs it: how to build
// its HTTP request from tool arguments and what schemas to advertise. It is read
// from the embedded OpenAPI document, the same one internal/permissioncatalog reads
// (ADR-0020), so a tool cannot describe anything the contract does not declare.
type operation struct {
	ID          string
	Method      string
	Path        string
	Description string
	// Params are the path and query parameters, with their references resolved.
	Params []parameter
	// Body is the request body's JSON Schema, nil when the operation takes none.
	Body         map[string]any
	BodyRequired bool
	// Multipart is set instead of a JSON Body when the request is a form with one
	// binary part (the attachment upload): a tool carries the file as base64
	// and the request synthesizer builds the multipart body from it.
	Multipart *multipartFile
	// Response is the JSON Schema of the first 2xx response that has a JSON body,
	// nil when the operation answers without one (a delete).
	Response map[string]any
	// Admin is true when the operation requires the AdminSession scheme: the one
	// whose security handler is the MCP checkpoint.
	Admin bool
}

// multipartFile is the single file part of a multipart/form-data request body.
type multipartFile struct {
	// Field is the part's name, as the contract names it.
	Field string
}

// parameter is a path or query parameter of an operation.
type parameter struct {
	Name        string
	In          string
	Required    bool
	Description string
	Schema      map[string]any
}

type contractDoc struct {
	Paths      map[string]map[string]any `yaml:"paths"`
	Components struct {
		Schemas    map[string]any `yaml:"schemas"`
		Parameters map[string]any `yaml:"parameters"`
	} `yaml:"components"`
}

// contract is parsed once: every tool of the registry reads it.
var contract = sync.OnceValues(func() (*contractDoc, error) {
	var doc contractDoc
	if err := yaml.Unmarshal(apispec.OpenAPI, &doc); err != nil {
		return nil, fmt.Errorf("parse OpenAPI: %w", err)
	}
	return &doc, nil
})

// readOperation finds an operation by operationId and resolves what a tool
// needs of it. A schema that cannot be inlined, or a parameter that is neither a
// path nor a query one, is an error: the registry refuses to start rather than
// advertise a tool it cannot call.
func readOperation(operationID string) (*operation, error) {
	doc, err := contract()
	if err != nil {
		return nil, err
	}
	for path, methods := range doc.Paths {
		for method, raw := range methods {
			node, _ := raw.(map[string]any)
			if node["operationId"] != operationID {
				continue
			}
			return buildOperation(doc, operationID, strings.ToUpper(method), path, node)
		}
	}
	return nil, fmt.Errorf("operation %s is not in the contract", operationID)
}

func buildOperation(doc *contractDoc, id, method, path string, node map[string]any) (*operation, error) {
	op := &operation{ID: id, Method: method, Path: path, Description: asString(node["description"])}

	for _, raw := range asList(node["parameters"]) {
		p, _ := raw.(map[string]any)
		if ref, ok := p["$ref"].(string); ok {
			name, found := strings.CutPrefix(ref, "#/components/parameters/")
			if !found || doc.Components.Parameters[name] == nil {
				return nil, fmt.Errorf("operation %s: unresolvable parameter reference %q", id, ref)
			}
			p, _ = doc.Components.Parameters[name].(map[string]any)
		}
		param := parameter{Name: asString(p["name"]), In: asString(p["in"]), Description: asString(p["description"])}
		param.Required, _ = p["required"].(bool)
		if param.In != "path" && param.In != "query" {
			return nil, fmt.Errorf("operation %s: parameter %s is in %q, which a tool cannot carry", id, param.Name, param.In)
		}
		schema, err := inline(p["schema"], doc.Components.Schemas, nil)
		if err != nil {
			return nil, fmt.Errorf("operation %s: parameter %s: %w", id, param.Name, err)
		}
		param.Schema, _ = schema.(map[string]any)
		op.Params = append(op.Params, param)
	}
	for _, name := range pathParams(path) {
		if !slices.ContainsFunc(op.Params, func(p parameter) bool { return p.In == "path" && p.Name == name }) {
			return nil, fmt.Errorf("operation %s: path parameter %s is not declared", id, name)
		}
	}

	if body, ok := node["requestBody"].(map[string]any); ok {
		file, err := multipartFileOf(body["content"], doc)
		if err != nil {
			return nil, fmt.Errorf("operation %s: request body: %w", id, err)
		}
		if file != nil {
			op.Multipart = file
		} else {
			schema, err := jsonSchemaOf(body["content"], doc)
			if err != nil {
				return nil, fmt.Errorf("operation %s: request body: %w", id, err)
			}
			if schema == nil {
				return nil, fmt.Errorf("operation %s: request body is not JSON", id)
			}
			op.Body = schema
		}
		op.BodyRequired, _ = body["required"].(bool)
	}

	responses, _ := node["responses"].(map[string]any)
	codes := make([]string, 0, len(responses))
	for code := range responses {
		if strings.HasPrefix(code, "2") {
			codes = append(codes, code)
		}
	}
	sort.Strings(codes)
	for _, code := range codes {
		response, _ := responses[code].(map[string]any)
		schema, err := jsonSchemaOf(response["content"], doc)
		if err != nil {
			return nil, fmt.Errorf("operation %s: response %s: %w", id, code, err)
		}
		if schema != nil {
			op.Response = schema
			break
		}
	}

	for _, requirement := range asList(node["security"]) {
		if m, ok := requirement.(map[string]any); ok {
			if _, admin := m["AdminSession"]; admin {
				op.Admin = true
			}
		}
	}
	return op, nil
}

// jsonSchemaOf is the inlined schema of an application/json content entry, nil
// when there is none.
func jsonSchemaOf(content any, doc *contractDoc) (map[string]any, error) {
	entries, _ := content.(map[string]any)
	entry, ok := entries["application/json"].(map[string]any)
	if !ok || entry["schema"] == nil {
		return nil, nil
	}
	resolved, err := inline(entry["schema"], doc.Components.Schemas, nil)
	if err != nil {
		return nil, err
	}
	schema, _ := resolved.(map[string]any)
	return schema, nil
}

// multipartFileOf finds the file part of a multipart/form-data entry, nil when
// the content has none. A form that is not exactly one binary part is an error:
// the base64 tool input carries one file and nothing else, so the registry
// refuses an operation it would build a wrong request for.
func multipartFileOf(content any, doc *contractDoc) (*multipartFile, error) {
	entries, _ := content.(map[string]any)
	entry, ok := entries["multipart/form-data"].(map[string]any)
	if !ok {
		return nil, nil
	}
	resolved, err := inline(entry["schema"], doc.Components.Schemas, nil)
	if err != nil {
		return nil, err
	}
	schema, _ := resolved.(map[string]any)
	properties, _ := schema["properties"].(map[string]any)
	if len(properties) != 1 {
		return nil, errors.New("a multipart form must have exactly one part")
	}
	for name, raw := range properties {
		part, _ := raw.(map[string]any)
		if part["type"] != "string" || part["format"] != "binary" {
			return nil, fmt.Errorf("multipart part %s is not a binary file", name)
		}
		return &multipartFile{Field: name}, nil
	}
	return nil, nil
}

func asList(v any) []any {
	l, _ := v.([]any)
	return l
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}

// pathParams are the names a path template carries, in order.
func pathParams(path string) []string {
	var names []string
	for rest := path; ; {
		start := strings.Index(rest, "{")
		if start < 0 {
			return names
		}
		end := strings.Index(rest[start:], "}")
		if end < 0 {
			return names
		}
		names = append(names, rest[start+1:start+end])
		rest = rest[start+end+1:]
	}
}

// ogenName is the operation's name as the generated server and the permissions
// package know it (AdminListBanners).
func (o *operation) ogenName() string { return strings.ToUpper(o.ID[:1]) + o.ID[1:] }
