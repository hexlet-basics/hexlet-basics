class RenameColumns < ActiveRecord::Migration[6.0]
  def change
    rename_column :language_lesson_version_data, :language_lesson_version_id, :version_id
    rename_column :language_lesson_versions, :language_lesson_id, :lesson_id
    rename_column :language_lessons, :language_module_id, :module_id
    rename_column :language_module_version_data, :language_module_version_id, :version_id
    rename_column :language_lesson_versions, :language_module_version_id, :module_version_id
  end
end
