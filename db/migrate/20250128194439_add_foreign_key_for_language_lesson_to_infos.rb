# frozen_string_literal: true

class AddForeignKeyForLanguageLessonToInfos < ActiveRecord::Migration[7.2]
  def change
    add_foreign_key :language_lesson_version_infos, :language_lessons, column: :language_lesson_id, validate: false
  end
end
