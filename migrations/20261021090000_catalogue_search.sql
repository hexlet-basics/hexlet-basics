-- Catalogue search (#877): full-text search over the author-controlled text of
-- the public catalogue, which is course landing pages, lessons and blog posts.
-- Nothing a learner wrote is indexed, and no table of learner content is touched.
--
-- One function builds the document of every kind of record, so the indexes below
-- and the queries of internal/search name the same expression and cannot drift
-- apart: a title counts most (weight A), a summary next (B), the body least (C).
-- It is IMMUTABLE, which an index expression requires, because to_tsvector with an
-- explicit configuration is.
CREATE FUNCTION "catalogue_search_vector"("config" regconfig, "title" text, "summary" text, "body" text)
RETURNS tsvector
LANGUAGE sql IMMUTABLE PARALLEL SAFE
AS $$
  SELECT setweight(to_tsvector("config", coalesce("title", '')), 'A')
      || setweight(to_tsvector("config", coalesce("summary", '')), 'B')
      || setweight(to_tsvector("config", coalesce("body", '')), 'C')
$$;

-- One index per configuration the site speaks (ru and en). The query picks the
-- configuration of the request locale and writes it as a literal, so the planner
-- can match the expression; a locale with no index here (es) falls back to a scan,
-- which is cheap while the catalogue has no content in it.
CREATE INDEX "index_language_landing_pages_on_search_ru" ON "language_landing_pages"
  USING gin (catalogue_search_vector('russian', coalesce("name", '') || ' ' || coalesce("header", ''), "description", "outcomes_description"));
CREATE INDEX "index_language_landing_pages_on_search_en" ON "language_landing_pages"
  USING gin (catalogue_search_vector('english', coalesce("name", '') || ' ' || coalesce("header", ''), "description", "outcomes_description"));

CREATE INDEX "index_language_lesson_version_infos_on_search_ru" ON "language_lesson_version_infos"
  USING gin (catalogue_search_vector('russian', "name", "description", coalesce("theory", '') || ' ' || coalesce("instructions", '')));
CREATE INDEX "index_language_lesson_version_infos_on_search_en" ON "language_lesson_version_infos"
  USING gin (catalogue_search_vector('english', "name", "description", coalesce("theory", '') || ' ' || coalesce("instructions", '')));

CREATE INDEX "index_blog_posts_on_search_ru" ON "blog_posts"
  USING gin (catalogue_search_vector('russian', "name", "description", "rich_body"));
CREATE INDEX "index_blog_posts_on_search_en" ON "blog_posts"
  USING gin (catalogue_search_vector('english', "name", "description", "rich_body"));
