# frozen_string_literal: true

class AddIndexForLanguageLessonToInfos < ActiveRecord::Migration[7.2]
  disable_ddl_transaction!

  def change
    add_index :language_lesson_version_infos, :language_lesson_id, algorithm: :concurrently
  end
end
