-- Server-side sessions (ADR-0017). Not `sessions`: the Rails baseline already has a
-- table of that name from Rails 8's own authentication, which stays untouched. `token` holds the digest of the cookie value,
-- never the value itself. `data` is the session codec's opaque encoding.
CREATE TABLE "user_sessions" (
  "token" text NOT NULL,
  "data" bytea NOT NULL,
  "expiry" timestamp NOT NULL,
  PRIMARY KEY ("token")
);
-- Expired sessions are removed by a periodic job that selects on this column.
CREATE INDEX "index_user_sessions_on_expiry" ON "user_sessions" ("expiry");

-- The moment a password last changed. A session remembers the value it started
-- under, so changing or resetting a password ends the user's other sessions.
ALTER TABLE "users" ADD COLUMN "password_changed_at" timestamp NULL;
