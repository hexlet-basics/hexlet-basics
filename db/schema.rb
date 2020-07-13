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

ActiveRecord::Schema.define(version: 2020_07_09_171856) do

  create_table "language_lesson_members", force: :cascade do |t|
    t.integer "language_id", null: false
    t.integer "lesson_id", null: false
    t.integer "lesson_version_id", null: false
    t.integer "user_id", null: false
    t.string "state"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_lesson_members_on_language_id"
    t.index ["lesson_id"], name: "index_language_lesson_members_on_lesson_id"
    t.index ["lesson_version_id"], name: "index_language_lesson_members_on_lesson_version_id"
    t.index ["user_id"], name: "index_language_lesson_members_on_user_id"
  end

  create_table "language_lesson_version_infos", force: :cascade do |t|
    t.string "name"
    t.string "description"
    t.string "locale"
    t.string "theory"
    t.string "tips"
    t.string "definitions"
    t.string "instructions"
    t.integer "language_id", null: false
    t.integer "language_version_id", null: false
    t.integer "version_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_lesson_version_infos_on_language_id"
    t.index ["language_version_id"], name: "index_language_lesson_version_infos_on_language_version_id"
  end

  create_table "language_lesson_versions", force: :cascade do |t|
    t.string "order"
    t.string "original_code"
    t.string "prepared_code"
    t.string "test_code"
    t.string "path_to_code"
    t.integer "language_version_id", null: false
    t.integer "language_id", null: false
    t.integer "lesson_id", null: false
    t.integer "module_version_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_lesson_versions_on_language_id"
    t.index ["language_version_id"], name: "index_language_lesson_versions_on_language_version_id"
    t.index ["lesson_id"], name: "index_language_lesson_versions_on_lesson_id"
    t.index ["module_version_id"], name: "index_language_lesson_versions_on_module_version_id"
  end

  create_table "language_lessons", force: :cascade do |t|
    t.string "slug"
    t.integer "language_id", null: false
    t.integer "module_id", null: false
    t.integer "current_version_id"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["current_version_id"], name: "index_language_lessons_on_current_version_id"
    t.index ["language_id"], name: "index_language_lessons_on_language_id"
    t.index ["module_id"], name: "index_language_lessons_on_module_id"
  end

  create_table "language_module_version_infos", force: :cascade do |t|
    t.string "name"
    t.string "description"
    t.string "locale"
    t.integer "language_id", null: false
    t.integer "version_id", null: false
    t.integer "language_version_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_module_version_infos_on_language_id"
    t.index ["language_version_id"], name: "index_language_module_version_infos_on_language_version_id"
  end

  create_table "language_module_versions", force: :cascade do |t|
    t.integer "language_id", null: false
    t.integer "language_version_id", null: false
    t.integer "module_id", null: false
    t.string "order"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_module_versions_on_language_id"
    t.index ["language_version_id"], name: "index_language_module_versions_on_language_version_id"
    t.index ["module_id"], name: "index_language_module_versions_on_module_id"
  end

  create_table "language_modules", force: :cascade do |t|
    t.string "slug"
    t.integer "language_id", null: false
    t.integer "current_version_id"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["current_version_id"], name: "index_language_modules_on_current_version_id"
    t.index ["language_id"], name: "index_language_modules_on_language_id"
  end

  create_table "language_versions", force: :cascade do |t|
    t.string "docker_image"
    t.string "exercise_filename"
    t.string "exercise_test_filename"
    t.string "extension"
    t.string "name"
    t.string "state"
    t.string "result"
    t.integer "language_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_versions_on_language_id"
  end

  create_table "languages", force: :cascade do |t|
    t.string "slug"
    t.integer "current_version_id"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["current_version_id"], name: "index_languages_on_current_version_id"
  end

  create_table "users", force: :cascade do |t|
    t.string "email", null: false
    t.string "password_digest", null: false
    t.boolean "admin", default: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["email"], name: "index_users_on_email", unique: true
  end

  add_foreign_key "language_lesson_members", "language_lesson_versions", column: "lesson_version_id"
  add_foreign_key "language_lesson_members", "language_lessons", column: "lesson_id"
  add_foreign_key "language_lesson_members", "languages"
  add_foreign_key "language_lesson_members", "users"
  add_foreign_key "language_lesson_version_infos", "language_lesson_versions", column: "version_id"
  add_foreign_key "language_lesson_version_infos", "language_versions"
  add_foreign_key "language_lesson_version_infos", "languages"
  add_foreign_key "language_lesson_versions", "language_lessons", column: "lesson_id"
  add_foreign_key "language_lesson_versions", "language_module_versions", column: "module_version_id"
  add_foreign_key "language_lesson_versions", "language_versions"
  add_foreign_key "language_lesson_versions", "languages"
  add_foreign_key "language_lessons", "language_modules", column: "module_id"
  add_foreign_key "language_lessons", "languages"
  add_foreign_key "language_module_version_infos", "language_module_versions", column: "version_id"
  add_foreign_key "language_module_version_infos", "language_versions"
  add_foreign_key "language_module_version_infos", "languages"
  add_foreign_key "language_module_versions", "language_modules", column: "module_id"
  add_foreign_key "language_module_versions", "language_versions"
  add_foreign_key "language_module_versions", "languages"
  add_foreign_key "language_modules", "languages"
  add_foreign_key "language_versions", "languages"
end
