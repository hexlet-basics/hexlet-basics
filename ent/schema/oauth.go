package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/lib/pq"
)

// OAuthClient is an application that may ask users for access (ADR-0022): the
// row behind a client_id. Every client is public (no secret) and must use PKCE
// S256, so the table holds nothing a leak could turn into a credential. The
// OAuth tables belong to the oauthserver package and carry no foreign-key edges
// in ent: the package reaches users only through ids it is handed, which is
// what keeps it liftable out of the application.
//
// The table name is pinned because ent would split "OAuth" into o_auth_*.
type OAuthClient struct {
	ent.Schema
}

func (OAuthClient) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "oauth_clients"}}
}

func (OAuthClient) Fields() []ent.Field {
	return []ent.Field{
		// The public identifier the client sends as client_id.
		field.String("client_id").NotEmpty().Unique().Immutable(),
		// What the consent screen calls the application.
		field.String("name").NotEmpty(),
		// Exact-match redirect URIs.
		field.Other("redirect_uris", pq.StringArray{}).
			SchemaType(map[string]string{dialect.Postgres: "varchar[]"}),
		// The custom scopes this client may request. Scopes outside the list are
		// dropped by the library before the consent screen, so a client can never
		// be shown, or granted, more than it was registered for.
		field.Other("scopes", pq.StringArray{}).
			SchemaType(map[string]string{dialect.Postgres: "varchar[]"}),
		// Set only for a client that was created from a Client ID Metadata
		// Document (its client_id is the document's https URL): the instant the
		// cached copy of the document stops being trusted and is fetched again
		// at the next authorization request. Null for a statically seeded or
		// dynamically registered client, which nothing refreshes.
		field.Time("metadata_expires_at").Optional().Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (OAuthClient) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}

// OAuthAuthRequest is one authorization request between /oauth/authorize and
// the token exchange: created when the client sends the user over, completed
// when the user approves, and consumed (deleted) when the code is exchanged.
// The id is a random 128-bit value that travels in the browser URL to the
// consent screen, so it is unguessable, not merely unique.
type OAuthAuthRequest struct {
	ent.Schema
}

func (OAuthAuthRequest) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "oauth_authorization_requests"}}
}

func (OAuthAuthRequest) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Immutable(),
		field.String("client_id").NotEmpty().Immutable(),
		field.String("redirect_uri").NotEmpty().Immutable(),
		field.String("state").Default("").Immutable(),
		field.Other("scopes", pq.StringArray{}).
			SchemaType(map[string]string{dialect.Postgres: "varchar[]"}),
		field.String("code_challenge").NotEmpty().Immutable(),
		field.String("code_challenge_method").NotEmpty().Immutable(),
		field.String("nonce").Default("").Immutable(),
		field.String("response_type").NotEmpty().Immutable(),
		// Set when the user approves; until then the request has no subject.
		field.Int("user_id").Optional().Nillable(),
		field.Bool("done").Default(false),
		field.Time("auth_time").Optional().Nillable(),
		// SHA-256 of the authorization code, hex. The code is single-use and
		// short-lived, but a database read must not hand out a redeemable one.
		field.String("code_hash").Optional().Nillable().Sensitive(),
		field.Time("expires_at").Immutable(),
	}
}

func (OAuthAuthRequest) Indexes() []ent.Index {
	return []ent.Index{index.Fields("code_hash").Unique()}
}

func (OAuthAuthRequest) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}

// OAuthAccessToken is the server-side record of an issued access token. The
// bearer value handed to the client is the encrypted pair (id, subject); the
// table holds only the id, so a database leak yields no usable token, and the
// resource server verifies by looking the id up, so deleting the row revokes
// the token at once (ADR-0022).
type OAuthAccessToken struct {
	ent.Schema
}

func (OAuthAccessToken) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "oauth_access_tokens"}}
}

func (OAuthAccessToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Immutable(),
		field.String("client_id").NotEmpty().Immutable(),
		field.Int("user_id").Immutable(),
		field.Other("scopes", pq.StringArray{}).
			SchemaType(map[string]string{dialect.Postgres: "varchar[]"}),
		field.String("refresh_token_id").Optional().Nillable(),
		// The authorization this token descends from; revoking the family
		// deletes every access token carrying it.
		field.String("family_id").NotEmpty().Immutable(),
		field.Time("expires_at").Immutable(),
	}
}

func (OAuthAccessToken) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}

// OAuthRefreshToken is a refresh token the client holds. Its value is chosen by
// the server's storage, so only a SHA-256 hash is kept: the row can verify a
// presented token but never reproduce one. A rotated token is marked used, not
// deleted: it is the tombstone that makes a replay detectable.
type OAuthRefreshToken struct {
	ent.Schema
}

func (OAuthRefreshToken) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "oauth_refresh_tokens"}}
}

func (OAuthRefreshToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Immutable(),
		field.String("token_hash").NotEmpty().Unique().Sensitive(),
		field.String("client_id").NotEmpty().Immutable(),
		field.Int("user_id").Immutable(),
		field.Other("scopes", pq.StringArray{}).
			SchemaType(map[string]string{dialect.Postgres: "varchar[]"}),
		field.Time("auth_time"),
		field.String("access_token_id").Optional().Nillable(),
		// The authorization the token descends from. Every rotation keeps it, so
		// a replayed refresh token can revoke everything issued since.
		field.String("family_id").NotEmpty().Immutable(),
		// Set when the token is exchanged. A used token stays on file (until its
		// expiry) so that presenting it again is recognised as reuse rather than
		// as an unknown token.
		field.Time("used_at").Optional().Nillable(),
		field.Time("expires_at"),
	}
}

func (OAuthRefreshToken) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
