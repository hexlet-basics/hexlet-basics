-- Tables of the in-process OAuth authorization server (ADR-0022). They are owned
-- by internal/oauthserver, which reaches users only through ids it is handed.
-- The foreign keys to users make removing an account revoke every grant it made;
-- the one to oauth_clients does the same when a client is deleted.
--
-- Instants that are compared with "now" (expires_at, auth_time) are timestamptz:
-- a plain timestamp drops the offset, and a token's lifetime must not depend on
-- the server's zone. created_at/updated_at follow the other tables.

CREATE TABLE "oauth_clients" (
  "id" bigserial NOT NULL,
  "client_id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "redirect_uris" character varying[] NOT NULL DEFAULT '{}',
  "scopes" character varying[] NOT NULL DEFAULT '{}',
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "index_oauth_clients_on_client_id" ON "oauth_clients" ("client_id");

-- One authorization request between /oauth/authorize and the token exchange.
-- code_hash is the SHA-256 of the authorization code, never the code itself.
CREATE TABLE "oauth_authorization_requests" (
  "id" character varying NOT NULL,
  "client_id" character varying NOT NULL REFERENCES "oauth_clients" ("client_id") ON DELETE CASCADE,
  "redirect_uri" character varying NOT NULL,
  "state" character varying NOT NULL DEFAULT '',
  "scopes" character varying[] NOT NULL DEFAULT '{}',
  "code_challenge" character varying NOT NULL,
  "code_challenge_method" character varying NOT NULL,
  "nonce" character varying NOT NULL DEFAULT '',
  "response_type" character varying NOT NULL,
  "user_id" bigint NULL REFERENCES "users" ("id") ON DELETE CASCADE,
  "done" boolean NOT NULL DEFAULT false,
  "auth_time" timestamptz NULL,
  "code_hash" character varying NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "index_oauth_authorization_requests_on_code_hash" ON "oauth_authorization_requests" ("code_hash");
CREATE INDEX "index_oauth_authorization_requests_on_expires_at" ON "oauth_authorization_requests" ("expires_at");

-- An issued access token. The bearer value is the encrypted pair (id, subject);
-- only the id is stored, and the resource server verifies by looking it up.
CREATE TABLE "oauth_access_tokens" (
  "id" character varying NOT NULL,
  "client_id" character varying NOT NULL REFERENCES "oauth_clients" ("client_id") ON DELETE CASCADE,
  "user_id" bigint NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
  "scopes" character varying[] NOT NULL DEFAULT '{}',
  "refresh_token_id" character varying NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE INDEX "index_oauth_access_tokens_on_user_id" ON "oauth_access_tokens" ("user_id");
CREATE INDEX "index_oauth_access_tokens_on_refresh_token_id" ON "oauth_access_tokens" ("refresh_token_id");

-- A refresh token the client holds, kept as its SHA-256 hash.
CREATE TABLE "oauth_refresh_tokens" (
  "id" character varying NOT NULL,
  "token_hash" character varying NOT NULL,
  "client_id" character varying NOT NULL REFERENCES "oauth_clients" ("client_id") ON DELETE CASCADE,
  "user_id" bigint NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
  "scopes" character varying[] NOT NULL DEFAULT '{}',
  "auth_time" timestamptz NOT NULL,
  "access_token_id" character varying NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "index_oauth_refresh_tokens_on_token_hash" ON "oauth_refresh_tokens" ("token_hash");
CREATE INDEX "index_oauth_refresh_tokens_on_user_id" ON "oauth_refresh_tokens" ("user_id");
