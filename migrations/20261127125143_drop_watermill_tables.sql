-- Domain events no longer go through a Watermill outbox (ADR-0031): publishing a
-- fact inserts the River jobs of its consumers in the caller's transaction. The
-- topic, its consumer offsets and its poison queue have no reader or writer left,
-- and the Go stack has not served production, so no row is carried over.

SET lock_timeout = '5s';
SET statement_timeout = '60s';

-- squawk-ignore ban-drop-table
DROP TABLE "watermill_offsets_domain_events";
-- squawk-ignore ban-drop-table
DROP TABLE "watermill_domain_events_poison";
-- squawk-ignore ban-drop-table
DROP TABLE "watermill_domain_events";
