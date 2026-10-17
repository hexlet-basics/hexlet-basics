// Package apispec carries the compiled OpenAPI document into the Go binary.
//
// ogen ignores vendor extensions, so the Permission each admin operation
// declares (`x-permission`, ADR-0018) is invisible in internal/api. The server
// reads it from the same generated document the contract is built from, which
// is embedded here because go:embed cannot reach up out of a package directory.
// There is no second generated file to keep in step: regenerating the contract
// (`mise run gen:spec`) is what updates the catalog.
package apispec

import _ "embed"

// OpenAPI is api-spec/dist/openapi.yaml as generated from the TypeSpec.
//
//go:embed dist/openapi.yaml
var OpenAPI []byte
