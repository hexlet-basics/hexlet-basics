-- Create the Go stack's own blob-record table (ADR-0005), backing the ent
-- Attachment schema. One row per uploaded file: `storage_key` locates the bytes
-- in the gocloud bucket and is unique (it is both the write key and the read
-- URL). This is a NEW table (no legacy `attachments`); ActiveStorage's
-- active_storage_* tables are left untouched — migrating their blobs into this
-- table is downstream (consumer FK columns land with BlogPosts/covers).
CREATE TABLE "attachments" (
  "id" bigserial NOT NULL,
  "storage_key" character varying NOT NULL,
  "filename" character varying NOT NULL,
  "content_type" character varying NOT NULL,
  "byte_size" bigint NOT NULL,
  "created_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "attachments_storage_key_uniq" ON "attachments" ("storage_key");
