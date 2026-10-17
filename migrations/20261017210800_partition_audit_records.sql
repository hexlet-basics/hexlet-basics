-- ADR-0020 / ticket #846: audit_records becomes a table partitioned by month on
-- occurred_at, so retention is a DROP of a whole month instead of a DELETE
-- (which the append-only trigger refuses, and which would bloat the table).
--
-- The three properties the original migration set up for this stay: the primary
-- key is (id, occurred_at), there are no foreign keys, and the append-only
-- trigger is ROW-level, so PostgreSQL clones it onto every partition. The
-- application's insert path is unchanged: it inserts into audit_records and
-- PostgreSQL routes the row to its month.

-- 1. Set the unpartitioned table aside. Index, sequence and constraint names
-- are schema-wide, so they are renamed to free the names for the new table.
ALTER TABLE "audit_records" RENAME TO "audit_records_unpartitioned";
ALTER SEQUENCE "audit_records_id_seq" RENAME TO "audit_records_unpartitioned_id_seq";
ALTER INDEX "audit_records_pkey" RENAME TO "audit_records_unpartitioned_pkey";
ALTER INDEX "index_audit_records_on_occurred_at" RENAME TO "audit_records_unpartitioned_occurred_at_idx";
ALTER INDEX "index_audit_records_on_actor_user_id_and_occurred_at" RENAME TO "audit_records_unpartitioned_actor_idx";
ALTER INDEX "index_audit_records_on_entity" RENAME TO "audit_records_unpartitioned_entity_idx";
ALTER TABLE "audit_records_unpartitioned" RENAME CONSTRAINT "audit_records_channel_check" TO "audit_records_unpartitioned_channel_check";
ALTER TABLE "audit_records_unpartitioned" RENAME CONSTRAINT "audit_records_outcome_check" TO "audit_records_unpartitioned_outcome_check";

-- 2. The partitioned parent: same columns, key, checks, indexes and trigger.
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
) PARTITION BY RANGE ("occurred_at");
CREATE INDEX "index_audit_records_on_occurred_at" ON "audit_records" ("occurred_at");
CREATE INDEX "index_audit_records_on_actor_user_id_and_occurred_at" ON "audit_records" ("actor_user_id", "occurred_at");
CREATE INDEX "index_audit_records_on_entity" ON "audit_records" ("entity_type", "entity_id", "occurred_at");

CREATE TRIGGER "audit_records_append_only"
  BEFORE UPDATE OR DELETE ON "audit_records"
  FOR EACH ROW EXECUTE FUNCTION "audit_records_refuse_change"();

-- The safety net: a row whose month has no partition (the maintenance job has
-- not run, or a record carries an odd timestamp) lands here instead of failing
-- the audited operation. audit_records_ensure_partition moves such rows out
-- when it later creates their month.
CREATE TABLE "audit_records_default" PARTITION OF "audit_records" DEFAULT;

-- 3. Partition DDL lives in SQL functions, not in Go strings: the application
-- calls two fixed statements and the DDL is written, reviewed and tested once.
-- Months are calendar months in UTC. Partitions are named audit_records_pYYYYMM.

-- Creates the partition for the month containing p_month. Idempotent and safe
-- to run from concurrent workers (advisory lock). Rows already parked in the
-- default partition for that month are moved into the new partition, because
-- PostgreSQL refuses to create a partition that overlaps default rows. The
-- move is the one place rows are deleted, and it is possible only because
-- DETACH removes the append-only trigger that was cloned onto the default
-- partition (ATTACH clones it back). It all happens inside this function's
-- transaction and the rows are re-inserted unchanged, so no record is lost or
-- altered, and the application cannot reach it: it has no such statement.
CREATE FUNCTION "audit_records_ensure_partition"(p_month timestamptz) RETURNS text
  LANGUAGE plpgsql AS $$
DECLARE
  month_start timestamptz := date_trunc('month', p_month AT TIME ZONE 'UTC') AT TIME ZONE 'UTC';
  month_end timestamptz := month_start + interval '1 month';
  part_name text := 'audit_records_p' || to_char(month_start AT TIME ZONE 'UTC', 'YYYYMM');
  bounds text;
BEGIN
  PERFORM pg_advisory_xact_lock(hashtext('audit_records_partitions'));
  IF to_regclass(quote_ident(part_name)) IS NOT NULL THEN
    RETURN part_name;
  END IF;
  bounds := format('FOR VALUES FROM (%L) TO (%L)',
    to_char(month_start AT TIME ZONE 'UTC', 'YYYY-MM-DD HH24:MI:SS+00'),
    to_char(month_end AT TIME ZONE 'UTC', 'YYYY-MM-DD HH24:MI:SS+00'));

  IF EXISTS (SELECT 1 FROM "audit_records_default"
             WHERE "occurred_at" >= month_start AND "occurred_at" < month_end) THEN
    ALTER TABLE "audit_records" DETACH PARTITION "audit_records_default";
    EXECUTE format('CREATE TABLE %I PARTITION OF "audit_records" %s', part_name, bounds);
    WITH moved AS (
      DELETE FROM "audit_records_default"
      WHERE "occurred_at" >= month_start AND "occurred_at" < month_end
      RETURNING *
    )
    INSERT INTO "audit_records" OVERRIDING SYSTEM VALUE SELECT * FROM moved;
    ALTER TABLE "audit_records" ATTACH PARTITION "audit_records_default" DEFAULT;
  ELSE
    EXECUTE format('CREATE TABLE %I PARTITION OF "audit_records" %s', part_name, bounds);
  END IF;
  RETURN part_name;
END;
$$;

-- Drops every monthly partition that ended at or before p_cutoff and returns
-- the names dropped. A partition is dropped whole or not at all: the month that
-- contains the cutoff, and the default partition, are never touched, so a
-- record is kept for at least the retention and up to a month longer. DROP
-- TABLE does not fire the append-only trigger; this is the only way a record
-- leaves the table.
CREATE FUNCTION "audit_records_drop_partitions_ended_before"(p_cutoff timestamptz) RETURNS SETOF text
  LANGUAGE plpgsql AS $$
DECLARE
  part record;
  month_start timestamptz;
BEGIN
  PERFORM pg_advisory_xact_lock(hashtext('audit_records_partitions'));
  FOR part IN
    SELECT c.relname::text AS name
    FROM pg_inherits i
    JOIN pg_class c ON c.oid = i.inhrelid
    WHERE i.inhparent = 'audit_records'::regclass
      AND c.relname ~ '^audit_records_p[0-9]{6}$'
    ORDER BY c.relname
  LOOP
    month_start := to_timestamp(substr(part.name, 16) || '01', 'YYYYMMDD') AT TIME ZONE 'UTC';
    IF month_start + interval '1 month' <= p_cutoff THEN
      EXECUTE format('DROP TABLE %I', part.name);
      RETURN NEXT part.name;
    END IF;
  END LOOP;
END;
$$;

-- 4. Carry over what the unpartitioned table holds: partitions for every month
-- it has rows in and for the current one, then the rows themselves.
SELECT "audit_records_ensure_partition"(m)
FROM generate_series(
  date_trunc('month', COALESCE((SELECT min("occurred_at") FROM "audit_records_unpartitioned"), now())),
  date_trunc('month', now()) + interval '3 months',
  interval '1 month'
) AS m;
INSERT INTO "audit_records" OVERRIDING SYSTEM VALUE SELECT * FROM "audit_records_unpartitioned";
SELECT setval(pg_get_serial_sequence('audit_records', 'id'),
  COALESCE((SELECT max("id") FROM "audit_records_unpartitioned"), 0) + 1, false);
DROP TABLE "audit_records_unpartitioned";
