class CreateLanguageCategoryItems < ActiveRecord::Migration[8.0]
  def change
    create_table :language_category_items do |t|
      t.references :language_category, null: false, foreign_key: true
      t.references :language_landing_page, null: false, foreign_key: true

      t.timestamps
    end
  end
end
