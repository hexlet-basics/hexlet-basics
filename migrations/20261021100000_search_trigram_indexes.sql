-- The admin lists take a text filter `q` (a case-insensitive substring match,
-- ILIKE '%term%'). A btree cannot serve a leading wildcard, so the tables that
-- grow without bound get a trigram GIN index on the columns the filter reads;
-- the small editorial tables (banners, blog posts, courses, landing pages,
-- reviews) scan, which costs nothing at their size.
--
-- pg_trgm ships with PostgreSQL and is a trusted extension, so a role that may
-- create objects in the database may create it.
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Users: the first name, the last name and the email (the typeahead and the list).
CREATE INDEX "index_users_on_first_name_trgm" ON "users" USING gin ("first_name" gin_trgm_ops);
CREATE INDEX "index_users_on_last_name_trgm" ON "users" USING gin ("last_name" gin_trgm_ops);
CREATE INDEX "index_users_on_email_trgm" ON "users" USING gin ("email" gin_trgm_ops);

-- Leads: how an editor finds a person's request.
CREATE INDEX "index_leads_on_email_trgm" ON "leads" USING gin ("email" gin_trgm_ops);
CREATE INDEX "index_leads_on_phone_trgm" ON "leads" USING gin ("phone" gin_trgm_ops);
CREATE INDEX "index_leads_on_whatsapp_trgm" ON "leads" USING gin ("whatsapp" gin_trgm_ops);
CREATE INDEX "index_leads_on_telegram_trgm" ON "leads" USING gin ("telegram" gin_trgm_ops);

-- Assistant messages: one row per message of every learner's chat.
CREATE INDEX "index_ai_messages_on_content_trgm" ON "ai_messages" USING gin ("content" gin_trgm_ops);

-- Lesson translations: a row per lesson, locale and course version.
CREATE INDEX "index_language_lesson_version_infos_on_name_trgm" ON "language_lesson_version_infos" USING gin ("name" gin_trgm_ops);
CREATE INDEX "index_language_lesson_version_infos_on_description_trgm" ON "language_lesson_version_infos" USING gin ("description" gin_trgm_ops);
