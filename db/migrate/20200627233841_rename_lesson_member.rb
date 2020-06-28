class RenameLessonMember < ActiveRecord::Migration[6.0]
  def change
    rename_column :lesson_members, :language_module_lesson_id, :lesson_id
    rename_column :lesson_members, :language_module_lesson_version_id, :lesson_version_id
  end
end
