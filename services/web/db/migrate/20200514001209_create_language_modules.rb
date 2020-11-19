class CreateLanguageModules < ActiveRecord::Migration[6.0]
  def change
    create_table :language_modules do |t|
      t.string :slug
      t.references :language, null: false, foreign_key: true

      t.timestamps
    end
  end
end
