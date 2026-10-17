-- ADR-0019 / ticket #850: a read of personal data is recorded with what it
-- touched (a row count and, up to a cap, the ids), never with field values.
-- `changes` is the old and new value of each changed field and means a change,
-- so a read gets a column of its own. Adding a column to the partitioned parent
-- adds it to every partition, the default one included.
ALTER TABLE "audit_records" ADD COLUMN "reads" jsonb NULL;
