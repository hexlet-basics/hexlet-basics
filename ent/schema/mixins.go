package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// TimestampsMixin supplies the Rails-owned `created_at`/`updated_at` pair for
// tables ent writes to. The columns are NOT NULL with no DB default in the
// atlas baseline (Rails supplied the values application-side), so ent must set
// them on insert/update exactly like ActiveRecord did. Read-only schemas that
// never insert (e.g. Lead, BlogPost) declare their timestamps directly without
// defaults instead of using this mixin.
type TimestampsMixin struct {
	mixin.Schema
}

func (TimestampsMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
