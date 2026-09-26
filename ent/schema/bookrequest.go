package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BookRequest maps the legacy `book_requests` table: a user's request for the
// free book and whether they have downloaded it since. It is an ent schema on
// the table the atlas baseline already has, so it adds no migration. The table
// name already matches ent's default plural of `BookRequest`.
type BookRequest struct {
	ent.Schema
}

// Mixin supplies Rails' application-side timestamps: both columns are NOT NULL
// with no DB default, so ent must fill them on insert like ActiveRecord did.
func (BookRequest) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}

func (BookRequest) Fields() []ent.Field {
	return []ent.Field{
		// NOT NULL in the baseline (FK to users), hence a value field.
		field.Int("user_id"),
		// Legacy's typed_enum: `requested` on creation, `downloaded` once the
		// PDF was fetched. The column is nullable varchar with no DB default, so
		// the value set lives in internal/books rather than in a DB enum.
		field.String("state").Optional().Nillable(),
	}
}

// Indexes mirrors `index_book_requests_on_user_id`: one request per user,
// which is what makes a repeated request a no-op.
func (BookRequest) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id").Unique(),
	}
}
