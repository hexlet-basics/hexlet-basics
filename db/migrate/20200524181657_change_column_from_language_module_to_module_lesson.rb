class ChangeColumnFromLanguageModuleToModuleLesson < ActiveRecord::Migration[6.0]
  def change
    rename_column :language_module_lessons, :language_module_id, :module_id
  end
end
