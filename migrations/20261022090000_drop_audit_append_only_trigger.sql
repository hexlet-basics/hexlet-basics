-- No trigger of ours: nothing in the application updates or deletes an Audit
-- Record, and a rule that lives in a trigger is logic hidden from the code.
-- Retention still removes whole partitions (DROP TABLE), and
-- audit_records_ensure_partition still moves default-partition rows, neither of
-- which depended on the trigger beyond having to get around it.
DROP TRIGGER "audit_records_append_only" ON "audit_records";
DROP FUNCTION "audit_records_refuse_change"();
