-- Watermill's durable PostgreSQL topic for domain events (ADR-0004).
-- The shape mirrors watermill-sql v4.1.5 DefaultPostgreSQLSchema; Atlas owns
-- creation, so Watermill runtime schema initialization stays disabled.
CREATE TABLE "watermill_domain_events" (
  "offset" bigserial,
  "uuid" character varying(36) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "payload" json NULL,
  "metadata" json NULL,
  "transaction_id" xid8 NOT NULL,
  PRIMARY KEY ("transaction_id", "offset")
);

-- Every subscriber has its own stable consumer_group row. Independent offsets
-- are what make one stored domain fact fan out durably to multiple handlers.
CREATE TABLE "watermill_offsets_domain_events" (
  "consumer_group" character varying(255) NOT NULL,
  "offset_acked" bigint NULL,
  "last_processed_transaction_id" xid8 NOT NULL,
  PRIMARY KEY ("consumer_group")
);
