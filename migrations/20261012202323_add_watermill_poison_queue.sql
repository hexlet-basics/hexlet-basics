-- Watermill's poison queue for domain events (ADR-0004): a message whose handler
-- still fails after the retries is moved here instead of being redelivered
-- forever. Same shape as watermill_domain_events; only a publisher writes it,
-- so there is no offsets table.
CREATE TABLE "watermill_domain_events_poison" (
  "offset" bigserial,
  "uuid" character varying(36) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "payload" json NULL,
  "metadata" json NULL,
  "transaction_id" xid8 NOT NULL,
  PRIMARY KEY ("transaction_id", "offset")
);
