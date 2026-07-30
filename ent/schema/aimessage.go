package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AiMessage maps the legacy `ai_messages` table: one message of an in-lesson
// assistant chat. Only the product columns are declared: role/content identify
// the message, user_id marks who asked (null on assistant replies), and the
// input/output token counts keep per-message cost visible. The RubyLLM
// internals (content_raw, thinking_*, cache tokens, ai_model_id,
// ai_tool_call_id) stay in the table but out of the Go domain — modeling them
// would freeze one Ruby library's persistence format into the new stack (see
// the assistant port design). The assistant will write this table, so the
// Rails-owned timestamps come from TimestampsMixin. No @Table annotation: the
// struct name's snake-plural already matches the table.
type AiMessage struct {
	ent.Schema
}

func (AiMessage) Fields() []ent.Field {
	return []ent.Field{
		field.Int("ai_chat_id"),
		field.String("role"),
		field.Text("content").Optional().Nillable(),
		field.Int("user_id").Optional().Nillable(),
		field.Int("input_tokens").Optional().Nillable(),
		field.Int("output_tokens").Optional().Nillable(),
	}
}

func (AiMessage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("chat", AiChat.Type).
			Field("ai_chat_id").
			Unique().
			Required(),
	}
}

func (AiMessage) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
