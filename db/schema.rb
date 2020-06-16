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

ActiveRecord::Schema.define(version: 2020_06_20_214520) do

  create_table "language_module_descriptions", force: :cascade do |t|
    t.string "name"
    t.text "description"
    t.string "locale"
    t.integer "module_id", null: false
    t.integer "language_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_module_descriptions_on_language_id"
    t.index ["module_id"], name: "index_language_module_descriptions_on_module_id"
  end

  create_table "language_module_lesson_descriptions", force: :cascade do |t|
    t.string "instructions"
    t.string "locale"
    t.string "name"
    t.string "theory"
    t.string "tips"
    t.string "definitions"
    t.integer "lesson_id", null: false
    t.integer "language_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_module_lesson_descriptions_on_language_id"
    t.index ["lesson_id"], name: "index_language_module_lesson_descriptions_on_lesson_id"
  end

  create_table "language_module_lesson_versions", force: :cascade do |t|
    t.integer "language_version_id", null: false
    t.integer "module_version_id", null: false
    t.integer "lesson_id", null: false
    t.string "order"
    t.string "original_code"
    t.string "prepared_code"
    t.string "test_code"
    t.string "path_to_code"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.integer "upload_id"
    t.index ["language_version_id"], name: "index_language_module_lesson_versions_on_language_version_id"
    t.index ["lesson_id"], name: "index_language_module_lesson_version_on_lesson_id"
    t.index ["module_version_id"], name: "index_language_module_lesson_version_on_module_version_id"
    t.index ["upload_id"], name: "index_language_module_lesson_versions_on_upload_id"
  end

  create_table "language_module_lessons", force: :cascade do |t|
    t.string "slug"
    t.integer "language_id", null: false
    t.integer "module_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.integer "current_version_id"
    t.index ["current_version_id"], name: "index_language_module_lessons_on_current_version_id"
    t.index ["language_id"], name: "index_language_module_lessons_on_language_id"
    t.index ["module_id"], name: "index_language_module_lessons_on_module_id"
  end

  create_table "language_module_versions", force: :cascade do |t|
    t.integer "language_version_id", null: false
    t.integer "module_id", null: false
    t.string "order"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.integer "upload_id"
    t.index ["language_version_id"], name: "index_language_module_versions_on_language_version_id"
    t.index ["module_id"], name: "index_language_module_versions_on_module_id"
    t.index ["upload_id"], name: "index_language_module_versions_on_upload_id"
  end

  create_table "language_modules", force: :cascade do |t|
    t.string "slug"
    t.integer "language_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.integer "current_version_id"
    t.index ["current_version_id"], name: "index_language_modules_on_current_version_id"
    t.index ["language_id"], name: "index_language_modules_on_language_id"
  end

  create_table "language_uploads", force: :cascade do |t|
    t.string "state"
    t.string "uploader"
    t.integer "language_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_uploads_on_language_id"
  end

  create_table "language_versions", force: :cascade do |t|
    t.string "docker_image"
    t.string "exercise_filename"
    t.string "exercise_test_filename"
    t.string "extension"
    t.string "name"
    t.integer "language_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.integer "upload_id"
    t.index ["language_id"], name: "index_language_versions_on_language_id"
    t.index ["upload_id"], name: "index_language_versions_on_upload_id"
  end

  create_table "languages", force: :cascade do |t|
    t.string "slug"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.integer "current_version_id"
    t.index ["current_version_id"], name: "index_languages_on_current_version_id"
  end

  create_table "users", force: :cascade do |t|
    t.string "email", null: false
    t.string "password_digest", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.boolean "admin", default: false
    t.index ["email"], name: "index_users_on_email", unique: true
  end

  add_foreign_key "language_module_descriptions", "language_modules", column: "module_id"
  add_foreign_key "language_module_descriptions", "languages"
  add_foreign_key "language_module_lesson_descriptions", "language_module_lessons", column: "lesson_id"
  add_foreign_key "language_module_lesson_descriptions", "languages"
  add_foreign_key "language_module_lesson_versions", "language_module_lessons", column: "lesson_id"
  add_foreign_key "language_module_lesson_versions", "language_module_versions", column: "module_version_id"
  add_foreign_key "language_module_lesson_versions", "language_uploads", column: "upload_id"
  add_foreign_key "language_module_lesson_versions", "language_versions"
  add_foreign_key "language_module_lessons", "language_modules", column: "module_id"
  add_foreign_key "language_module_lessons", "languages"
  add_foreign_key "language_module_versions", "language_modules", column: "module_id"
  add_foreign_key "language_module_versions", "language_uploads", column: "upload_id"
  add_foreign_key "language_module_versions", "language_versions"
  add_foreign_key "language_modules", "languages"
  add_foreign_key "language_uploads", "languages"
  add_foreign_key "language_versions", "language_uploads", column: "upload_id"
  add_foreign_key "language_versions", "languages"
end
