-- name: InsertEnrollmentIfAbsent :one
-- Upsert as ent's OnConflictColumns(...).DoNothing().ID(): no row back on conflict.
INSERT INTO language_members (user_id, language_id, state, created_at, updated_at)
VALUES (@user_id, @course_id, @state, now(), now())
ON CONFLICT (user_id, language_id) DO NOTHING
RETURNING id;
