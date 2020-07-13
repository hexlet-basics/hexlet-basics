class RenameLessonMemberColumns < ActiveRecord::Migration[6.0]
  def change
    rename_column :language_lesson_members, :language_lesson_id, :lesson_id
    rename_column :language_lesson_members, :language_lesson_version_id, :lesson_version_id
  end
end
