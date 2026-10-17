-- name: ListUsers :many
SELECT * FROM users ORDER BY id DESC LIMIT @lim OFFSET @off;

-- name: CountUsers :one
SELECT count(*) FROM users;

-- name: SearchUsers :many
-- Typeahead: case-insensitive substring on first name, last name or email.
SELECT * FROM users
WHERE first_name ILIKE '%' || @term::text || '%'
   OR last_name  ILIKE '%' || @term::text || '%'
   OR email      ILIKE '%' || @term::text || '%'
ORDER BY id DESC LIMIT 20;

-- name: FilterUsers :many
-- Dynamic admin filter: every argument optional (sqlc.narg), so one static
-- statement covers every combination. Contrast with ent's composable predicates.
SELECT * FROM users
WHERE (sqlc.narg('term')::text IS NULL OR first_name ILIKE '%' || sqlc.narg('term') || '%' OR last_name ILIKE '%' || sqlc.narg('term') || '%' OR email ILIKE '%' || sqlc.narg('term') || '%')
  AND (sqlc.narg('locale')::text IS NULL OR locale = sqlc.narg('locale'))
  AND (sqlc.narg('admin')::bool IS NULL OR admin = sqlc.narg('admin'))
ORDER BY id DESC LIMIT @lim OFFSET @off;
