package schema

// AdminInput marks a schema whose admin create/update handlers accept a
// contract input type (api.<Schema>Input by default). The adminput template
// (ent/template/adminput.tmpl) generates SetInput methods on the schema's
// Create and UpdateOne builders from this marker, so handlers never hand-map
// contract fields onto builder setters — the mapping regenerates with
// `make gen-ent` whenever the schema changes.
//
// Default field mapping: a Nillable ent field is expected to arrive as an
// ogen Nil* wrapper (SetNillable on create, null-clears-the-column on
// update, matching the legacy assign_attributes semantics); a non-Nillable
// field is expected as a plain value. Contract nullability often diverges
// from legacy column nullability — override per field with AdminInputField.
type AdminInput struct {
	// Type overrides the api input type name; empty means "<Schema>Input".
	Type string `json:"Type"`
	// Explicit flips the default from opt-out to opt-in: only fields carrying
	// an AdminInputField annotation are mapped. Use it on wide schemas (User,
	// Course) where the contract input covers a small slice of the columns,
	// so a new column never leaks into SetInput by omission.
	Explicit bool `json:"Explicit"`
}

// Name implements schema.Annotation.
func (AdminInput) Name() string { return "AdminInput" }

// AdminInputField overrides how a single field of an AdminInput-annotated
// schema maps onto the contract input. The zero value keeps the defaults
// derived from the ent field (see AdminInput) — under Explicit it acts as the
// opt-in marker. created_at/updated_at are always skipped — they are
// ent-owned, never part of a contract input. SetOnly wins over Required, so a
// NOT NULL FK column can still take the "apply only when non-null" semantics.
type AdminInputField struct {
	// Skip excludes the field from SetInput entirely (the contract input has
	// no such field, e.g. legacy columns the admin UI never edits).
	Skip bool `json:"Skip"`
	// Required treats the field as a plain contract value even though the
	// legacy column (and thus the ent field) is nullable.
	Required bool `json:"Required"`
	// SetOnly applies a nullable input field only when non-null and never
	// clears the column — the association semantics used by FK ids, where a
	// null means "leave as is", not "detach".
	SetOnly bool `json:"SetOnly"`
	// Rename overrides the input struct field name when it differs from the
	// ent StructField (e.g. legacy column language_id vs contract CourseId).
	Rename string `json:"Rename"`
}

// Name implements schema.Annotation.
func (AdminInputField) Name() string { return "AdminInputField" }
