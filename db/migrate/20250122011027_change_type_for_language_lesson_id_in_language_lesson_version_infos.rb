class ChangeTypeForLanguageLessonIdInLanguageLessonVersionInfos < ActiveRecord::Migration[8.0]
  def change
    change_column :language_lesson_version_infos, :language_lesson_id, :bigint, null: false
  end
end
