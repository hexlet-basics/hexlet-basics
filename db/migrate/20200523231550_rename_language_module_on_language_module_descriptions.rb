class RenameLanguageModuleOnLanguageModuleDescriptions < ActiveRecord::Migration[6.0]
  def change
    rename_column :language_module_descriptions, :language_module_id, :module_id
  end
end
