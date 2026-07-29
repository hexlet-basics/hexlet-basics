-- A module is a stable course-scoped identity reused by every course version.
-- Enforce the identity key in PostgreSQL so concurrent builds cannot create two
-- rows before either transaction observes the other.
--
-- Existing duplicates must be checked before this migration reaches production:
-- SELECT language_id, slug, count(*)
-- FROM language_modules
-- WHERE language_id IS NOT NULL AND slug IS NOT NULL
-- GROUP BY language_id, slug
-- HAVING count(*) > 1;
CREATE UNIQUE INDEX "language_modules_language_id_slug_uniq"
ON "language_modules" ("language_id", "slug");
