class RenameLanguageModuleLessonVersionField < ActiveRecord::Migration[6.0]
  def change
    rename_column :language_module_lesson_versions, :language_module_lesson_id, :lesson_id
  end
end
