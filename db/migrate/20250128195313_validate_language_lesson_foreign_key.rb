# frozen_string_literal: true

class ValidateLanguageLessonForeignKey < ActiveRecord::Migration[8.0]
  def change
    validate_foreign_key :language_lesson_version_infos, column: :language_lesson_id
  end
end
