-- Enforce the legacy Language::Category uniqueness validations as real DB
-- constraints (previously checked in the handler with per-field Exist queries,
-- which raced under concurrency). Parity with legacy/app/models/language/category.rb:
--   validates :name,   uniqueness: true
--   validates :header, uniqueness: true
--   validates :slug,   uniqueness: { scope: :locale }
--
-- The handler no longer pre-checks: a violation surfaces as ent.IsConstraintError
-- and the central ogen ErrorHandler maps it to 409.
--
-- PRE-DEPLOY: these indexes are added to a table holding real legacy data, and
-- the old app-level check could race, so prod may contain duplicates that will
-- make CREATE UNIQUE INDEX fail. Before shipping, verify (and dedup) on prod:
--   SELECT name,   count(*) FROM language_categories GROUP BY name   HAVING count(*) > 1;
--   SELECT header, count(*) FROM language_categories GROUP BY header HAVING count(*) > 1;
--   SELECT slug, locale, count(*) FROM language_categories GROUP BY slug, locale HAVING count(*) > 1;

-- Global uniqueness for name and header. Both columns are presence-validated at
-- the API (@minLength(1)), so the "multiple NULLs allowed" default is moot.
CREATE UNIQUE INDEX "language_categories_name_uniq" ON "language_categories" ("name");
CREATE UNIQUE INDEX "language_categories_header_uniq" ON "language_categories" ("header");

-- slug is unique per locale. NULLS NOT DISTINCT (Postgres 15+) makes two rows
-- with the same slug and a NULL locale collide — matching Rails' `scope: :locale`,
-- which compares a nil locale with `locale IS NULL` rather than treating NULLs as
-- distinct. Without it, null-locale rows (the shape the current input creates)
-- could duplicate slugs freely.
CREATE UNIQUE INDEX "language_categories_slug_locale_uniq" ON "language_categories" ("slug", "locale") NULLS NOT DISTINCT;
