class AddLanguageCategoryToLanguage < ActiveRecord::Migration[7.0]
  def change
    add_reference :languages, :category, null: true, foreign_key: { to_table: :language_categories }
  end
end
