-- ADR-0019 addendum (#854): an MCP call records the access token it came with,
-- and an OAuth access token's id is a random string (oauth_access_tokens.id),
-- not a number. The column was reserved as bigint before the token existed; no
-- row has ever filled it, so the cast keeps nothing but null. ALTER on the
-- partitioned parent applies to every partition; the column is not the key.
ALTER TABLE "audit_records" ALTER COLUMN "token_id" TYPE character varying USING "token_id"::character varying;
