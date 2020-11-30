class AddCreatedAtToLanguageModuleDescriptions < ActiveRecord::Migration[6.0]
  def change
    add_column :language_modules, :created_at, :datetime, precision: 6, null: false
    Language::Module.update_all(created_at: :inserted_at)
    remove_column :language_modules, :inserted_at
  end
end
