-- Create "action_text_rich_texts" table
CREATE TABLE "action_text_rich_texts" (
  "id" bigserial NOT NULL,
  "body" text NULL,
  "created_at" timestamp NOT NULL,
  "name" character varying NOT NULL,
  "record_id" bigint NOT NULL,
  "record_type" character varying NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_action_text_rich_texts_uniqueness" to table: "action_text_rich_texts"
CREATE UNIQUE INDEX "index_action_text_rich_texts_uniqueness" ON "action_text_rich_texts" ("record_type", "record_id", "name");
-- Create "active_storage_attachments" table
CREATE TABLE "active_storage_attachments" (
  "id" bigserial NOT NULL,
  "blob_id" bigint NOT NULL,
  "created_at" timestamp NOT NULL,
  "name" character varying NOT NULL,
  "record_id" bigint NOT NULL,
  "record_type" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_active_storage_attachments_on_blob_id" to table: "active_storage_attachments"
CREATE INDEX "index_active_storage_attachments_on_blob_id" ON "active_storage_attachments" ("blob_id");
-- Create index "index_active_storage_attachments_uniqueness" to table: "active_storage_attachments"
CREATE UNIQUE INDEX "index_active_storage_attachments_uniqueness" ON "active_storage_attachments" ("record_type", "record_id", "name", "blob_id");
-- Create "active_storage_blobs" table
CREATE TABLE "active_storage_blobs" (
  "id" bigserial NOT NULL,
  "byte_size" bigint NOT NULL,
  "checksum" character varying NULL,
  "content_type" character varying NULL,
  "created_at" timestamp NOT NULL,
  "filename" character varying NOT NULL,
  "key" character varying NOT NULL,
  "metadata" text NULL,
  "service_name" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_active_storage_blobs_on_key" to table: "active_storage_blobs"
CREATE UNIQUE INDEX "index_active_storage_blobs_on_key" ON "active_storage_blobs" ("key");
-- Create "active_storage_variant_records" table
CREATE TABLE "active_storage_variant_records" (
  "id" bigserial NOT NULL,
  "blob_id" bigint NOT NULL,
  "variation_digest" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_active_storage_variant_records_uniqueness" to table: "active_storage_variant_records"
CREATE UNIQUE INDEX "index_active_storage_variant_records_uniqueness" ON "active_storage_variant_records" ("blob_id", "variation_digest");
-- Create "ahoy_events" table
CREATE TABLE "ahoy_events" (
  "id" bigserial NOT NULL,
  "name" character varying NULL,
  "properties" jsonb NULL,
  "time" timestamp NULL,
  "user_id" bigint NULL,
  "visit_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_ahoy_events_on_name_and_time" to table: "ahoy_events"
CREATE INDEX "index_ahoy_events_on_name_and_time" ON "ahoy_events" ("name", "time");
-- Create index "index_ahoy_events_on_properties" to table: "ahoy_events"
CREATE INDEX "index_ahoy_events_on_properties" ON "ahoy_events" USING GIN ("properties" jsonb_path_ops);
-- Create index "index_ahoy_events_on_user_id" to table: "ahoy_events"
CREATE INDEX "index_ahoy_events_on_user_id" ON "ahoy_events" ("user_id");
-- Create index "index_ahoy_events_on_visit_id" to table: "ahoy_events"
CREATE INDEX "index_ahoy_events_on_visit_id" ON "ahoy_events" ("visit_id");
-- Create "ahoy_visits" table
CREATE TABLE "ahoy_visits" (
  "id" bigserial NOT NULL,
  "app_version" character varying NULL,
  "browser" character varying NULL,
  "city" character varying NULL,
  "country" character varying NULL,
  "device_type" character varying NULL,
  "ip" character varying NULL,
  "landing_page" text NULL,
  "latitude" double precision NULL,
  "longitude" double precision NULL,
  "os" character varying NULL,
  "os_version" character varying NULL,
  "platform" character varying NULL,
  "referrer" text NULL,
  "referring_domain" character varying NULL,
  "region" character varying NULL,
  "started_at" timestamp NULL,
  "user_agent" text NULL,
  "user_id" bigint NULL,
  "utm_campaign" character varying NULL,
  "utm_content" character varying NULL,
  "utm_medium" character varying NULL,
  "utm_source" character varying NULL,
  "utm_term" character varying NULL,
  "visit_token" character varying NULL,
  "visitor_token" character varying NULL,
  "ym_client_id" character varying NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_ahoy_visits_on_user_id" to table: "ahoy_visits"
CREATE INDEX "index_ahoy_visits_on_user_id" ON "ahoy_visits" ("user_id");
-- Create index "index_ahoy_visits_on_visit_token" to table: "ahoy_visits"
CREATE UNIQUE INDEX "index_ahoy_visits_on_visit_token" ON "ahoy_visits" ("visit_token");
-- Create index "index_ahoy_visits_on_visitor_token_and_started_at" to table: "ahoy_visits"
CREATE INDEX "index_ahoy_visits_on_visitor_token_and_started_at" ON "ahoy_visits" ("visitor_token", "started_at");
-- Create "ai_chats" table
CREATE TABLE "ai_chats" (
  "id" bigserial NOT NULL,
  "ai_model_id" bigint NULL,
  "created_at" timestamp NOT NULL,
  "language_lesson_member_id" bigint NOT NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_ai_chats_on_ai_model_id" to table: "ai_chats"
CREATE INDEX "index_ai_chats_on_ai_model_id" ON "ai_chats" ("ai_model_id");
-- Create index "index_ai_chats_on_language_lesson_member_id" to table: "ai_chats"
CREATE INDEX "index_ai_chats_on_language_lesson_member_id" ON "ai_chats" ("language_lesson_member_id");
-- Create index "index_ai_chats_on_user_id_and_language_lesson_member_id" to table: "ai_chats"
CREATE UNIQUE INDEX "index_ai_chats_on_user_id_and_language_lesson_member_id" ON "ai_chats" ("user_id", "language_lesson_member_id");
-- Create "ai_messages" table
CREATE TABLE "ai_messages" (
  "id" bigserial NOT NULL,
  "ai_chat_id" bigint NOT NULL,
  "ai_model_id" bigint NULL,
  "ai_tool_call_id" bigint NULL,
  "cache_creation_tokens" integer NULL,
  "cached_tokens" integer NULL,
  "content" text NULL,
  "content_raw" json NULL,
  "created_at" timestamp NOT NULL,
  "input_tokens" integer NULL,
  "output_tokens" integer NULL,
  "role" character varying NOT NULL,
  "thinking_signature" text NULL,
  "thinking_text" text NULL,
  "thinking_tokens" integer NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_ai_messages_on_ai_chat_id" to table: "ai_messages"
CREATE INDEX "index_ai_messages_on_ai_chat_id" ON "ai_messages" ("ai_chat_id");
-- Create index "index_ai_messages_on_ai_model_id" to table: "ai_messages"
CREATE INDEX "index_ai_messages_on_ai_model_id" ON "ai_messages" ("ai_model_id");
-- Create index "index_ai_messages_on_ai_tool_call_id" to table: "ai_messages"
CREATE INDEX "index_ai_messages_on_ai_tool_call_id" ON "ai_messages" ("ai_tool_call_id");
-- Create index "index_ai_messages_on_role" to table: "ai_messages"
CREATE INDEX "index_ai_messages_on_role" ON "ai_messages" ("role");
-- Create index "index_ai_messages_on_user_id" to table: "ai_messages"
CREATE INDEX "index_ai_messages_on_user_id" ON "ai_messages" ("user_id");
-- Create "ai_models" table
CREATE TABLE "ai_models" (
  "id" bigserial NOT NULL,
  "capabilities" jsonb NULL DEFAULT '[]',
  "context_window" integer NULL,
  "created_at" timestamp NOT NULL,
  "family" character varying NULL,
  "knowledge_cutoff" date NULL,
  "max_output_tokens" integer NULL,
  "metadata" jsonb NULL DEFAULT '{}',
  "modalities" jsonb NULL DEFAULT '{}',
  "model_created_at" timestamp NULL,
  "model_id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "pricing" jsonb NULL DEFAULT '{}',
  "provider" character varying NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_ai_models_on_capabilities" to table: "ai_models"
CREATE INDEX "index_ai_models_on_capabilities" ON "ai_models" USING GIN ("capabilities");
-- Create index "index_ai_models_on_family" to table: "ai_models"
CREATE INDEX "index_ai_models_on_family" ON "ai_models" ("family");
-- Create index "index_ai_models_on_modalities" to table: "ai_models"
CREATE INDEX "index_ai_models_on_modalities" ON "ai_models" USING GIN ("modalities");
-- Create index "index_ai_models_on_provider_and_model_id" to table: "ai_models"
CREATE UNIQUE INDEX "index_ai_models_on_provider_and_model_id" ON "ai_models" ("provider", "model_id");
-- Create "ai_tool_calls" table
CREATE TABLE "ai_tool_calls" (
  "id" bigserial NOT NULL,
  "ai_message_id" bigint NOT NULL,
  "arguments" jsonb NULL DEFAULT '{}',
  "created_at" timestamp NOT NULL,
  "name" character varying NOT NULL,
  "thought_signature" text NULL,
  "tool_call_id" character varying NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_ai_tool_calls_on_ai_message_id" to table: "ai_tool_calls"
CREATE INDEX "index_ai_tool_calls_on_ai_message_id" ON "ai_tool_calls" ("ai_message_id");
-- Create index "index_ai_tool_calls_on_name" to table: "ai_tool_calls"
CREATE INDEX "index_ai_tool_calls_on_name" ON "ai_tool_calls" ("name");
-- Create index "index_ai_tool_calls_on_tool_call_id" to table: "ai_tool_calls"
CREATE UNIQUE INDEX "index_ai_tool_calls_on_tool_call_id" ON "ai_tool_calls" ("tool_call_id");
-- Create "ar_internal_metadata" table
CREATE TABLE "ar_internal_metadata" (
  "key" character varying NOT NULL,
  "value" character varying NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("key")
);
-- Create "banners" table
CREATE TABLE "banners" (
  "id" bigserial NOT NULL,
  "background" character varying NOT NULL DEFAULT 'cta_gradient',
  "body" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "finishes_at" timestamp NULL,
  "locale" character varying NOT NULL,
  "starts_at" timestamp NULL,
  "state" character varying NOT NULL DEFAULT 'draft',
  "updated_at" timestamp NOT NULL,
  "url" character varying NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_banners_on_locale_and_state" to table: "banners"
CREATE INDEX "index_banners_on_locale_and_state" ON "banners" ("locale", "state");
-- Create "blog_post_likes" table
CREATE TABLE "blog_post_likes" (
  "id" bigserial NOT NULL,
  "blog_post_id" bigint NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_blog_post_likes_on_blog_post_id" to table: "blog_post_likes"
CREATE INDEX "index_blog_post_likes_on_blog_post_id" ON "blog_post_likes" ("blog_post_id");
-- Create index "index_blog_post_likes_on_user_id" to table: "blog_post_likes"
CREATE INDEX "index_blog_post_likes_on_user_id" ON "blog_post_likes" ("user_id");
-- Create "blog_post_related_language_items" table
CREATE TABLE "blog_post_related_language_items" (
  "id" bigserial NOT NULL,
  "blog_post_id" bigint NOT NULL,
  "created_at" timestamp NOT NULL,
  "language_id" bigint NOT NULL,
  "order" integer NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_blog_post_related_language_items_on_blog_post_id" to table: "blog_post_related_language_items"
CREATE INDEX "index_blog_post_related_language_items_on_blog_post_id" ON "blog_post_related_language_items" ("blog_post_id");
-- Create index "index_blog_post_related_language_items_on_language_id" to table: "blog_post_related_language_items"
CREATE INDEX "index_blog_post_related_language_items_on_language_id" ON "blog_post_related_language_items" ("language_id");
-- Create "blog_posts" table
CREATE TABLE "blog_posts" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "creator_id" bigint NOT NULL,
  "description" character varying NULL,
  "language_id" bigint NULL,
  "locale" character varying NULL,
  "name" character varying NULL,
  "related_language_items_count" integer NOT NULL DEFAULT 0,
  "slug" character varying NULL,
  "state" character varying NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_blog_posts_on_creator_id" to table: "blog_posts"
CREATE INDEX "index_blog_posts_on_creator_id" ON "blog_posts" ("creator_id");
-- Create index "index_blog_posts_on_language_id" to table: "blog_posts"
CREATE INDEX "index_blog_posts_on_language_id" ON "blog_posts" ("language_id");
-- Create index "index_blog_posts_on_slug" to table: "blog_posts"
CREATE UNIQUE INDEX "index_blog_posts_on_slug" ON "blog_posts" ("slug");
-- Create "book_requests" table
CREATE TABLE "book_requests" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "state" character varying NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_book_requests_on_user_id" to table: "book_requests"
CREATE UNIQUE INDEX "index_book_requests_on_user_id" ON "book_requests" ("user_id");
-- Create "course_categories" table
CREATE TABLE "course_categories" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "name" character varying NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "event_store_events" table
CREATE TABLE "event_store_events" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "data" bytea NOT NULL,
  "event_id" character varying(36) NOT NULL,
  "event_type" character varying NOT NULL,
  "metadata" bytea NULL,
  "valid_at" timestamp NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_event_store_events_on_created_at" to table: "event_store_events"
CREATE INDEX "index_event_store_events_on_created_at" ON "event_store_events" ("created_at");
-- Create index "index_event_store_events_on_event_id" to table: "event_store_events"
CREATE UNIQUE INDEX "index_event_store_events_on_event_id" ON "event_store_events" ("event_id");
-- Create index "index_event_store_events_on_event_type" to table: "event_store_events"
CREATE INDEX "index_event_store_events_on_event_type" ON "event_store_events" ("event_type");
-- Create index "index_event_store_events_on_valid_at" to table: "event_store_events"
CREATE INDEX "index_event_store_events_on_valid_at" ON "event_store_events" ("valid_at");
-- Create "event_store_events_in_streams" table
CREATE TABLE "event_store_events_in_streams" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "event_id" character varying(36) NOT NULL,
  "position" integer NULL,
  "stream" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_event_store_events_in_streams_on_created_at" to table: "event_store_events_in_streams"
CREATE INDEX "index_event_store_events_in_streams_on_created_at" ON "event_store_events_in_streams" ("created_at");
-- Create index "index_event_store_events_in_streams_on_event_id" to table: "event_store_events_in_streams"
CREATE INDEX "index_event_store_events_in_streams_on_event_id" ON "event_store_events_in_streams" ("event_id");
-- Create index "index_event_store_events_in_streams_on_stream_and_event_id" to table: "event_store_events_in_streams"
CREATE UNIQUE INDEX "index_event_store_events_in_streams_on_stream_and_event_id" ON "event_store_events_in_streams" ("stream", "event_id");
-- Create index "index_event_store_events_in_streams_on_stream_and_position" to table: "event_store_events_in_streams"
CREATE UNIQUE INDEX "index_event_store_events_in_streams_on_stream_and_position" ON "event_store_events_in_streams" ("stream", "position");
-- Create "flipper_features" table
CREATE TABLE "flipper_features" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "key" character varying NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_flipper_features_on_key" to table: "flipper_features"
CREATE UNIQUE INDEX "index_flipper_features_on_key" ON "flipper_features" ("key");
-- Create "flipper_gates" table
CREATE TABLE "flipper_gates" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "feature_key" character varying NOT NULL,
  "key" character varying NOT NULL,
  "updated_at" timestamp NOT NULL,
  "value" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_flipper_gates_on_feature_key_and_key_and_value" to table: "flipper_gates"
CREATE UNIQUE INDEX "index_flipper_gates_on_feature_key_and_key_and_value" ON "flipper_gates" ("feature_key", "key", "value");
-- Create "language_categories" table
CREATE TABLE "language_categories" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "description" character varying NULL,
  "header" character varying NULL,
  "locale" character varying NULL,
  "name" character varying NULL,
  "name_en" character varying NULL,
  "name_ru" character varying NULL,
  "slug" character varying NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "language_category_items" table
CREATE TABLE "language_category_items" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "language_category_id" bigint NOT NULL,
  "language_landing_page_id" bigint NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_language_category_items_on_language_category_id" to table: "language_category_items"
CREATE INDEX "index_language_category_items_on_language_category_id" ON "language_category_items" ("language_category_id");
-- Create index "index_language_category_items_on_language_landing_page_id" to table: "language_category_items"
CREATE INDEX "index_language_category_items_on_language_landing_page_id" ON "language_category_items" ("language_landing_page_id");
-- Create "language_category_qna_items" table
CREATE TABLE "language_category_qna_items" (
  "id" bigserial NOT NULL,
  "answer" character varying NULL,
  "created_at" timestamp NOT NULL,
  "language_category_id" bigint NOT NULL,
  "question" character varying NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_language_category_qna_items_on_language_category_id" to table: "language_category_qna_items"
CREATE INDEX "index_language_category_qna_items_on_language_category_id" ON "language_category_qna_items" ("language_category_id");
-- Create "language_landing_page_qna_items" table
CREATE TABLE "language_landing_page_qna_items" (
  "id" bigserial NOT NULL,
  "answer" character varying NULL,
  "created_at" timestamp NOT NULL,
  "language_landing_page_id" bigint NOT NULL,
  "question" character varying NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_on_language_landing_page_id_98023e1f90" to table: "language_landing_page_qna_items"
CREATE INDEX "idx_on_language_landing_page_id_98023e1f90" ON "language_landing_page_qna_items" ("language_landing_page_id");
-- Create "language_landing_pages" table
CREATE TABLE "language_landing_pages" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "description" character varying NULL,
  "footer" boolean NULL,
  "footer_name" character varying NULL,
  "header" character varying NULL,
  "landing_page_to_redirect_id" bigint NULL,
  "language_category_id" bigint NULL,
  "language_id" bigint NOT NULL,
  "listed" boolean NULL,
  "locale" character varying NULL,
  "main" boolean NULL,
  "meta_description" character varying NULL,
  "meta_title" character varying NULL,
  "name" character varying NULL,
  "order" character varying NULL,
  "outcomes_description" character varying NULL,
  "outcomes_header" character varying NULL,
  "slug" character varying NULL,
  "state" character varying NULL,
  "updated_at" timestamp NOT NULL,
  "used_in_description" character varying NULL,
  "used_in_header" character varying NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_rails_00381fb5f4" FOREIGN KEY ("landing_page_to_redirect_id") REFERENCES "language_landing_pages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "index_language_landing_pages_on_landing_page_to_redirect_id" to table: "language_landing_pages"
CREATE INDEX "index_language_landing_pages_on_landing_page_to_redirect_id" ON "language_landing_pages" ("landing_page_to_redirect_id");
-- Create index "index_language_landing_pages_on_language_category_id" to table: "language_landing_pages"
CREATE INDEX "index_language_landing_pages_on_language_category_id" ON "language_landing_pages" ("language_category_id");
-- Create index "index_language_landing_pages_on_language_id" to table: "language_landing_pages"
CREATE INDEX "index_language_landing_pages_on_language_id" ON "language_landing_pages" ("language_id");
-- Create "language_lesson_members" table
CREATE TABLE "language_lesson_members" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "language_id" bigint NOT NULL,
  "language_member_id" bigint NOT NULL,
  "lesson_id" bigint NOT NULL,
  "messages_count" integer NULL DEFAULT 0,
  "state" character varying NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_language_lesson_members_on_language_member_id" to table: "language_lesson_members"
CREATE INDEX "index_language_lesson_members_on_language_member_id" ON "language_lesson_members" ("language_member_id");
-- Create index "user_finished_lessons_language_module_lesson_id_index" to table: "language_lesson_members"
CREATE INDEX "user_finished_lessons_language_module_lesson_id_index" ON "language_lesson_members" ("lesson_id");
-- Create index "user_finished_lessons_user_id_language_module_lesson_id_index" to table: "language_lesson_members"
CREATE UNIQUE INDEX "user_finished_lessons_user_id_language_module_lesson_id_index" ON "language_lesson_members" ("user_id", "lesson_id");
-- Create "language_lesson_reviews" table
CREATE TABLE "language_lesson_reviews" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "language_id" bigint NOT NULL,
  "language_lesson_id" bigint NOT NULL,
  "language_lesson_version_id" bigint NOT NULL,
  "language_lesson_version_info_id" bigint NOT NULL,
  "locale" character varying NOT NULL,
  "summary" text NOT NULL DEFAULT '',
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_on_language_lesson_version_info_id_e5ef52eeca" to table: "language_lesson_reviews"
CREATE INDEX "idx_on_language_lesson_version_info_id_e5ef52eeca" ON "language_lesson_reviews" ("language_lesson_version_info_id");
-- Create index "index_language_lesson_reviews_on_language_id" to table: "language_lesson_reviews"
CREATE INDEX "index_language_lesson_reviews_on_language_id" ON "language_lesson_reviews" ("language_id");
-- Create index "index_language_lesson_reviews_on_language_lesson_id" to table: "language_lesson_reviews"
CREATE INDEX "index_language_lesson_reviews_on_language_lesson_id" ON "language_lesson_reviews" ("language_lesson_id");
-- Create index "index_language_lesson_reviews_on_language_lesson_version_id" to table: "language_lesson_reviews"
CREATE INDEX "index_language_lesson_reviews_on_language_lesson_version_id" ON "language_lesson_reviews" ("language_lesson_version_id");
-- Create "language_lesson_version_infos" table
CREATE TABLE "language_lesson_version_infos" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "definitions" character varying NULL,
  "description" character varying NULL,
  "instructions" character varying NULL,
  "language_id" bigint NOT NULL,
  "language_lesson_id" bigint NOT NULL,
  "language_version_id" bigint NOT NULL,
  "locale" character varying NULL,
  "name" character varying NULL,
  "theory" character varying NULL,
  "tips" character varying NULL,
  "updated_at" timestamp NOT NULL,
  "version_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_language_lesson_version_infos_on_language_id" to table: "language_lesson_version_infos"
CREATE INDEX "index_language_lesson_version_infos_on_language_id" ON "language_lesson_version_infos" ("language_id");
-- Create index "index_language_lesson_version_infos_on_language_lesson_id" to table: "language_lesson_version_infos"
CREATE INDEX "index_language_lesson_version_infos_on_language_lesson_id" ON "language_lesson_version_infos" ("language_lesson_id");
-- Create index "index_language_lesson_version_infos_on_language_version_id" to table: "language_lesson_version_infos"
CREATE INDEX "index_language_lesson_version_infos_on_language_version_id" ON "language_lesson_version_infos" ("language_version_id");
-- Create index "index_language_lesson_version_infos_on_version_id_and_locale" to table: "language_lesson_version_infos"
CREATE INDEX "index_language_lesson_version_infos_on_version_id_and_locale" ON "language_lesson_version_infos" ("version_id", "locale");
-- Create "language_lesson_versions" table
CREATE TABLE "language_lesson_versions" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "language_id" bigint NOT NULL,
  "language_version_id" bigint NOT NULL,
  "lesson_id" bigint NOT NULL,
  "module_version_id" bigint NOT NULL,
  "natural_order" integer NULL,
  "order" integer NULL,
  "original_code" character varying NULL,
  "path_to_code" character varying NULL,
  "prepared_code" character varying NULL,
  "test_code" character varying NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_language_lesson_versions_on_language_id" to table: "language_lesson_versions"
CREATE INDEX "index_language_lesson_versions_on_language_id" ON "language_lesson_versions" ("language_id");
-- Create index "index_language_lesson_versions_on_language_version_id" to table: "language_lesson_versions"
CREATE INDEX "index_language_lesson_versions_on_language_version_id" ON "language_lesson_versions" ("language_version_id");
-- Create index "index_language_lesson_versions_on_lesson_id" to table: "language_lesson_versions"
CREATE INDEX "index_language_lesson_versions_on_lesson_id" ON "language_lesson_versions" ("lesson_id");
-- Create index "index_language_lesson_versions_on_module_version_id" to table: "language_lesson_versions"
CREATE INDEX "index_language_lesson_versions_on_module_version_id" ON "language_lesson_versions" ("module_version_id");
-- Create "language_lessons" table
CREATE TABLE "language_lessons" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "language_id" bigint NULL,
  "module_id" bigint NULL,
  "natural_order" integer NULL,
  "order" integer NULL,
  "original_code" text NULL,
  "path_to_code" character varying(255) NULL,
  "prepared_code" text NULL,
  "slug" character varying(255) NULL,
  "state" character varying(255) NULL,
  "test_code" text NULL,
  "updated_at" timestamp NOT NULL,
  "upload_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_language_lessons_on_language_id_and_slug" to table: "language_lessons"
CREATE UNIQUE INDEX "index_language_lessons_on_language_id_and_slug" ON "language_lessons" ("language_id", "slug");
-- Create index "language_module_lessons_module_id_index" to table: "language_lessons"
CREATE INDEX "language_module_lessons_module_id_index" ON "language_lessons" ("module_id");
-- Create index "language_module_lessons_upload_id_index" to table: "language_lessons"
CREATE INDEX "language_module_lessons_upload_id_index" ON "language_lessons" ("upload_id");
-- Create "language_members" table
CREATE TABLE "language_members" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "finished_lessons_count" integer NOT NULL DEFAULT 0,
  "language_id" bigint NOT NULL,
  "state" character varying NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_language_members_on_language_id" to table: "language_members"
CREATE INDEX "index_language_members_on_language_id" ON "language_members" ("language_id");
-- Create index "index_language_members_on_user_id" to table: "language_members"
CREATE INDEX "index_language_members_on_user_id" ON "language_members" ("user_id");
-- Create "language_module_descriptions" table
CREATE TABLE "language_module_descriptions" (
  "id" bigserial NOT NULL,
  "description" text NULL,
  "inserted_at" timestamp NOT NULL,
  "language_id" bigint NULL,
  "locale" character varying(255) NULL,
  "module_id" bigint NULL,
  "name" character varying(255) NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "language_module_descriptions_module_id_index" to table: "language_module_descriptions"
CREATE INDEX "language_module_descriptions_module_id_index" ON "language_module_descriptions" ("module_id");
-- Create "language_module_version_infos" table
CREATE TABLE "language_module_version_infos" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "description" character varying NULL,
  "language_id" bigint NOT NULL,
  "language_version_id" bigint NOT NULL,
  "locale" character varying NULL,
  "name" character varying NULL,
  "updated_at" timestamp NOT NULL,
  "version_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_language_module_version_infos_on_language_id" to table: "language_module_version_infos"
CREATE INDEX "index_language_module_version_infos_on_language_id" ON "language_module_version_infos" ("language_id");
-- Create index "index_language_module_version_infos_on_language_version_id" to table: "language_module_version_infos"
CREATE INDEX "index_language_module_version_infos_on_language_version_id" ON "language_module_version_infos" ("language_version_id");
-- Create "language_module_versions" table
CREATE TABLE "language_module_versions" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "language_id" bigint NOT NULL,
  "language_version_id" bigint NOT NULL,
  "module_id" bigint NOT NULL,
  "order" integer NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_language_module_versions_on_language_id" to table: "language_module_versions"
CREATE INDEX "index_language_module_versions_on_language_id" ON "language_module_versions" ("language_id");
-- Create index "index_language_module_versions_on_language_version_id" to table: "language_module_versions"
CREATE INDEX "index_language_module_versions_on_language_version_id" ON "language_module_versions" ("language_version_id");
-- Create index "index_language_module_versions_on_module_id" to table: "language_module_versions"
CREATE INDEX "index_language_module_versions_on_module_id" ON "language_module_versions" ("module_id");
-- Create "language_modules" table
CREATE TABLE "language_modules" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "language_id" bigint NULL,
  "order" integer NULL,
  "slug" character varying(255) NULL,
  "state" character varying(255) NULL,
  "updated_at" timestamp NOT NULL,
  "upload_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "language_modules_language_id_index" to table: "language_modules"
CREATE INDEX "language_modules_language_id_index" ON "language_modules" ("language_id");
-- Create index "language_modules_upload_id_index" to table: "language_modules"
CREATE INDEX "language_modules_upload_id_index" ON "language_modules" ("upload_id");
-- Create "language_version_infos" table
CREATE TABLE "language_version_infos" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "description" character varying NULL,
  "header" character varying NULL,
  "keywords" character varying NULL,
  "language_id" bigint NOT NULL,
  "language_version_id" bigint NOT NULL,
  "locale" character varying NULL,
  "seo_description" text NULL,
  "title" character varying NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_language_version_infos_on_language_id" to table: "language_version_infos"
CREATE INDEX "index_language_version_infos_on_language_id" ON "language_version_infos" ("language_id");
-- Create index "index_language_version_infos_on_language_version_id" to table: "language_version_infos"
CREATE INDEX "index_language_version_infos_on_language_version_id" ON "language_version_infos" ("language_version_id");
-- Create "language_versions" table
CREATE TABLE "language_versions" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "docker_image" character varying NULL,
  "exercise_filename" character varying NULL,
  "exercise_test_filename" character varying NULL,
  "extension" character varying NULL,
  "language_id" bigint NOT NULL,
  "learn_as" character varying NULL,
  "lessons_count" integer NOT NULL DEFAULT 0,
  "name" character varying NULL,
  "progress" character varying NULL,
  "result" character varying NULL,
  "state" character varying NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_language_versions_on_language_id" to table: "language_versions"
CREATE INDEX "index_language_versions_on_language_id" ON "language_versions" ("language_id");
-- Create "languages" table
CREATE TABLE "languages" (
  "id" bigserial NOT NULL,
  "category_id" bigint NULL,
  "created_at" timestamp NOT NULL,
  "current_version_id" bigint NULL,
  "docker_image" character varying(255) NULL,
  "exercise_filename" character varying(255) NULL,
  "exercise_test_filename" character varying(255) NULL,
  "extension" character varying(255) NULL,
  "hexlet_program_landing_page" character varying NULL,
  "learn_as" character varying NULL,
  "lessons_count" integer NOT NULL DEFAULT 0,
  "members_count" integer NOT NULL DEFAULT 0,
  "name" character varying(255) NULL,
  "order" integer NULL,
  "progress" character varying NULL,
  "slug" character varying(255) NULL,
  "state" character varying(255) NULL,
  "updated_at" timestamp NOT NULL,
  "upload_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_languages_on_category_id" to table: "languages"
CREATE INDEX "index_languages_on_category_id" ON "languages" ("category_id");
-- Create index "index_languages_on_current_version_id" to table: "languages"
CREATE INDEX "index_languages_on_current_version_id" ON "languages" ("current_version_id");
-- Create index "languages_slug_index" to table: "languages"
CREATE UNIQUE INDEX "languages_slug_index" ON "languages" ("slug");
-- Create index "languages_upload_id_index" to table: "languages"
CREATE INDEX "languages_upload_id_index" ON "languages" ("upload_id");
-- Create "leads" table
CREATE TABLE "leads" (
  "id" bigserial NOT NULL,
  "ahoy_visit_id" bigint NULL,
  "courses_data" text NULL,
  "created_at" timestamp NOT NULL,
  "email" character varying NULL,
  "phone" character varying NULL,
  "state" character varying NULL,
  "survey_answers_data" text NULL,
  "telegram" character varying NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NOT NULL,
  "whatsapp" character varying NULL,
  "ym_client_id" character varying NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_leads_on_ahoy_visit_id" to table: "leads"
CREATE INDEX "index_leads_on_ahoy_visit_id" ON "leads" ("ahoy_visit_id");
-- Create index "index_leads_on_user_id" to table: "leads"
CREATE INDEX "index_leads_on_user_id" ON "leads" ("user_id");
-- Create "reviews" table
CREATE TABLE "reviews" (
  "id" bigserial NOT NULL,
  "body" text NULL,
  "created_at" timestamp NOT NULL,
  "first_name" character varying NULL,
  "language_id" bigint NOT NULL,
  "last_name" character varying NULL,
  "locale" character varying NULL,
  "pinned" boolean NULL,
  "state" character varying NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_reviews_on_language_id" to table: "reviews"
CREATE INDEX "index_reviews_on_language_id" ON "reviews" ("language_id");
-- Create index "index_reviews_on_user_id" to table: "reviews"
CREATE INDEX "index_reviews_on_user_id" ON "reviews" ("user_id");
-- Create "schema_migrations" table
CREATE TABLE "schema_migrations" (
  "version" character varying NOT NULL,
  PRIMARY KEY ("version")
);
-- Create "sessions" table
CREATE TABLE "sessions" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "ip_address" character varying NULL,
  "updated_at" timestamp NOT NULL,
  "user_agent" character varying NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_sessions_on_user_id" to table: "sessions"
CREATE INDEX "index_sessions_on_user_id" ON "sessions" ("user_id");
-- Create "solid_cable_messages" table
CREATE TABLE "solid_cable_messages" (
  "id" bigserial NOT NULL,
  "channel" bytea NOT NULL,
  "channel_hash" bigint NOT NULL,
  "created_at" timestamp NOT NULL,
  "payload" bytea NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_cable_messages_on_channel" to table: "solid_cable_messages"
CREATE INDEX "index_solid_cable_messages_on_channel" ON "solid_cable_messages" ("channel");
-- Create index "index_solid_cable_messages_on_channel_hash" to table: "solid_cable_messages"
CREATE INDEX "index_solid_cable_messages_on_channel_hash" ON "solid_cable_messages" ("channel_hash");
-- Create index "index_solid_cable_messages_on_created_at" to table: "solid_cable_messages"
CREATE INDEX "index_solid_cable_messages_on_created_at" ON "solid_cable_messages" ("created_at");
-- Create "solid_cache_entries" table
CREATE TABLE "solid_cache_entries" (
  "id" bigserial NOT NULL,
  "byte_size" integer NOT NULL,
  "created_at" timestamp NOT NULL,
  "key" bytea NOT NULL,
  "key_hash" bigint NOT NULL,
  "value" bytea NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_cache_entries_on_byte_size" to table: "solid_cache_entries"
CREATE INDEX "index_solid_cache_entries_on_byte_size" ON "solid_cache_entries" ("byte_size");
-- Create index "index_solid_cache_entries_on_key_hash" to table: "solid_cache_entries"
CREATE UNIQUE INDEX "index_solid_cache_entries_on_key_hash" ON "solid_cache_entries" ("key_hash");
-- Create index "index_solid_cache_entries_on_key_hash_and_byte_size" to table: "solid_cache_entries"
CREATE INDEX "index_solid_cache_entries_on_key_hash_and_byte_size" ON "solid_cache_entries" ("key_hash", "byte_size");
-- Create "solid_queue_blocked_executions" table
CREATE TABLE "solid_queue_blocked_executions" (
  "id" bigserial NOT NULL,
  "concurrency_key" character varying NOT NULL,
  "created_at" timestamp NOT NULL,
  "expires_at" timestamp NOT NULL,
  "job_id" bigint NOT NULL,
  "priority" integer NOT NULL DEFAULT 0,
  "queue_name" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_queue_blocked_executions_for_maintenance" to table: "solid_queue_blocked_executions"
CREATE INDEX "index_solid_queue_blocked_executions_for_maintenance" ON "solid_queue_blocked_executions" ("expires_at", "concurrency_key");
-- Create index "index_solid_queue_blocked_executions_for_release" to table: "solid_queue_blocked_executions"
CREATE INDEX "index_solid_queue_blocked_executions_for_release" ON "solid_queue_blocked_executions" ("concurrency_key", "priority", "job_id");
-- Create index "index_solid_queue_blocked_executions_on_job_id" to table: "solid_queue_blocked_executions"
CREATE UNIQUE INDEX "index_solid_queue_blocked_executions_on_job_id" ON "solid_queue_blocked_executions" ("job_id");
-- Create "solid_queue_claimed_executions" table
CREATE TABLE "solid_queue_claimed_executions" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "job_id" bigint NOT NULL,
  "process_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_queue_claimed_executions_on_job_id" to table: "solid_queue_claimed_executions"
CREATE UNIQUE INDEX "index_solid_queue_claimed_executions_on_job_id" ON "solid_queue_claimed_executions" ("job_id");
-- Create index "index_solid_queue_claimed_executions_on_process_id_and_job_id" to table: "solid_queue_claimed_executions"
CREATE INDEX "index_solid_queue_claimed_executions_on_process_id_and_job_id" ON "solid_queue_claimed_executions" ("process_id", "job_id");
-- Create "solid_queue_failed_executions" table
CREATE TABLE "solid_queue_failed_executions" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "error" text NULL,
  "job_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_queue_failed_executions_on_job_id" to table: "solid_queue_failed_executions"
CREATE UNIQUE INDEX "index_solid_queue_failed_executions_on_job_id" ON "solid_queue_failed_executions" ("job_id");
-- Create "solid_queue_jobs" table
CREATE TABLE "solid_queue_jobs" (
  "id" bigserial NOT NULL,
  "active_job_id" character varying NULL,
  "arguments" text NULL,
  "class_name" character varying NOT NULL,
  "concurrency_key" character varying NULL,
  "created_at" timestamp NOT NULL,
  "finished_at" timestamp NULL,
  "priority" integer NOT NULL DEFAULT 0,
  "queue_name" character varying NOT NULL,
  "scheduled_at" timestamp NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_queue_jobs_for_alerting" to table: "solid_queue_jobs"
CREATE INDEX "index_solid_queue_jobs_for_alerting" ON "solid_queue_jobs" ("scheduled_at", "finished_at");
-- Create index "index_solid_queue_jobs_for_filtering" to table: "solid_queue_jobs"
CREATE INDEX "index_solid_queue_jobs_for_filtering" ON "solid_queue_jobs" ("queue_name", "finished_at");
-- Create index "index_solid_queue_jobs_on_active_job_id" to table: "solid_queue_jobs"
CREATE INDEX "index_solid_queue_jobs_on_active_job_id" ON "solid_queue_jobs" ("active_job_id");
-- Create index "index_solid_queue_jobs_on_class_name" to table: "solid_queue_jobs"
CREATE INDEX "index_solid_queue_jobs_on_class_name" ON "solid_queue_jobs" ("class_name");
-- Create index "index_solid_queue_jobs_on_finished_at" to table: "solid_queue_jobs"
CREATE INDEX "index_solid_queue_jobs_on_finished_at" ON "solid_queue_jobs" ("finished_at");
-- Create "solid_queue_pauses" table
CREATE TABLE "solid_queue_pauses" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "queue_name" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_queue_pauses_on_queue_name" to table: "solid_queue_pauses"
CREATE UNIQUE INDEX "index_solid_queue_pauses_on_queue_name" ON "solid_queue_pauses" ("queue_name");
-- Create "solid_queue_processes" table
CREATE TABLE "solid_queue_processes" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "hostname" character varying NULL,
  "kind" character varying NOT NULL,
  "last_heartbeat_at" timestamp NOT NULL,
  "metadata" text NULL,
  "name" character varying NOT NULL,
  "pid" integer NOT NULL,
  "supervisor_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_queue_processes_on_last_heartbeat_at" to table: "solid_queue_processes"
CREATE INDEX "index_solid_queue_processes_on_last_heartbeat_at" ON "solid_queue_processes" ("last_heartbeat_at");
-- Create index "index_solid_queue_processes_on_name_and_supervisor_id" to table: "solid_queue_processes"
CREATE UNIQUE INDEX "index_solid_queue_processes_on_name_and_supervisor_id" ON "solid_queue_processes" ("name", "supervisor_id");
-- Create index "index_solid_queue_processes_on_supervisor_id" to table: "solid_queue_processes"
CREATE INDEX "index_solid_queue_processes_on_supervisor_id" ON "solid_queue_processes" ("supervisor_id");
-- Create "solid_queue_ready_executions" table
CREATE TABLE "solid_queue_ready_executions" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "job_id" bigint NOT NULL,
  "priority" integer NOT NULL DEFAULT 0,
  "queue_name" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_queue_poll_all" to table: "solid_queue_ready_executions"
CREATE INDEX "index_solid_queue_poll_all" ON "solid_queue_ready_executions" ("priority", "job_id");
-- Create index "index_solid_queue_poll_by_queue" to table: "solid_queue_ready_executions"
CREATE INDEX "index_solid_queue_poll_by_queue" ON "solid_queue_ready_executions" ("queue_name", "priority", "job_id");
-- Create index "index_solid_queue_ready_executions_on_job_id" to table: "solid_queue_ready_executions"
CREATE UNIQUE INDEX "index_solid_queue_ready_executions_on_job_id" ON "solid_queue_ready_executions" ("job_id");
-- Create "solid_queue_recurring_executions" table
CREATE TABLE "solid_queue_recurring_executions" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "job_id" bigint NOT NULL,
  "run_at" timestamp NOT NULL,
  "task_key" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_queue_recurring_executions_on_job_id" to table: "solid_queue_recurring_executions"
CREATE UNIQUE INDEX "index_solid_queue_recurring_executions_on_job_id" ON "solid_queue_recurring_executions" ("job_id");
-- Create index "index_solid_queue_recurring_executions_on_task_key_and_run_at" to table: "solid_queue_recurring_executions"
CREATE UNIQUE INDEX "index_solid_queue_recurring_executions_on_task_key_and_run_at" ON "solid_queue_recurring_executions" ("task_key", "run_at");
-- Create "solid_queue_recurring_tasks" table
CREATE TABLE "solid_queue_recurring_tasks" (
  "id" bigserial NOT NULL,
  "arguments" text NULL,
  "class_name" character varying NULL,
  "command" character varying(2048) NULL,
  "created_at" timestamp NOT NULL,
  "description" text NULL,
  "key" character varying NOT NULL,
  "priority" integer NULL DEFAULT 0,
  "queue_name" character varying NULL,
  "schedule" character varying NOT NULL,
  "static" boolean NOT NULL DEFAULT true,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_queue_recurring_tasks_on_key" to table: "solid_queue_recurring_tasks"
CREATE UNIQUE INDEX "index_solid_queue_recurring_tasks_on_key" ON "solid_queue_recurring_tasks" ("key");
-- Create index "index_solid_queue_recurring_tasks_on_static" to table: "solid_queue_recurring_tasks"
CREATE INDEX "index_solid_queue_recurring_tasks_on_static" ON "solid_queue_recurring_tasks" ("static");
-- Create "solid_queue_scheduled_executions" table
CREATE TABLE "solid_queue_scheduled_executions" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "job_id" bigint NOT NULL,
  "priority" integer NOT NULL DEFAULT 0,
  "queue_name" character varying NOT NULL,
  "scheduled_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_queue_dispatch_all" to table: "solid_queue_scheduled_executions"
CREATE INDEX "index_solid_queue_dispatch_all" ON "solid_queue_scheduled_executions" ("scheduled_at", "priority", "job_id");
-- Create index "index_solid_queue_scheduled_executions_on_job_id" to table: "solid_queue_scheduled_executions"
CREATE UNIQUE INDEX "index_solid_queue_scheduled_executions_on_job_id" ON "solid_queue_scheduled_executions" ("job_id");
-- Create "solid_queue_semaphores" table
CREATE TABLE "solid_queue_semaphores" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "expires_at" timestamp NOT NULL,
  "key" character varying NOT NULL,
  "updated_at" timestamp NOT NULL,
  "value" integer NOT NULL DEFAULT 1,
  PRIMARY KEY ("id")
);
-- Create index "index_solid_queue_semaphores_on_expires_at" to table: "solid_queue_semaphores"
CREATE INDEX "index_solid_queue_semaphores_on_expires_at" ON "solid_queue_semaphores" ("expires_at");
-- Create index "index_solid_queue_semaphores_on_key" to table: "solid_queue_semaphores"
CREATE UNIQUE INDEX "index_solid_queue_semaphores_on_key" ON "solid_queue_semaphores" ("key");
-- Create index "index_solid_queue_semaphores_on_key_and_value" to table: "solid_queue_semaphores"
CREATE INDEX "index_solid_queue_semaphores_on_key_and_value" ON "solid_queue_semaphores" ("key", "value");
-- Create "staff_member_role_permissions" table
CREATE TABLE "staff_member_role_permissions" (
  "id" bigserial NOT NULL,
  "can_create" boolean NOT NULL DEFAULT false,
  "can_destroy" boolean NOT NULL DEFAULT false,
  "can_index" boolean NOT NULL DEFAULT false,
  "can_update" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL,
  "resource" character varying NOT NULL,
  "role_id" bigint NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_staff_member_role_permissions_on_role_id_and_resource" to table: "staff_member_role_permissions"
CREATE UNIQUE INDEX "index_staff_member_role_permissions_on_role_id_and_resource" ON "staff_member_role_permissions" ("role_id", "resource");
-- Create "staff_member_roles" table
CREATE TABLE "staff_member_roles" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "description" text NULL,
  "name" character varying NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "staff_members" table
CREATE TABLE "staff_members" (
  "id" bigserial NOT NULL,
  "allowed_locales" character varying[] NOT NULL DEFAULT '{ru}',
  "created_at" timestamp NOT NULL,
  "role_id" bigint NOT NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_staff_members_on_role_id" to table: "staff_members"
CREATE INDEX "index_staff_members_on_role_id" ON "staff_members" ("role_id");
-- Create index "index_staff_members_on_user_id" to table: "staff_members"
CREATE UNIQUE INDEX "index_staff_members_on_user_id" ON "staff_members" ("user_id");
-- Create "survey_answers" table
CREATE TABLE "survey_answers" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "state" character varying NULL,
  "survey_id" bigint NOT NULL,
  "survey_item_id" bigint NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_survey_answers_on_survey_id_and_user_id" to table: "survey_answers"
CREATE UNIQUE INDEX "index_survey_answers_on_survey_id_and_user_id" ON "survey_answers" ("survey_id", "user_id");
-- Create index "index_survey_answers_on_survey_item_id" to table: "survey_answers"
CREATE INDEX "index_survey_answers_on_survey_item_id" ON "survey_answers" ("survey_item_id");
-- Create index "index_survey_answers_on_user_id" to table: "survey_answers"
CREATE INDEX "index_survey_answers_on_user_id" ON "survey_answers" ("user_id");
-- Create "survey_items" table
CREATE TABLE "survey_items" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "order" integer NOT NULL,
  "slug" character varying NULL,
  "state" character varying NULL,
  "survey_id" bigint NOT NULL,
  "updated_at" timestamp NOT NULL,
  "value" character varying NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_survey_items_on_survey_id" to table: "survey_items"
CREATE INDEX "index_survey_items_on_survey_id" ON "survey_items" ("survey_id");
-- Create "survey_scenario_items" table
CREATE TABLE "survey_scenario_items" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "order" integer NULL,
  "scenario_id" bigint NOT NULL,
  "survey_id" bigint NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_survey_scenario_items_on_scenario_id" to table: "survey_scenario_items"
CREATE INDEX "index_survey_scenario_items_on_scenario_id" ON "survey_scenario_items" ("scenario_id");
-- Create index "index_survey_scenario_items_on_survey_id_and_scenario_id" to table: "survey_scenario_items"
CREATE UNIQUE INDEX "index_survey_scenario_items_on_survey_id_and_scenario_id" ON "survey_scenario_items" ("survey_id", "scenario_id");
-- Create "survey_scenario_members" table
CREATE TABLE "survey_scenario_members" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "event_name" character varying NULL,
  "scenario_id" bigint NOT NULL,
  "state" character varying NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_survey_scenario_members_on_scenario_id" to table: "survey_scenario_members"
CREATE INDEX "index_survey_scenario_members_on_scenario_id" ON "survey_scenario_members" ("scenario_id");
-- Create index "index_survey_scenario_members_on_user_id" to table: "survey_scenario_members"
CREATE INDEX "index_survey_scenario_members_on_user_id" ON "survey_scenario_members" ("user_id");
-- Create "survey_scenario_triggers" table
CREATE TABLE "survey_scenario_triggers" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "event_name" character varying NULL,
  "event_threshold_count" integer NULL,
  "scenario_id" bigint NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_survey_scenario_triggers_on_event_name_and_scenario_id" to table: "survey_scenario_triggers"
CREATE UNIQUE INDEX "index_survey_scenario_triggers_on_event_name_and_scenario_id" ON "survey_scenario_triggers" ("event_name", "scenario_id");
-- Create index "index_survey_scenario_triggers_on_scenario_id" to table: "survey_scenario_triggers"
CREATE INDEX "index_survey_scenario_triggers_on_scenario_id" ON "survey_scenario_triggers" ("scenario_id");
-- Create "survey_scenarios" table
CREATE TABLE "survey_scenarios" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "locale" character varying NULL,
  "name" character varying NULL,
  "state" character varying NULL,
  "survey_item_id" bigint NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_survey_scenarios_on_survey_item_id" to table: "survey_scenarios"
CREATE INDEX "index_survey_scenarios_on_survey_item_id" ON "survey_scenarios" ("survey_item_id");
-- Create "surveys" table
CREATE TABLE "surveys" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "description" character varying NULL,
  "locale" character varying NULL,
  "parent_survey_id" bigint NULL,
  "parent_survey_item_id" bigint NULL,
  "question" character varying NULL,
  "run_after_finishing_lessons_count" integer NULL DEFAULT 0,
  "run_always" boolean NULL DEFAULT false,
  "slug" character varying NULL,
  "state" character varying NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_rails_c84299bbb9" FOREIGN KEY ("parent_survey_id") REFERENCES "surveys" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "index_surveys_on_parent_survey_id" to table: "surveys"
CREATE INDEX "index_surveys_on_parent_survey_id" ON "surveys" ("parent_survey_id");
-- Create index "index_surveys_on_parent_survey_item_id" to table: "surveys"
CREATE INDEX "index_surveys_on_parent_survey_item_id" ON "surveys" ("parent_survey_item_id");
-- Create index "index_surveys_on_slug_and_locale" to table: "surveys"
CREATE UNIQUE INDEX "index_surveys_on_slug_and_locale" ON "surveys" ("slug", "locale");
-- Create "taggings" table
CREATE TABLE "taggings" (
  "id" bigserial NOT NULL,
  "context" character varying(128) NULL,
  "created_at" timestamp NULL,
  "tag_id" bigint NULL,
  "taggable_id" bigint NULL,
  "taggable_type" character varying NULL,
  "tagger_id" bigint NULL,
  "tagger_type" character varying NULL,
  "tenant" character varying(128) NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_taggings_on_context" to table: "taggings"
CREATE INDEX "index_taggings_on_context" ON "taggings" ("context");
-- Create index "index_taggings_on_tag_id" to table: "taggings"
CREATE INDEX "index_taggings_on_tag_id" ON "taggings" ("tag_id");
-- Create index "index_taggings_on_taggable_id" to table: "taggings"
CREATE INDEX "index_taggings_on_taggable_id" ON "taggings" ("taggable_id");
-- Create index "index_taggings_on_taggable_type" to table: "taggings"
CREATE INDEX "index_taggings_on_taggable_type" ON "taggings" ("taggable_type");
-- Create index "index_taggings_on_taggable_type_and_taggable_id" to table: "taggings"
CREATE INDEX "index_taggings_on_taggable_type_and_taggable_id" ON "taggings" ("taggable_type", "taggable_id");
-- Create index "index_taggings_on_tagger_id" to table: "taggings"
CREATE INDEX "index_taggings_on_tagger_id" ON "taggings" ("tagger_id");
-- Create index "index_taggings_on_tagger_id_and_tagger_type" to table: "taggings"
CREATE INDEX "index_taggings_on_tagger_id_and_tagger_type" ON "taggings" ("tagger_id", "tagger_type");
-- Create index "index_taggings_on_tagger_type_and_tagger_id" to table: "taggings"
CREATE INDEX "index_taggings_on_tagger_type_and_tagger_id" ON "taggings" ("tagger_type", "tagger_id");
-- Create index "index_taggings_on_tenant" to table: "taggings"
CREATE INDEX "index_taggings_on_tenant" ON "taggings" ("tenant");
-- Create index "taggings_idx" to table: "taggings"
CREATE UNIQUE INDEX "taggings_idx" ON "taggings" ("tag_id", "taggable_id", "taggable_type", "context", "tagger_id", "tagger_type");
-- Create index "taggings_idy" to table: "taggings"
CREATE INDEX "taggings_idy" ON "taggings" ("taggable_id", "taggable_type", "tagger_id", "context");
-- Create index "taggings_taggable_context_idx" to table: "taggings"
CREATE INDEX "taggings_taggable_context_idx" ON "taggings" ("taggable_id", "taggable_type", "context");
-- Create "tags" table
CREATE TABLE "tags" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "name" character varying NULL,
  "taggings_count" integer NULL DEFAULT 0,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_tags_on_name" to table: "tags"
CREATE UNIQUE INDEX "index_tags_on_name" ON "tags" ("name");
-- Create "uploads" table
CREATE TABLE "uploads" (
  "id" bigserial NOT NULL,
  "inserted_at" timestamp NOT NULL,
  "language_name" character varying(255) NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "user_accounts" table
CREATE TABLE "user_accounts" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "provider" character varying(255) NOT NULL,
  "uid" character varying(255) NOT NULL,
  "updated_at" timestamp(0) NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_user_accounts_on_provider_and_uid" to table: "user_accounts"
CREATE UNIQUE INDEX "index_user_accounts_on_provider_and_uid" ON "user_accounts" ("provider", "uid");
-- Create "user_credentials" table
CREATE TABLE "user_credentials" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "external_id" character varying NOT NULL,
  "nickname" character varying NULL,
  "public_key" character varying NOT NULL,
  "sign_count" bigint NOT NULL DEFAULT 0,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_user_credentials_on_external_id" to table: "user_credentials"
CREATE UNIQUE INDEX "index_user_credentials_on_external_id" ON "user_credentials" ("external_id");
-- Create index "index_user_credentials_on_user_id" to table: "user_credentials"
CREATE INDEX "index_user_credentials_on_user_id" ON "user_credentials" ("user_id");
-- Create "user_survey_pivots" table
CREATE TABLE "user_survey_pivots" (
  "id" bigserial NOT NULL,
  "coding_experience_item_id" bigint NULL,
  "created_at" timestamp NOT NULL,
  "goal_item_id" bigint NULL,
  "study_plan_item_id" bigint NULL,
  "updated_at" timestamp NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_user_survey_pivots_on_coding_experience_item_id" to table: "user_survey_pivots"
CREATE INDEX "index_user_survey_pivots_on_coding_experience_item_id" ON "user_survey_pivots" ("coding_experience_item_id");
-- Create index "index_user_survey_pivots_on_goal_item_id" to table: "user_survey_pivots"
CREATE INDEX "index_user_survey_pivots_on_goal_item_id" ON "user_survey_pivots" ("goal_item_id");
-- Create index "index_user_survey_pivots_on_study_plan_item_id" to table: "user_survey_pivots"
CREATE INDEX "index_user_survey_pivots_on_study_plan_item_id" ON "user_survey_pivots" ("study_plan_item_id");
-- Create index "index_user_survey_pivots_on_user_id" to table: "user_survey_pivots"
CREATE INDEX "index_user_survey_pivots_on_user_id" ON "user_survey_pivots" ("user_id");
-- Create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "admin" boolean NULL,
  "assistant_messages_count" integer NULL DEFAULT 0,
  "confirmation_token" character varying(255) NULL,
  "contact_method" character varying NULL,
  "contact_value" character varying NULL,
  "created_at" timestamp NOT NULL,
  "email" character varying(255) NULL,
  "email_delivery_state" character varying(255) NULL,
  "facebook_uid" character varying(255) NULL,
  "first_name" character varying(255) NULL,
  "github_uid" integer NULL,
  "help" boolean NULL,
  "last_name" character varying(255) NULL,
  "locale" character varying(255) NULL,
  "nickname" character varying(255) NULL,
  "password_digest" character varying(255) NULL,
  "phone" character varying NULL,
  "phone_verified_at" timestamp NULL,
  "state" character varying(255) NULL,
  "updated_at" timestamp NOT NULL,
  "webauthn_id" character varying NULL,
  PRIMARY KEY ("id")
);
-- Create index "index_users_on_LOWER_email" to table: "users"
CREATE UNIQUE INDEX "index_users_on_LOWER_email" ON "users" ((lower((email)::text)));
-- Create index "index_users_on_email" to table: "users"
CREATE UNIQUE INDEX "index_users_on_email" ON "users" ("email");
-- Create index "index_users_on_phone" to table: "users"
CREATE UNIQUE INDEX "index_users_on_phone" ON "users" ("phone") WHERE (phone IS NOT NULL);
-- Create index "index_users_on_webauthn_id" to table: "users"
CREATE UNIQUE INDEX "index_users_on_webauthn_id" ON "users" ("webauthn_id") WHERE (webauthn_id IS NOT NULL);
-- Modify "active_storage_attachments" table
ALTER TABLE "active_storage_attachments" ADD CONSTRAINT "fk_rails_c3b3935057" FOREIGN KEY ("blob_id") REFERENCES "active_storage_blobs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "active_storage_variant_records" table
ALTER TABLE "active_storage_variant_records" ADD CONSTRAINT "fk_rails_993965df05" FOREIGN KEY ("blob_id") REFERENCES "active_storage_blobs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "ai_chats" table
ALTER TABLE "ai_chats" ADD CONSTRAINT "fk_rails_4831695bfa" FOREIGN KEY ("language_lesson_member_id") REFERENCES "language_lesson_members" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_768e14b856" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_e112a98abc" FOREIGN KEY ("ai_model_id") REFERENCES "ai_models" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "ai_messages" table
ALTER TABLE "ai_messages" ADD CONSTRAINT "fk_rails_438986158e" FOREIGN KEY ("ai_model_id") REFERENCES "ai_models" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_69368b3cd9" FOREIGN KEY ("ai_tool_call_id") REFERENCES "ai_tool_calls" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_958af724a1" FOREIGN KEY ("ai_chat_id") REFERENCES "ai_chats" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_eaeb97aaaa" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "ai_tool_calls" table
ALTER TABLE "ai_tool_calls" ADD CONSTRAINT "fk_rails_1ebf54e503" FOREIGN KEY ("ai_message_id") REFERENCES "ai_messages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "blog_post_likes" table
ALTER TABLE "blog_post_likes" ADD CONSTRAINT "fk_rails_2c9e0d4a09" FOREIGN KEY ("blog_post_id") REFERENCES "blog_posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_326553c5d9" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "blog_post_related_language_items" table
ALTER TABLE "blog_post_related_language_items" ADD CONSTRAINT "fk_rails_423f6248fd" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_d74a547da2" FOREIGN KEY ("blog_post_id") REFERENCES "blog_posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "blog_posts" table
ALTER TABLE "blog_posts" ADD CONSTRAINT "fk_rails_399b6b9958" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_45731ba462" FOREIGN KEY ("creator_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "book_requests" table
ALTER TABLE "book_requests" ADD CONSTRAINT "fk_rails_1e4327d03b" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "event_store_events_in_streams" table
ALTER TABLE "event_store_events_in_streams" ADD CONSTRAINT "fk_rails_c8d52b5857" FOREIGN KEY ("event_id") REFERENCES "event_store_events" ("event_id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_category_items" table
ALTER TABLE "language_category_items" ADD CONSTRAINT "fk_rails_1653d28321" FOREIGN KEY ("language_category_id") REFERENCES "language_categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_50f4828843" FOREIGN KEY ("language_landing_page_id") REFERENCES "language_landing_pages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_category_qna_items" table
ALTER TABLE "language_category_qna_items" ADD CONSTRAINT "fk_rails_3ba389c392" FOREIGN KEY ("language_category_id") REFERENCES "language_categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_landing_page_qna_items" table
ALTER TABLE "language_landing_page_qna_items" ADD CONSTRAINT "fk_rails_c3a22366db" FOREIGN KEY ("language_landing_page_id") REFERENCES "language_landing_pages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_landing_pages" table
ALTER TABLE "language_landing_pages" ADD CONSTRAINT "fk_rails_77f646d57e" FOREIGN KEY ("language_category_id") REFERENCES "language_categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_f28dfb0a77" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_lesson_members" table
ALTER TABLE "language_lesson_members" ADD CONSTRAINT "fk_rails_1ff6af1a44" FOREIGN KEY ("lesson_id") REFERENCES "language_lessons" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_7e60189e01" FOREIGN KEY ("language_member_id") REFERENCES "language_members" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_eb00b8f36c" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_lesson_reviews" table
ALTER TABLE "language_lesson_reviews" ADD CONSTRAINT "fk_rails_58fe2f4daf" FOREIGN KEY ("language_lesson_version_info_id") REFERENCES "language_lesson_version_infos" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_8bf293e8a9" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_9041c890b0" FOREIGN KEY ("language_lesson_version_id") REFERENCES "language_lesson_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_93e176fee3" FOREIGN KEY ("language_lesson_id") REFERENCES "language_lessons" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_lesson_version_infos" table
ALTER TABLE "language_lesson_version_infos" ADD CONSTRAINT "fk_rails_2b2de76835" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_3e2e50d44d" FOREIGN KEY ("language_version_id") REFERENCES "language_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_749db9acb4" FOREIGN KEY ("language_lesson_id") REFERENCES "language_lessons" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_8872f568fe" FOREIGN KEY ("version_id") REFERENCES "language_lesson_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_lesson_versions" table
ALTER TABLE "language_lesson_versions" ADD CONSTRAINT "fk_rails_28cb0a511a" FOREIGN KEY ("lesson_id") REFERENCES "language_lessons" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_4cae0d7625" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_6b52e26355" FOREIGN KEY ("language_version_id") REFERENCES "language_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_fd3c3cf805" FOREIGN KEY ("module_version_id") REFERENCES "language_module_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_lessons" table
ALTER TABLE "language_lessons" ADD CONSTRAINT "fk_rails_36dee53ae2" FOREIGN KEY ("module_id") REFERENCES "language_modules" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_439e3990bb" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_9ca0d6ceaa" FOREIGN KEY ("upload_id") REFERENCES "uploads" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_members" table
ALTER TABLE "language_members" ADD CONSTRAINT "fk_rails_4343116778" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_86f99b4837" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_module_descriptions" table
ALTER TABLE "language_module_descriptions" ADD CONSTRAINT "fk_rails_1cc3025f38" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_dac1bb6244" FOREIGN KEY ("module_id") REFERENCES "language_modules" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_module_version_infos" table
ALTER TABLE "language_module_version_infos" ADD CONSTRAINT "fk_rails_03a490d994" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_9c7b71797b" FOREIGN KEY ("version_id") REFERENCES "language_module_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_e4269a6401" FOREIGN KEY ("language_version_id") REFERENCES "language_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_module_versions" table
ALTER TABLE "language_module_versions" ADD CONSTRAINT "fk_rails_2a94cf2f1a" FOREIGN KEY ("module_id") REFERENCES "language_modules" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_413afabb56" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_4b0d9ee90f" FOREIGN KEY ("language_version_id") REFERENCES "language_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_modules" table
ALTER TABLE "language_modules" ADD CONSTRAINT "fk_rails_39957b735c" FOREIGN KEY ("upload_id") REFERENCES "uploads" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_5c574dddfe" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_version_infos" table
ALTER TABLE "language_version_infos" ADD CONSTRAINT "fk_rails_11a0eeec04" FOREIGN KEY ("language_version_id") REFERENCES "language_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_6ffda837cf" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "language_versions" table
ALTER TABLE "language_versions" ADD CONSTRAINT "fk_rails_6cc776ff38" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "languages" table
ALTER TABLE "languages" ADD CONSTRAINT "fk_rails_bcc060b35a" FOREIGN KEY ("upload_id") REFERENCES "uploads" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_cd49b170ef" FOREIGN KEY ("category_id") REFERENCES "language_categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_f53e1946e0" FOREIGN KEY ("current_version_id") REFERENCES "language_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "leads" table
ALTER TABLE "leads" ADD CONSTRAINT "fk_rails_1d08b36969" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_daafd5abf0" FOREIGN KEY ("ahoy_visit_id") REFERENCES "ahoy_visits" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "reviews" table
ALTER TABLE "reviews" ADD CONSTRAINT "fk_rails_2d5506f396" FOREIGN KEY ("language_id") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_74a66bd6c5" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "sessions" table
ALTER TABLE "sessions" ADD CONSTRAINT "fk_rails_758836b4f0" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "solid_queue_blocked_executions" table
ALTER TABLE "solid_queue_blocked_executions" ADD CONSTRAINT "fk_rails_4cd34e2228" FOREIGN KEY ("job_id") REFERENCES "solid_queue_jobs" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "solid_queue_claimed_executions" table
ALTER TABLE "solid_queue_claimed_executions" ADD CONSTRAINT "fk_rails_9cfe4d4944" FOREIGN KEY ("job_id") REFERENCES "solid_queue_jobs" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "solid_queue_failed_executions" table
ALTER TABLE "solid_queue_failed_executions" ADD CONSTRAINT "fk_rails_39bbc7a631" FOREIGN KEY ("job_id") REFERENCES "solid_queue_jobs" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "solid_queue_ready_executions" table
ALTER TABLE "solid_queue_ready_executions" ADD CONSTRAINT "fk_rails_81fcbd66af" FOREIGN KEY ("job_id") REFERENCES "solid_queue_jobs" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "solid_queue_recurring_executions" table
ALTER TABLE "solid_queue_recurring_executions" ADD CONSTRAINT "fk_rails_318a5533ed" FOREIGN KEY ("job_id") REFERENCES "solid_queue_jobs" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "solid_queue_scheduled_executions" table
ALTER TABLE "solid_queue_scheduled_executions" ADD CONSTRAINT "fk_rails_c4316f352d" FOREIGN KEY ("job_id") REFERENCES "solid_queue_jobs" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "staff_member_role_permissions" table
ALTER TABLE "staff_member_role_permissions" ADD CONSTRAINT "fk_rails_5d1aecb760" FOREIGN KEY ("role_id") REFERENCES "staff_member_roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "staff_members" table
ALTER TABLE "staff_members" ADD CONSTRAINT "fk_rails_2925f894a8" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_7fe1e4b2a6" FOREIGN KEY ("role_id") REFERENCES "staff_member_roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "survey_answers" table
ALTER TABLE "survey_answers" ADD CONSTRAINT "fk_rails_3869acb601" FOREIGN KEY ("survey_id") REFERENCES "surveys" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_621f80522c" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_8cb45b53a9" FOREIGN KEY ("survey_item_id") REFERENCES "survey_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "survey_items" table
ALTER TABLE "survey_items" ADD CONSTRAINT "fk_rails_bdccd8c655" FOREIGN KEY ("survey_id") REFERENCES "surveys" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "survey_scenario_items" table
ALTER TABLE "survey_scenario_items" ADD CONSTRAINT "fk_rails_2eb91bfb82" FOREIGN KEY ("survey_id") REFERENCES "surveys" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_924a4eb1ae" FOREIGN KEY ("scenario_id") REFERENCES "survey_scenarios" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "survey_scenario_members" table
ALTER TABLE "survey_scenario_members" ADD CONSTRAINT "fk_rails_463c98f012" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_5db85e9ca1" FOREIGN KEY ("scenario_id") REFERENCES "survey_scenarios" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "survey_scenario_triggers" table
ALTER TABLE "survey_scenario_triggers" ADD CONSTRAINT "fk_rails_c8040740a8" FOREIGN KEY ("scenario_id") REFERENCES "survey_scenarios" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "survey_scenarios" table
ALTER TABLE "survey_scenarios" ADD CONSTRAINT "fk_rails_f2414b87c2" FOREIGN KEY ("survey_item_id") REFERENCES "survey_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "surveys" table
ALTER TABLE "surveys" ADD CONSTRAINT "fk_rails_7c2d7aa117" FOREIGN KEY ("parent_survey_item_id") REFERENCES "survey_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "taggings" table
ALTER TABLE "taggings" ADD CONSTRAINT "fk_rails_9fcd2e236b" FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "user_accounts" table
ALTER TABLE "user_accounts" ADD CONSTRAINT "fk_rails_d64ac9bcc2" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "user_credentials" table
ALTER TABLE "user_credentials" ADD CONSTRAINT "fk_rails_9b162a81f6" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "user_survey_pivots" table
ALTER TABLE "user_survey_pivots" ADD CONSTRAINT "fk_rails_23a507246e" FOREIGN KEY ("goal_item_id") REFERENCES "survey_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_535f0749ed" FOREIGN KEY ("coding_experience_item_id") REFERENCES "survey_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_5ca50bd2d9" FOREIGN KEY ("study_plan_item_id") REFERENCES "survey_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_rails_eb41950872" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
