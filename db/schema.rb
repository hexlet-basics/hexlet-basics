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

ActiveRecord::Schema.define(version: 2020_05_30_225230) do

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

  create_table "language_module_lesson_exercises", force: :cascade do |t|
    t.string "original_code"
    t.string "prepared_code"
    t.string "test_code"
    t.string "path_to_code"
    t.integer "lesson_id", null: false
    t.integer "language_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_module_lesson_exercises_on_language_id"
    t.index ["lesson_id"], name: "index_language_module_lesson_exercise_on_lesson_id"
  end

  create_table "language_module_lesson_members", force: :cascade do |t|
    t.string "state"
    t.integer "user_id", null: false
    t.integer "lesson_id", null: false
    t.integer "language_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["language_id"], name: "index_language_module_lesson_members_on_language_id"
    t.index ["lesson_id"], name: "index_language_module_lesson_member_on_lesson_id"
    t.index ["user_id"], name: "index_language_module_lesson_members_on_user_id"
  end

  create_table "language_module_lessons", force: :cascade do |t|
    t.integer "order"
    t.string "slug"
    t.integer "language_id", null: false
    t.integer "module_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.string "state"
    t.integer "upload_id", null: false
    t.integer "current_exercise_id"
    t.index ["current_exercise_id"], name: "index_language_module_lessons_on_current_exercise_id"
    t.index ["language_id"], name: "index_language_module_lessons_on_language_id"
    t.index ["module_id"], name: "index_language_module_lessons_on_module_id"
    t.index ["upload_id"], name: "index_language_module_lessons_on_upload_id"
  end

  create_table "language_modules", force: :cascade do |t|
    t.string "slug"
    t.integer "language_id", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.integer "order"
    t.integer "upload_id", null: false
    t.index ["language_id"], name: "index_language_modules_on_language_id"
    t.index ["upload_id"], name: "index_language_modules_on_upload_id"
  end

  create_table "languages", force: :cascade do |t|
    t.string "name"
    t.string "slug"
    t.string "extension"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.string "docker_image"
    t.string "exercise_filename"
    t.string "exercise_test_filename"
    t.integer "upload_id", null: false
    t.index ["upload_id"], name: "index_languages_on_upload_id"
  end

  create_table "uploads", force: :cascade do |t|
    t.string "language_name"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

  create_table "users", force: :cascade do |t|
    t.string "email", null: false
    t.string "password_digest", null: false
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["email"], name: "index_users_on_email", unique: true
  end

  create_table "version_associations", force: :cascade do |t|
    t.integer "version_id"
    t.string "foreign_key_name", null: false
    t.integer "foreign_key_id"
    t.string "foreign_type"
    t.index ["foreign_key_name", "foreign_key_id", "foreign_type"], name: "index_version_associations_on_foreign_key"
    t.index ["version_id"], name: "index_version_associations_on_version_id"
  end

  create_table "versions", force: :cascade do |t|
    t.string "item_type", null: false
    t.integer "item_id", limit: 8, null: false
    t.string "event", null: false
    t.string "whodunnit"
    t.text "object", limit: 1073741823
    t.datetime "created_at"
    t.integer "transaction_id"
    t.index ["item_type", "item_id"], name: "index_versions_on_item_type_and_item_id"
    t.index ["transaction_id"], name: "index_versions_on_transaction_id"
  end

  add_foreign_key "language_module_descriptions", "language_modules", column: "module_id"
  add_foreign_key "language_module_descriptions", "languages"
  add_foreign_key "language_module_lesson_descriptions", "language_module_lessons", column: "lesson_id"
  add_foreign_key "language_module_lesson_descriptions", "languages"
  add_foreign_key "language_module_lesson_exercises", "language_module_lessons", column: "lesson_id"
  add_foreign_key "language_module_lesson_exercises", "languages"
  add_foreign_key "language_module_lesson_members", "language_module_lessons", column: "lesson_id"
  add_foreign_key "language_module_lesson_members", "languages"
  add_foreign_key "language_module_lesson_members", "users"
  add_foreign_key "language_module_lessons", "language_modules", column: "module_id"
  add_foreign_key "language_module_lessons", "languages"
  add_foreign_key "language_module_lessons", "uploads"
  add_foreign_key "language_modules", "languages"
  add_foreign_key "language_modules", "uploads"
  add_foreign_key "languages", "uploads"
end
