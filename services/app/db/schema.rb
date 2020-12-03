# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `rails
# db:schema:load`. When creating a new database, `rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 2020_12_03_110336) do

  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "language_lesson_members", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.bigint "lesson_id", null: false
    t.datetime "updated_at", null: false
    t.string "state"
    t.bigint "language_id", null: false
    t.datetime "created_at", precision: 6, null: false
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
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_lesson_version_infos_on_language_id"
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
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
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
    t.datetime "inserted_at", null: false
    t.datetime "updated_at", null: false
    t.integer "natural_order"
    t.index ["language_id"], name: "language_module_lessons_language_id_index"
    t.index ["module_id"], name: "language_module_lessons_module_id_index"
    t.index ["upload_id"], name: "language_module_lessons_upload_id_index"
  end

  create_table "language_members", force: :cascade do |t|
    t.bigint "language_id", null: false
    t.bigint "user_id", null: false
    t.string "state"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_members_on_language_id"
    t.index ["user_id"], name: "index_language_members_on_user_id"
  end

  create_table "language_module_descriptions", force: :cascade do |t|
    t.string "name", limit: 255
    t.text "description"
    t.string "locale", limit: 255
    t.bigint "module_id"
    t.datetime "inserted_at", null: false
    t.datetime "updated_at", null: false
    t.bigint "language_id"
    t.index ["module_id"], name: "language_module_descriptions_module_id_index"
  end

  create_table "language_module_lesson_descriptions", force: :cascade do |t|
    t.string "name", limit: 255
    t.text "theory"
    t.text "instructions"
    t.string "locale", limit: 255
    t.string "tips", limit: 255, null: false, array: true
    t.bigint "lesson_id"
    t.datetime "inserted_at", null: false
    t.datetime "updated_at", null: false
    t.bigint "language_id"
    t.jsonb "definitions", null: false, array: true
    t.index ["lesson_id"], name: "language_module_lesson_descriptions_lesson_id_index"
  end

  create_table "language_module_version_infos", force: :cascade do |t|
    t.string "name"
    t.string "description"
    t.string "locale"
    t.bigint "language_id", null: false
    t.bigint "version_id", null: false
    t.bigint "language_version_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_module_version_infos_on_language_id"
    t.index ["language_version_id"], name: "index_language_module_version_infos_on_language_version_id"
  end

  create_table "language_module_versions", force: :cascade do |t|
    t.bigint "language_id", null: false
    t.bigint "language_version_id", null: false
    t.bigint "module_id", null: false
    t.integer "order"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
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
    t.datetime "updated_at", null: false
    t.datetime "created_at", precision: 6, null: false
    t.index ["language_id"], name: "language_modules_language_id_index"
    t.index ["upload_id"], name: "language_modules_upload_id_index"
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
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
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
    t.datetime "updated_at", null: false
    t.string "progress"
    t.bigint "current_version_id"
    t.datetime "created_at", precision: 6, null: false
    t.index ["current_version_id"], name: "index_languages_on_current_version_id"
    t.index ["slug"], name: "languages_slug_index", unique: true
    t.index ["upload_id"], name: "languages_upload_id_index"
  end

  create_table "uploads", force: :cascade do |t|
    t.string "language_name", limit: 255
    t.datetime "inserted_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "user_accounts", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.string "provider", limit: 255, null: false
    t.string "uid", limit: 255, null: false
    t.datetime "updated_at", precision: 0, null: false
    t.datetime "created_at", precision: 6, null: false
  end

  create_table "users", force: :cascade do |t|
    t.string "first_name", limit: 255
    t.string "last_name", limit: 255
    t.string "email", limit: 255
    t.string "nickname", limit: 255
    t.integer "github_uid"
    t.datetime "updated_at", null: false
    t.string "facebook_uid", limit: 255
    t.string "password_digest", limit: 255
    t.string "confirmation_token", limit: 255
    t.string "reset_password_token", limit: 255
    t.string "state", limit: 255
    t.string "locale", limit: 255
    t.string "email_delivery_state", limit: 255
    t.boolean "admin"
    t.datetime "created_at", precision: 6, null: false
    t.index ["email"], name: "users_email_index", unique: true
  end

  add_foreign_key "language_lesson_members", "language_lessons", column: "lesson_id", name: "user_finished_lessons_language_module_lesson_id_fkey"
  add_foreign_key "language_lesson_members", "users", name: "user_finished_lessons_user_id_fkey"
  add_foreign_key "language_lesson_version_infos", "language_lesson_versions", column: "version_id"
  add_foreign_key "language_lesson_version_infos", "language_versions"
  add_foreign_key "language_lesson_version_infos", "languages"
  add_foreign_key "language_lesson_versions", "language_lessons", column: "lesson_id"
  add_foreign_key "language_lesson_versions", "language_module_versions", column: "module_version_id"
  add_foreign_key "language_lesson_versions", "language_versions"
  add_foreign_key "language_lesson_versions", "languages"
  add_foreign_key "language_lessons", "language_modules", column: "module_id", name: "language_module_lessons_module_id_fkey"
  add_foreign_key "language_lessons", "languages", name: "language_module_lessons_language_id_fkey"
  add_foreign_key "language_lessons", "uploads", name: "language_module_lessons_upload_id_fkey"
  add_foreign_key "language_members", "languages"
  add_foreign_key "language_members", "users"
  add_foreign_key "language_module_descriptions", "language_modules", column: "module_id", name: "language_module_descriptions_module_id_fkey"
  add_foreign_key "language_module_descriptions", "languages", name: "language_module_descriptions_language_id_fkey"
  add_foreign_key "language_module_lesson_descriptions", "language_lessons", column: "lesson_id", name: "language_module_lesson_descriptions_lesson_id_fkey"
  add_foreign_key "language_module_lesson_descriptions", "languages", name: "language_module_lesson_descriptions_language_id_fkey"
  add_foreign_key "language_module_version_infos", "language_module_versions", column: "version_id"
  add_foreign_key "language_module_version_infos", "language_versions"
  add_foreign_key "language_module_version_infos", "languages"
  add_foreign_key "language_module_versions", "language_modules", column: "module_id"
  add_foreign_key "language_module_versions", "language_versions"
  add_foreign_key "language_module_versions", "languages"
  add_foreign_key "language_modules", "languages", name: "language_modules_language_id_fkey"
  add_foreign_key "language_modules", "uploads", name: "language_modules_upload_id_fkey"
  add_foreign_key "language_versions", "languages"
  add_foreign_key "languages", "language_versions", column: "current_version_id"
  add_foreign_key "languages", "uploads", name: "languages_upload_id_fkey"
  add_foreign_key "user_accounts", "users", name: "user_accounts_user_id_fkey"
end
