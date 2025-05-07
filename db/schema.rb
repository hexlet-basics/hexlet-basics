# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[8.0].define(version: 2025_05_08_175604) do
  # These are extensions that must be enabled in order to support this database
  enable_extension "pg_catalog.plpgsql"

  create_table "active_storage_attachments", force: :cascade do |t|
    t.string "name", null: false
    t.string "record_type", null: false
    t.bigint "record_id", null: false
    t.bigint "blob_id", null: false
    t.datetime "created_at", null: false
    t.index ["blob_id"], name: "index_active_storage_attachments_on_blob_id"
    t.index ["record_type", "record_id", "name", "blob_id"], name: "index_active_storage_attachments_uniqueness", unique: true
  end

  create_table "active_storage_blobs", force: :cascade do |t|
    t.string "key", null: false
    t.string "filename", null: false
    t.string "content_type"
    t.text "metadata"
    t.string "service_name", null: false
    t.bigint "byte_size", null: false
    t.string "checksum"
    t.datetime "created_at", null: false
    t.index ["key"], name: "index_active_storage_blobs_on_key", unique: true
  end

  create_table "active_storage_variant_records", force: :cascade do |t|
    t.bigint "blob_id", null: false
    t.string "variation_digest", null: false
    t.index ["blob_id", "variation_digest"], name: "index_active_storage_variant_records_uniqueness", unique: true
  end

  create_table "blog_posts", force: :cascade do |t|
    t.bigint "language_id"
    t.string "locale"
    t.string "state"
    t.string "slug"
    t.string "name"
    t.text "body"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.bigint "creator_id", null: false
    t.string "description"
    t.index ["creator_id"], name: "index_blog_posts_on_creator_id"
    t.index ["language_id"], name: "index_blog_posts_on_language_id"
    t.index ["slug"], name: "index_blog_posts_on_slug", unique: true
  end

  create_table "book_requests", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.string "state"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["user_id"], name: "index_book_requests_on_user_id", unique: true
  end

  create_table "course_categories", force: :cascade do |t|
    t.string "name"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "event_store_events", force: :cascade do |t|
    t.string "event_id", limit: 36, null: false
    t.string "event_type", null: false
    t.binary "metadata"
    t.binary "data", null: false
    t.datetime "created_at", null: false
    t.datetime "valid_at"
    t.index ["created_at"], name: "index_event_store_events_on_created_at"
    t.index ["event_id"], name: "index_event_store_events_on_event_id", unique: true
    t.index ["event_type"], name: "index_event_store_events_on_event_type"
    t.index ["valid_at"], name: "index_event_store_events_on_valid_at"
  end

  create_table "event_store_events_in_streams", force: :cascade do |t|
    t.string "stream", null: false
    t.integer "position"
    t.string "event_id", limit: 36, null: false
    t.datetime "created_at", null: false
    t.index ["created_at"], name: "index_event_store_events_in_streams_on_created_at"
    t.index ["event_id"], name: "index_event_store_events_in_streams_on_event_id"
    t.index ["stream", "event_id"], name: "index_event_store_events_in_streams_on_stream_and_event_id", unique: true
    t.index ["stream", "position"], name: "index_event_store_events_in_streams_on_stream_and_position", unique: true
  end

  create_table "language_categories", force: :cascade do |t|
    t.string "name_ru"
    t.string "name_en"
    t.string "slug"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.string "name"
    t.string "locale"
  end

  create_table "language_landing_page_qna_items", force: :cascade do |t|
    t.integer "language_landing_page_id", null: false
    t.string "question"
    t.string "answer"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["language_landing_page_id"], name: "idx_on_language_landing_page_id_98023e1f90"
  end

  create_table "language_landing_pages", force: :cascade do |t|
    t.integer "language_id", null: false
    t.integer "language_category_id"
    t.string "meta_title"
    t.string "locale"
    t.string "header"
    t.string "slug"
    t.string "order"
    t.boolean "main"
    t.string "state"
    t.string "description"
    t.string "meta_description"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.boolean "listed"
    t.string "used_in_header"
    t.string "used_in_description"
    t.string "outcomes_header"
    t.string "outcomes_description"
    t.boolean "footer"
    t.string "footer_name"
    t.index ["language_category_id"], name: "index_language_landing_pages_on_language_category_id"
    t.index ["language_id"], name: "index_language_landing_pages_on_language_id"
  end

  create_table "language_lesson_member_messages", force: :cascade do |t|
    t.integer "language_id", null: false
    t.integer "language_lesson_id", null: false
    t.integer "language_lesson_member_id", null: false
    t.string "role"
    t.text "body"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.bigint "user_id"
    t.index ["language_id"], name: "index_language_lesson_member_messages_on_language_id"
    t.index ["language_lesson_id"], name: "index_language_lesson_member_messages_on_language_lesson_id"
    t.index ["language_lesson_member_id"], name: "idx_on_language_lesson_member_id_fe254654e9"
    t.index ["user_id"], name: "index_language_lesson_member_messages_on_user_id"
  end

  create_table "language_lesson_members", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.bigint "lesson_id", null: false
    t.datetime "updated_at", precision: nil, null: false
    t.string "state"
    t.bigint "language_id", null: false
    t.datetime "created_at", null: false
    t.integer "language_member_id", null: false
    t.string "openai_thread_id"
    t.integer "messages_count", default: 0
    t.index ["language_member_id"], name: "index_language_lesson_members_on_language_member_id"
    t.index ["lesson_id"], name: "user_finished_lessons_language_module_lesson_id_index"
    t.index ["user_id", "lesson_id"], name: "user_finished_lessons_user_id_language_module_lesson_id_index", unique: true
    t.index ["user_id"], name: "user_finished_lessons_user_id_index"
  end

  create_table "language_lesson_version_infos", force: :cascade do |t|
    t.string "name"
    t.string "description"
    t.string "locale"
    t.string "theory"
    t.string "tips"
    t.string "definitions"
    t.string "instructions"
    t.bigint "language_id", null: false
    t.bigint "language_version_id", null: false
    t.bigint "version_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.bigint "language_lesson_id", null: false
    t.index ["language_id"], name: "index_language_lesson_version_infos_on_language_id"
    t.index ["language_lesson_id"], name: "index_language_lesson_version_infos_on_language_lesson_id"
    t.index ["language_version_id"], name: "index_language_lesson_version_infos_on_language_version_id"
  end

  create_table "language_lesson_versions", force: :cascade do |t|
    t.integer "order"
    t.string "original_code"
    t.string "prepared_code"
    t.string "test_code"
    t.integer "natural_order"
    t.string "path_to_code"
    t.bigint "language_version_id", null: false
    t.bigint "language_id", null: false
    t.bigint "lesson_id", null: false
    t.bigint "module_version_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["language_id"], name: "index_language_lesson_versions_on_language_id"
    t.index ["language_version_id"], name: "index_language_lesson_versions_on_language_version_id"
    t.index ["lesson_id"], name: "index_language_lesson_versions_on_lesson_id"
    t.index ["module_version_id"], name: "index_language_lesson_versions_on_module_version_id"
  end

  create_table "language_lessons", force: :cascade do |t|
    t.string "slug", limit: 255
    t.string "state", limit: 255
    t.integer "order"
    t.text "original_code"
    t.text "prepared_code"
    t.text "test_code"
    t.string "path_to_code", limit: 255
    t.bigint "module_id"
    t.bigint "language_id"
    t.bigint "upload_id"
    t.datetime "updated_at", precision: nil, null: false
    t.integer "natural_order"
    t.datetime "created_at", null: false
    t.index ["language_id", "slug"], name: "index_language_lessons_on_language_id_and_slug", unique: true
    t.index ["language_id"], name: "language_module_lessons_language_id_index"
    t.index ["module_id"], name: "language_module_lessons_module_id_index"
    t.index ["upload_id"], name: "language_module_lessons_upload_id_index"
  end

  create_table "language_members", force: :cascade do |t|
    t.bigint "language_id", null: false
    t.bigint "user_id", null: false
    t.string "state"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.integer "finished_lessons_count", default: 0, null: false
    t.index ["language_id"], name: "index_language_members_on_language_id"
    t.index ["user_id"], name: "index_language_members_on_user_id"
  end

  create_table "language_module_descriptions", force: :cascade do |t|
    t.string "name", limit: 255
    t.text "description"
    t.string "locale", limit: 255
    t.bigint "module_id"
    t.datetime "inserted_at", precision: nil, null: false
    t.datetime "updated_at", precision: nil, null: false
    t.bigint "language_id"
    t.index ["module_id"], name: "language_module_descriptions_module_id_index"
  end

  create_table "language_module_version_infos", force: :cascade do |t|
    t.string "name"
    t.string "description"
    t.string "locale"
    t.bigint "language_id", null: false
    t.bigint "version_id", null: false
    t.bigint "language_version_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["language_id"], name: "index_language_module_version_infos_on_language_id"
    t.index ["language_version_id"], name: "index_language_module_version_infos_on_language_version_id"
  end

  create_table "language_module_versions", force: :cascade do |t|
    t.bigint "language_id", null: false
    t.bigint "language_version_id", null: false
    t.bigint "module_id", null: false
    t.integer "order"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["language_id"], name: "index_language_module_versions_on_language_id"
    t.index ["language_version_id"], name: "index_language_module_versions_on_language_version_id"
    t.index ["module_id"], name: "index_language_module_versions_on_module_id"
  end

  create_table "language_modules", force: :cascade do |t|
    t.string "slug", limit: 255
    t.string "state", limit: 255
    t.integer "order"
    t.bigint "language_id"
    t.bigint "upload_id"
    t.datetime "updated_at", precision: nil, null: false
    t.datetime "created_at", null: false
    t.index ["language_id"], name: "language_modules_language_id_index"
    t.index ["upload_id"], name: "language_modules_upload_id_index"
  end

  create_table "language_version_infos", force: :cascade do |t|
    t.bigint "language_id", null: false
    t.bigint "language_version_id", null: false
    t.string "locale"
    t.string "description"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.string "title"
    t.text "seo_description"
    t.string "header"
    t.string "keywords"
    t.index ["language_id"], name: "index_language_version_infos_on_language_id"
    t.index ["language_version_id"], name: "index_language_version_infos_on_language_version_id"
  end

  create_table "language_versions", force: :cascade do |t|
    t.string "docker_image"
    t.string "exercise_filename"
    t.string "exercise_test_filename"
    t.string "extension"
    t.string "name"
    t.string "state"
    t.string "result"
    t.bigint "language_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.string "learn_as"
    t.string "progress"
    t.index ["language_id"], name: "index_language_versions_on_language_id"
  end

  create_table "languages", force: :cascade do |t|
    t.string "name", limit: 255
    t.string "slug", limit: 255
    t.string "extension", limit: 255
    t.string "docker_image", limit: 255
    t.string "exercise_filename", limit: 255
    t.string "exercise_test_filename", limit: 255
    t.string "state", limit: 255
    t.bigint "upload_id"
    t.datetime "updated_at", precision: nil, null: false
    t.string "progress"
    t.bigint "current_version_id"
    t.datetime "created_at", null: false
    t.string "learn_as"
    t.integer "lessons_count", default: 0, null: false
    t.integer "members_count", default: 0, null: false
    t.integer "order"
    t.bigint "category_id"
    t.string "openai_assistant_id"
    t.index ["category_id"], name: "index_languages_on_category_id"
    t.index ["current_version_id"], name: "index_languages_on_current_version_id"
    t.index ["slug"], name: "languages_slug_index", unique: true
    t.index ["upload_id"], name: "languages_upload_id_index"
  end

  create_table "leads", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.string "state"
    t.string "email"
    t.string "phone"
    t.string "telegram"
    t.string "whatsapp"
    t.text "survey_answers_data"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["user_id"], name: "index_leads_on_user_id"
  end

  create_table "reviews", force: :cascade do |t|
    t.bigint "language_id", null: false
    t.bigint "user_id", null: false
    t.string "state"
    t.text "body"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.string "locale"
    t.string "first_name"
    t.string "last_name"
    t.boolean "pinned"
    t.index ["language_id"], name: "index_reviews_on_language_id"
    t.index ["user_id"], name: "index_reviews_on_user_id"
  end

  create_table "solid_cable_messages", force: :cascade do |t|
    t.binary "channel", null: false
    t.binary "payload", null: false
    t.datetime "created_at", null: false
    t.bigint "channel_hash", null: false
    t.index ["channel"], name: "index_solid_cable_messages_on_channel"
    t.index ["channel_hash"], name: "index_solid_cable_messages_on_channel_hash"
    t.index ["created_at"], name: "index_solid_cable_messages_on_created_at"
  end

  create_table "solid_cache_entries", force: :cascade do |t|
    t.binary "key", null: false
    t.binary "value", null: false
    t.datetime "created_at", null: false
    t.bigint "key_hash", null: false
    t.integer "byte_size", null: false
    t.index ["byte_size"], name: "index_solid_cache_entries_on_byte_size"
    t.index ["key_hash", "byte_size"], name: "index_solid_cache_entries_on_key_hash_and_byte_size"
    t.index ["key_hash"], name: "index_solid_cache_entries_on_key_hash", unique: true
  end

  create_table "solid_queue_blocked_executions", force: :cascade do |t|
    t.bigint "job_id", null: false
    t.string "queue_name", null: false
    t.integer "priority", default: 0, null: false
    t.string "concurrency_key", null: false
    t.datetime "expires_at", null: false
    t.datetime "created_at", null: false
    t.index ["concurrency_key", "priority", "job_id"], name: "index_solid_queue_blocked_executions_for_release"
    t.index ["expires_at", "concurrency_key"], name: "index_solid_queue_blocked_executions_for_maintenance"
    t.index ["job_id"], name: "index_solid_queue_blocked_executions_on_job_id", unique: true
  end

  create_table "solid_queue_claimed_executions", force: :cascade do |t|
    t.bigint "job_id", null: false
    t.bigint "process_id"
    t.datetime "created_at", null: false
    t.index ["job_id"], name: "index_solid_queue_claimed_executions_on_job_id", unique: true
    t.index ["process_id", "job_id"], name: "index_solid_queue_claimed_executions_on_process_id_and_job_id"
  end

  create_table "solid_queue_failed_executions", force: :cascade do |t|
    t.bigint "job_id", null: false
    t.text "error"
    t.datetime "created_at", null: false
    t.index ["job_id"], name: "index_solid_queue_failed_executions_on_job_id", unique: true
  end

  create_table "solid_queue_jobs", force: :cascade do |t|
    t.string "queue_name", null: false
    t.string "class_name", null: false
    t.text "arguments"
    t.integer "priority", default: 0, null: false
    t.string "active_job_id"
    t.datetime "scheduled_at"
    t.datetime "finished_at"
    t.string "concurrency_key"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["active_job_id"], name: "index_solid_queue_jobs_on_active_job_id"
    t.index ["class_name"], name: "index_solid_queue_jobs_on_class_name"
    t.index ["finished_at"], name: "index_solid_queue_jobs_on_finished_at"
    t.index ["queue_name", "finished_at"], name: "index_solid_queue_jobs_for_filtering"
    t.index ["scheduled_at", "finished_at"], name: "index_solid_queue_jobs_for_alerting"
  end

  create_table "solid_queue_pauses", force: :cascade do |t|
    t.string "queue_name", null: false
    t.datetime "created_at", null: false
    t.index ["queue_name"], name: "index_solid_queue_pauses_on_queue_name", unique: true
  end

  create_table "solid_queue_processes", force: :cascade do |t|
    t.string "kind", null: false
    t.datetime "last_heartbeat_at", null: false
    t.bigint "supervisor_id"
    t.integer "pid", null: false
    t.string "hostname"
    t.text "metadata"
    t.datetime "created_at", null: false
    t.string "name", null: false
    t.index ["last_heartbeat_at"], name: "index_solid_queue_processes_on_last_heartbeat_at"
    t.index ["name", "supervisor_id"], name: "index_solid_queue_processes_on_name_and_supervisor_id", unique: true
    t.index ["supervisor_id"], name: "index_solid_queue_processes_on_supervisor_id"
  end

  create_table "solid_queue_ready_executions", force: :cascade do |t|
    t.bigint "job_id", null: false
    t.string "queue_name", null: false
    t.integer "priority", default: 0, null: false
    t.datetime "created_at", null: false
    t.index ["job_id"], name: "index_solid_queue_ready_executions_on_job_id", unique: true
    t.index ["priority", "job_id"], name: "index_solid_queue_poll_all"
    t.index ["queue_name", "priority", "job_id"], name: "index_solid_queue_poll_by_queue"
  end

  create_table "solid_queue_recurring_executions", force: :cascade do |t|
    t.bigint "job_id", null: false
    t.string "task_key", null: false
    t.datetime "run_at", null: false
    t.datetime "created_at", null: false
    t.index ["job_id"], name: "index_solid_queue_recurring_executions_on_job_id", unique: true
    t.index ["task_key", "run_at"], name: "index_solid_queue_recurring_executions_on_task_key_and_run_at", unique: true
  end

  create_table "solid_queue_recurring_tasks", force: :cascade do |t|
    t.string "key", null: false
    t.string "schedule", null: false
    t.string "command", limit: 2048
    t.string "class_name"
    t.text "arguments"
    t.string "queue_name"
    t.integer "priority", default: 0
    t.boolean "static", default: true, null: false
    t.text "description"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["key"], name: "index_solid_queue_recurring_tasks_on_key", unique: true
    t.index ["static"], name: "index_solid_queue_recurring_tasks_on_static"
  end

  create_table "solid_queue_scheduled_executions", force: :cascade do |t|
    t.bigint "job_id", null: false
    t.string "queue_name", null: false
    t.integer "priority", default: 0, null: false
    t.datetime "scheduled_at", null: false
    t.datetime "created_at", null: false
    t.index ["job_id"], name: "index_solid_queue_scheduled_executions_on_job_id", unique: true
    t.index ["scheduled_at", "priority", "job_id"], name: "index_solid_queue_dispatch_all"
  end

  create_table "solid_queue_semaphores", force: :cascade do |t|
    t.string "key", null: false
    t.integer "value", default: 1, null: false
    t.datetime "expires_at", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["expires_at"], name: "index_solid_queue_semaphores_on_expires_at"
    t.index ["key", "value"], name: "index_solid_queue_semaphores_on_key_and_value"
    t.index ["key"], name: "index_solid_queue_semaphores_on_key", unique: true
  end

  create_table "survey_answers", force: :cascade do |t|
    t.integer "survey_id", null: false
    t.integer "survey_item_id"
    t.integer "user_id", null: false
    t.string "state"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["survey_id", "user_id"], name: "index_survey_answers_on_survey_id_and_user_id", unique: true
    t.index ["survey_id"], name: "index_survey_answers_on_survey_id"
    t.index ["survey_item_id"], name: "index_survey_answers_on_survey_item_id"
    t.index ["user_id"], name: "index_survey_answers_on_user_id"
  end

  create_table "survey_items", force: :cascade do |t|
    t.integer "survey_id", null: false
    t.integer "order", null: false
    t.string "value"
    t.string "state"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.string "slug"
    t.index ["survey_id"], name: "index_survey_items_on_survey_id"
  end

  create_table "survey_scenario_items", force: :cascade do |t|
    t.bigint "survey_id", null: false
    t.bigint "scenario_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["scenario_id"], name: "index_survey_scenario_items_on_scenario_id"
    t.index ["survey_id", "scenario_id"], name: "index_survey_scenario_items_on_survey_id_and_scenario_id", unique: true
    t.index ["survey_id"], name: "index_survey_scenario_items_on_survey_id"
  end

  create_table "survey_scenario_members", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.bigint "scenario_id", null: false
    t.string "state"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["scenario_id"], name: "index_survey_scenario_members_on_scenario_id"
    t.index ["user_id"], name: "index_survey_scenario_members_on_user_id"
  end

  create_table "survey_scenario_triggers", force: :cascade do |t|
    t.bigint "scenario_id", null: false
    t.string "event_name"
    t.integer "event_threshold_count"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["event_name", "scenario_id"], name: "index_survey_scenario_triggers_on_event_name_and_scenario_id", unique: true
    t.index ["scenario_id"], name: "index_survey_scenario_triggers_on_scenario_id"
  end

  create_table "survey_scenarios", force: :cascade do |t|
    t.bigint "survey_item_id"
    t.string "name"
    t.string "state"
    t.string "locale"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["survey_item_id"], name: "index_survey_scenarios_on_survey_item_id"
  end

  create_table "surveys", force: :cascade do |t|
    t.string "question"
    t.string "state"
    t.string "slug"
    t.string "locale"
    t.string "description"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.bigint "parent_survey_item_id"
    t.bigint "parent_survey_id"
    t.boolean "run_always", default: false
    t.integer "run_after_finishing_lessons_count", default: 0
    t.index ["parent_survey_id"], name: "index_surveys_on_parent_survey_id"
    t.index ["parent_survey_item_id"], name: "index_surveys_on_parent_survey_item_id"
    t.index ["slug", "locale"], name: "index_surveys_on_slug_and_locale", unique: true
  end

  create_table "uploads", force: :cascade do |t|
    t.string "language_name", limit: 255
    t.datetime "inserted_at", precision: nil, null: false
    t.datetime "updated_at", precision: nil, null: false
  end

  create_table "user_accounts", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.string "provider", limit: 255, null: false
    t.string "uid", limit: 255, null: false
    t.datetime "updated_at", precision: 0, null: false
    t.datetime "created_at", null: false
  end

  create_table "user_survey_pivots", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.bigint "coding_experience_item_id"
    t.bigint "goal_item_id"
    t.bigint "study_plan_item_id"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["coding_experience_item_id"], name: "index_user_survey_pivots_on_coding_experience_item_id"
    t.index ["goal_item_id"], name: "index_user_survey_pivots_on_goal_item_id"
    t.index ["study_plan_item_id"], name: "index_user_survey_pivots_on_study_plan_item_id"
    t.index ["user_id"], name: "index_user_survey_pivots_on_user_id"
  end

  create_table "users", force: :cascade do |t|
    t.string "first_name", limit: 255
    t.string "last_name", limit: 255
    t.string "email", limit: 255
    t.string "nickname", limit: 255
    t.integer "github_uid"
    t.datetime "updated_at", precision: nil, null: false
    t.string "facebook_uid", limit: 255
    t.string "password_digest", limit: 255
    t.string "confirmation_token", limit: 255
    t.string "reset_password_token", limit: 255
    t.string "state", limit: 255
    t.string "locale", limit: 255
    t.string "email_delivery_state", limit: 255
    t.boolean "admin"
    t.datetime "created_at", null: false
    t.boolean "help"
    t.integer "assistant_messages_count", default: 0
    t.string "contact_value"
    t.string "contact_method"
    t.index ["email"], name: "index_users_on_email", unique: true
  end

  add_foreign_key "active_storage_attachments", "active_storage_blobs", column: "blob_id"
  add_foreign_key "active_storage_variant_records", "active_storage_blobs", column: "blob_id"
  add_foreign_key "blog_posts", "languages"
  add_foreign_key "blog_posts", "users", column: "creator_id"
  add_foreign_key "book_requests", "users"
  add_foreign_key "event_store_events_in_streams", "event_store_events", column: "event_id", primary_key: "event_id"
  add_foreign_key "language_landing_page_qna_items", "language_landing_pages"
  add_foreign_key "language_landing_pages", "language_categories"
  add_foreign_key "language_landing_pages", "languages"
  add_foreign_key "language_lesson_member_messages", "language_lesson_members"
  add_foreign_key "language_lesson_member_messages", "language_lessons"
  add_foreign_key "language_lesson_member_messages", "languages"
  add_foreign_key "language_lesson_member_messages", "users"
  add_foreign_key "language_lesson_members", "language_lessons", column: "lesson_id"
  add_foreign_key "language_lesson_members", "language_members"
  add_foreign_key "language_lesson_members", "users"
  add_foreign_key "language_lesson_version_infos", "language_lesson_versions", column: "version_id"
  add_foreign_key "language_lesson_version_infos", "language_lessons"
  add_foreign_key "language_lesson_version_infos", "language_versions"
  add_foreign_key "language_lesson_version_infos", "languages"
  add_foreign_key "language_lesson_versions", "language_lessons", column: "lesson_id"
  add_foreign_key "language_lesson_versions", "language_module_versions", column: "module_version_id"
  add_foreign_key "language_lesson_versions", "language_versions"
  add_foreign_key "language_lesson_versions", "languages"
  add_foreign_key "language_lessons", "language_modules", column: "module_id"
  add_foreign_key "language_lessons", "languages"
  add_foreign_key "language_lessons", "uploads"
  add_foreign_key "language_members", "languages"
  add_foreign_key "language_members", "users"
  add_foreign_key "language_module_descriptions", "language_modules", column: "module_id"
  add_foreign_key "language_module_descriptions", "languages"
  add_foreign_key "language_module_version_infos", "language_module_versions", column: "version_id"
  add_foreign_key "language_module_version_infos", "language_versions"
  add_foreign_key "language_module_version_infos", "languages"
  add_foreign_key "language_module_versions", "language_modules", column: "module_id"
  add_foreign_key "language_module_versions", "language_versions"
  add_foreign_key "language_module_versions", "languages"
  add_foreign_key "language_modules", "languages"
  add_foreign_key "language_modules", "uploads"
  add_foreign_key "language_version_infos", "language_versions"
  add_foreign_key "language_version_infos", "languages"
  add_foreign_key "language_versions", "languages"
  add_foreign_key "languages", "language_categories", column: "category_id"
  add_foreign_key "languages", "language_versions", column: "current_version_id"
  add_foreign_key "languages", "uploads"
  add_foreign_key "leads", "users"
  add_foreign_key "reviews", "languages"
  add_foreign_key "reviews", "users"
  add_foreign_key "solid_queue_blocked_executions", "solid_queue_jobs", column: "job_id", on_delete: :cascade
  add_foreign_key "solid_queue_claimed_executions", "solid_queue_jobs", column: "job_id", on_delete: :cascade
  add_foreign_key "solid_queue_failed_executions", "solid_queue_jobs", column: "job_id", on_delete: :cascade
  add_foreign_key "solid_queue_ready_executions", "solid_queue_jobs", column: "job_id", on_delete: :cascade
  add_foreign_key "solid_queue_recurring_executions", "solid_queue_jobs", column: "job_id", on_delete: :cascade
  add_foreign_key "solid_queue_scheduled_executions", "solid_queue_jobs", column: "job_id", on_delete: :cascade
  add_foreign_key "survey_answers", "survey_items"
  add_foreign_key "survey_answers", "surveys"
  add_foreign_key "survey_answers", "users"
  add_foreign_key "survey_items", "surveys"
  add_foreign_key "survey_scenario_items", "survey_scenarios", column: "scenario_id"
  add_foreign_key "survey_scenario_items", "surveys"
  add_foreign_key "survey_scenario_members", "survey_scenarios", column: "scenario_id"
  add_foreign_key "survey_scenario_members", "users"
  add_foreign_key "survey_scenario_triggers", "survey_scenarios", column: "scenario_id"
  add_foreign_key "survey_scenarios", "survey_items"
  add_foreign_key "surveys", "survey_items", column: "parent_survey_item_id"
  add_foreign_key "surveys", "surveys", column: "parent_survey_id"
  add_foreign_key "user_accounts", "users"
  add_foreign_key "user_survey_pivots", "survey_items", column: "coding_experience_item_id"
  add_foreign_key "user_survey_pivots", "survey_items", column: "goal_item_id"
  add_foreign_key "user_survey_pivots", "survey_items", column: "study_plan_item_id"
  add_foreign_key "user_survey_pivots", "users"
end
