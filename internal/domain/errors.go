// Package domain holds the failures every use-case module answers with, so each
// adapter (HTTP, MCP) maps one small set of errors to its own wire format and
// never has to know which ent error means what. It imports nothing of ours:
// the translation from ent's errors lives with the client, in store.ToDomain,
// so a package the generated ent code imports (transition) may use these.
package domain

import (
	"errors"
	"fmt"
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

// FieldError is an ErrInvalid that says which field the input broke.
type FieldError struct {
	Field string
	cause error
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: field %s: %s", ErrInvalid, e.Field, e.cause)
}

// InvalidField is an ErrInvalid naming the field the input broke, with the
// store's failure kept in the chain as its cause.
func InvalidField(field string, cause error) error {
	return &FieldError{Field: field, cause: cause}
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

// IsDomain reports whether err is already in the domain's terms.
func IsDomain(err error) bool {
	return errors.Is(err, ErrNotFound) || errors.Is(err, ErrConflict) ||
		errors.Is(err, ErrInvalid) || errors.Is(err, ErrForbidden)
}
