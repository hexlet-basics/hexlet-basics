package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ActiveStorageAttachment maps the legacy `active_storage_attachments` table —
// ActiveStorage's polymorphic join binding an owning record to a blob. A
// BlogPost cover is the row with record_type='BlogPost', name='cover'; the
// handler queries this schema filtered by those columns and follows the `blob`
// edge to read the bucket key. Only the read-side columns are mapped.
type ActiveStorageAttachment struct {
	ent.Schema
}

func (ActiveStorageAttachment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "active_storage_attachments"},
	}
}

func (ActiveStorageAttachment) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Int("record_id"),
		field.String("record_type"),
		field.Int("blob_id"),
	}
}

func (ActiveStorageAttachment) Edges() []ent.Edge {
	return []ent.Edge{
		// belongs-to: the FK (`blob_id`) lives on this table, so the edge owns the
		// field. Loaded via WithBlob to read the blob's bucket key.
		edge.To("blob", ActiveStorageBlob.Type).
			Field("blob_id").
			Unique().
			Required(),
	}
}
