// Package apigen holds the code-generation directive for the API server.
// The ogen version is pinned by go.mod (go run uses the module's version).
package apigen

//go:generate go run github.com/ogen-go/ogen/cmd/ogen --target ../api --package api --clean ../../api-spec/dist/openapi.yaml
