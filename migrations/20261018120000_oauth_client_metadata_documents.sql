-- A client created from a Client ID Metadata Document (ADR-0019, #847) keeps the
-- instant its cached copy of the document expires; the next authorization
-- request after that fetches it again. Null for every other client. timestamptz
-- because it is compared with "now".
ALTER TABLE "oauth_clients" ADD COLUMN "metadata_expires_at" timestamptz NULL;
