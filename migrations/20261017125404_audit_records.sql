-- ADR-0017: the audit trail. One row per operation an administrator or staff
-- member performs, written in the same transaction as the change.
--
-- No foreign keys on purpose: a record outlives the user, the Connected App and
-- the entity it names. The primary key carries occurred_at so the table can be
-- turned into a monthly-partitioned one later (ticket #846) without a rewrite;
-- PostgreSQL requires the partition key inside every unique constraint.
CREATE TABLE "audit_records" (
  "id" bigint GENERATED ALWAYS AS IDENTITY,
  "occurred_at" timestamptz NOT NULL DEFAULT now(),
  "actor_user_id" bigint NULL,
  "channel" character varying NOT NULL,
  "connected_app_id" bigint NULL,
  "token_id" bigint NULL,
  "request_id" character varying NULL,
  "address" character varying NULL,
  "user_agent" character varying NULL,
  "operation" character varying NOT NULL,
  "entity_type" character varying NULL,
  "entity_id" character varying NULL,
  "locale" character varying NULL,
  "outcome" character varying NOT NULL,
  "changes" jsonb NULL,
  PRIMARY KEY ("id", "occurred_at"),
  CONSTRAINT "audit_records_channel_check" CHECK ("channel" IN ('web', 'mcp')),
  CONSTRAINT "audit_records_outcome_check" CHECK ("outcome" IN ('succeeded', 'forbidden', 'failed'))
);
CREATE INDEX "index_audit_records_on_occurred_at" ON "audit_records" ("occurred_at");
CREATE INDEX "index_audit_records_on_actor_user_id_and_occurred_at" ON "audit_records" ("actor_user_id", "occurred_at");
CREATE INDEX "index_audit_records_on_entity" ON "audit_records" ("entity_type", "entity_id", "occurred_at");

-- Append-only, enforced by the database so that no code path (a handler, a
-- console session, a future bug) can rewrite history. A row-level trigger is
-- cloned onto every partition of a partitioned table, which a statement-level
-- one is not. Retention drops whole partitions (DROP TABLE does not fire it).
CREATE FUNCTION "audit_records_refuse_change"() RETURNS trigger
  LANGUAGE plpgsql AS $$
BEGIN
  RAISE EXCEPTION 'audit_records is append-only: % refused', TG_OP
    USING ERRCODE = 'restrict_violation';
END;
$$;
CREATE TRIGGER "audit_records_append_only"
  BEFORE UPDATE OR DELETE ON "audit_records"
  FOR EACH ROW EXECUTE FUNCTION "audit_records_refuse_change"();
