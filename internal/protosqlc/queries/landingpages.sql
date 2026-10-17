-- name: CatalogPages :many
-- Listed landing pages of a locale with the Course and its current Version,
-- in catalog order: Course display order (NULLS LAST), Course id, page id.
SELECT sqlc.embed(p), sqlc.embed(c)
FROM language_landing_pages p
JOIN languages c ON c.id = p.language_id
WHERE p.state = 'published'
  AND p.locale = @locale::text
  AND c.progress = 'completed'
  AND p.listed = true
  AND (sqlc.narg('category_id')::bigint IS NULL OR p.language_category_id = sqlc.narg('category_id'))
ORDER BY c."order" NULLS LAST, p.language_id, p.id;

-- name: CatalogMainPages :many
SELECT sqlc.embed(p), sqlc.embed(c)
FROM language_landing_pages p
JOIN languages c ON c.id = p.language_id
WHERE p.state = 'published'
  AND p.locale = @locale::text
  AND c.progress = 'completed'
  AND p.main = true
  AND p.language_id = ANY(@course_ids::bigint[])
ORDER BY p.id;

-- name: VersionsByIDs :many
-- The eager-load step ent does itself (WithCurrentVersion = one IN query). A
-- LEFT JOIN + sqlc.embed scans NULL into a non-pointer struct and fails for a
-- Course with no current Version, so the second statement is written by hand.
SELECT * FROM language_versions WHERE id = ANY(@ids::bigint[]);
