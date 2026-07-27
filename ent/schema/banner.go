package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Banner maps the legacy `banners` table (site promo banners shown to visitors).
// Like CourseCategory, admin CRUD writes this table, so the Rails-owned
// `created_at`/`updated_at` (NOT NULL, no DB default) are supplied by ent on
// insert/update. The table name already matches, so no @Table annotation.
//
// `background`/`state`/`locale`/`body` are NOT NULL in the baseline, hence value
// fields (with the DB defaults for the enum columns); `url`/`starts_at`/
// `finishes_at` are nullable, hence Optional().Nillable(). The enum semantics
// live in the API layer (BannerState/BannerBackground/BannerLocale) — the column
// stays a plain string, mirroring the legacy string-backed enums.
type Banner struct {
	ent.Schema
}

func (Banner) Fields() []ent.Field {
	return []ent.Field{
		field.String("background").Default("cta_gradient"),
		field.String("body"),
		field.String("locale"),
		field.String("state").Default("draft"),
		field.String("url").Optional().Nillable(),
		field.Time("starts_at").Optional().Nillable(),
		field.Time("finishes_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
