// Package domain holds the failures every use-case module answers with, so each
// adapter (HTTP, MCP) maps one small set of errors to its own wire format and
// never has to know which ent error means what.
package domain

import (
	"errors"
	"fmt"

	"hexletbasics/ent"
)

var (
	// ErrNotFound is a request for a row that is not there.
	ErrNotFound = errors.New("not found")
	// ErrConflict is a change the data refuses: a unique index, a foreign key, or
	// a rule such as keeping the last administrator.
	ErrConflict = errors.New("conflict")
)

// FromStore turns ent's failures into the domain's. The cause stays in the
// chain, so logs keep the detail and code that still looks for ent's own error
// finds it.
func FromStore(err error) error {
	switch {
	case err == nil:
		return nil
	case ent.IsNotFound(err):
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	case ent.IsConstraintError(err):
		return fmt.Errorf("%w: %w", ErrConflict, err)
	default:
		return err
	}
}
