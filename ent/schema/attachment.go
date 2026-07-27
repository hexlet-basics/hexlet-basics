package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Attachment is the Go stack's own blob-record table (ADR-0005). It replaces the
// three-part ActiveStorage machinery (active_storage_blobs/attachments/
// variant_records) with a single flat table: one row per uploaded file, holding
// the storage key that locates the bytes in the gocloud `blob.Bucket` plus the
// metadata the admin forms echo back. Ownership (which model references which
// attachment) lives on the consumer side as a plain `*_attachment_id` FK, not in
// a polymorphic join — so this table knows nothing about its referrers.
//
// This is a NEW table (no legacy `attachments`), created by its own migration.
// `storage_key` is unique because it is the bucket key we write to and read from;
// a collision would mean two rows pointing at the same bytes.
type Attachment struct {
	ent.Schema
}

func (Attachment) Fields() []ent.Field {
	return []ent.Field{
		// The gocloud bucket key the bytes live under; also drives the read URL.
		field.String("storage_key").Unique().Immutable(),
		field.String("filename"),
		field.String("content_type"),
		field.Int64("byte_size"),
		// Set by ent on insert; attachments are immutable once uploaded.
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}
