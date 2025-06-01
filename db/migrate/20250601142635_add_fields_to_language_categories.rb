class AddFieldsToLanguageCategories < ActiveRecord::Migration[8.0]
  def change
    add_column :language_categories, :header, :string, null: true
    add_column :language_categories, :description, :string, null: true
  end
end
