class CreateLanguageCategories < ActiveRecord::Migration[7.0]
  def change
    create_table :language_categories do |t|
      t.string :name_ru
      t.string :name_en
      t.string :slug

      t.timestamps
    end
  end
end
