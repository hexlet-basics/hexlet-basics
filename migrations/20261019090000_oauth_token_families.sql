-- ADR-0022, refresh rotation (#839). Every access and refresh token descends
-- from one authorization: its family. Presenting a refresh token that was
-- already used revokes the whole family, so the server must remember a used
-- refresh token (used_at) instead of deleting it, until its own expiry.
ALTER TABLE "oauth_refresh_tokens"
  ADD COLUMN "family_id" character varying NULL,
  ADD COLUMN "used_at" timestamptz NULL;
ALTER TABLE "oauth_access_tokens" ADD COLUMN "family_id" character varying NULL;

-- Tokens issued before families existed each become a family of their own: the
-- pair of a refresh token shares its id, an access token without a refresh
-- token stands alone.
UPDATE "oauth_refresh_tokens" SET "family_id" = "id";
UPDATE "oauth_access_tokens" SET "family_id" = COALESCE("refresh_token_id", "id");

ALTER TABLE "oauth_refresh_tokens" ALTER COLUMN "family_id" SET NOT NULL;
ALTER TABLE "oauth_access_tokens" ALTER COLUMN "family_id" SET NOT NULL;

CREATE INDEX "index_oauth_refresh_tokens_on_family_id" ON "oauth_refresh_tokens" ("family_id");
CREATE INDEX "index_oauth_access_tokens_on_family_id" ON "oauth_access_tokens" ("family_id");
-- The periodic cleanup deletes by expiry.
CREATE INDEX "index_oauth_refresh_tokens_on_expires_at" ON "oauth_refresh_tokens" ("expires_at");
CREATE INDEX "index_oauth_access_tokens_on_expires_at" ON "oauth_access_tokens" ("expires_at");
