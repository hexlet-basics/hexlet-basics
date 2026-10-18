// Package ident generates identifiers.
package ident

import "github.com/oklog/ulid/v2"

// New returns a new lexicographically sortable ULID as a string.
func New() string {
	return ulid.Make().String()
}
