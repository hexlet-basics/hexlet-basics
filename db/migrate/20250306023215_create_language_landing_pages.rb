class CreateLanguageLandingPages < ActiveRecord::Migration[8.0]
  def change
    create_table :language_landing_pages do |t|
      t.references :language, null: false, foreign_key: true
      t.references :language_category, null: true, foreign_key: true
      t.string :meta_title
      t.string :locale
      t.string :header
      t.string :slug
      t.string :order
      t.boolean :main
      t.string :state
      t.string :description
      t.string :meta_description

      t.timestamps
    end
  end
end
