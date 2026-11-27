package store

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"hexletbasics/ent"
	"hexletbasics/internal/domain"
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

// ToDomain turns ent's failures into the domain's. The cause stays in the
// chain, so logs keep the detail and code that still looks for ent's own error
// finds it. A unique violation on one of the declared indexes is ErrInvalid
// naming its field; any other constraint is ErrConflict. It is the one place
// ent's errors are classified, and it is idempotent: an error already in the
// domain's terms passes through unchanged.
func ToDomain(err error, uniques ...Unique) error {
	switch {
	case err == nil:
		return nil
	case domain.IsDomain(err):
		return err
	case ent.IsNotFound(err):
		return errors.Join(domain.ErrNotFound, err)
	case ent.IsConstraintError(err):
		pgErr, ok := errors.AsType[*pgconn.PgError](err)
		if ok && pgErr.Code == uniqueViolation {
			for _, unique := range uniques {
				if unique.Constraint == pgErr.ConstraintName {
					return domain.InvalidField(unique.Field, err)
				}
			}
		}

		return errors.Join(domain.ErrConflict, err)
	default:
		return err
	}
}
