// Package domain holds the failures every use-case module answers with, so each
// adapter (HTTP, MCP) maps one small set of errors to its own wire format and
// never has to know which ent error means what.
package domain

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"hexletbasics/ent"
)

var (
	// ErrNotFound is a request for a row that is not there.
	ErrNotFound = errors.New("not found")
	// ErrConflict is a change the data refuses: a unique index, a foreign key, or
	// a rule such as keeping the last administrator.
	ErrConflict = errors.New("conflict")
	// ErrInvalid is input the operation refuses: a value that breaks a rule the
	// person can fix, such as an email another account already has.
	ErrInvalid = errors.New("invalid")
	// ErrForbidden is an operation the caller may not perform.
	ErrForbidden = errors.New("forbidden")
)

// uniqueViolation is PostgreSQL's SQLSTATE for a unique index refusing a row.
const uniqueViolation = "23505"

// Unique names the field behind a unique index, so a violation of it is the
// person's input being invalid (ErrInvalid, field named) instead of a bare
// conflict. A module declares the indexes of its own table; the translator
// stays one, here.
type Unique struct {
	Constraint string
	Field      string
}

// UniqueOn declares that the unique index named constraint guards field.
func UniqueOn(constraint, field string) Unique {
	return Unique{Constraint: constraint, Field: field}
}

// FieldError is an ErrInvalid that says which field the input broke.
type FieldError struct {
	Field string
	cause error
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: field %s: %s", ErrInvalid, e.Field, e.cause)
}

// Unwrap keeps both the domain class and the store's cause in the chain.
func (e *FieldError) Unwrap() []error { return []error{ErrInvalid, e.cause} }

// FieldOf returns the field an invalid-input error names, if it names one.
func FieldOf(err error) (string, bool) {
	fieldErr, ok := errors.AsType[*FieldError](err)
	if !ok {
		return "", false
	}
	return fieldErr.Field, true
}

// FromStore turns ent's failures into the domain's. The cause stays in the
// chain, so logs keep the detail and code that still looks for ent's own error
// finds it. A unique violation on one of the declared indexes is ErrInvalid
// naming its field; any other constraint is ErrConflict. It is the one place
// ent's errors are classified, and it is idempotent: an error already in the
// domain's terms passes through unchanged.
func FromStore(err error, uniques ...Unique) error {
	switch {
	case err == nil:
		return nil
	case isDomain(err):
		return err
	case ent.IsNotFound(err):
		return errors.Join(ErrNotFound, err)
	case ent.IsConstraintError(err):
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == uniqueViolation {
			for _, unique := range uniques {
				if unique.Constraint == pgErr.ConstraintName {
					return &FieldError{Field: unique.Field, cause: err}
				}
			}
		}
		return errors.Join(ErrConflict, err)
	default:
		return err
	}
}

func isDomain(err error) bool {
	return errors.Is(err, ErrNotFound) || errors.Is(err, ErrConflict) ||
		errors.Is(err, ErrInvalid) || errors.Is(err, ErrForbidden)
}
