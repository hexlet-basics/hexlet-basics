-- Emails are stored lowercase, as legacy's `normalizes :email` did, and the Go
-- stack now folds them on write and on every lookup. Rows written before that
-- (or by hand) may still be mixed case and would no longer be found.
-- The unique index on lower(users.email) already forbids two users differing
-- only by case, so this cannot collide.
UPDATE "users" SET "email" = lower("email") WHERE "email" <> lower("email");
UPDATE "leads" SET "email" = lower("email") WHERE "email" <> lower("email");
