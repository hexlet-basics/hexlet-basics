package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// ActiveStorageBlob maps the legacy `active_storage_blobs` table — the metadata
// record for one stored file (ActiveStorage). The read side needs `key`, which
// is the bucket key the bytes live under: because the Go blob bucket points at
// the same S3 bucket as ActiveStorage (ADR-0005), the shipped `/storage/{key}`
// read path serves these legacy blobs directly, so a cover URL is just
// `{origin}/storage/{key}`. Only the columns the read side reads are mapped.
type ActiveStorageBlob struct {
	ent.Schema
}

func (ActiveStorageBlob) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "active_storage_blobs"},
	}
}

func (ActiveStorageBlob) Fields() []ent.Field {
	return []ent.Field{
		field.String("key"),
		field.String("filename"),
		field.String("content_type").Optional().Nillable(),
	}
}
