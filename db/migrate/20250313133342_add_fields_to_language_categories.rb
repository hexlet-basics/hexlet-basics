class AddFieldsToLanguageCategories < ActiveRecord::Migration[8.0]
  def change
    add_column :language_categories, :name, :string
    add_column :language_categories, :locale, :string
  end
end
