class CreateLanguageModuleDescriptions < ActiveRecord::Migration[6.0]
  def change
    create_table :language_module_descriptions do |t|
      t.string :name
      t.text :description
      t.string :locale
      t.references :language_module, null: false, foreign_key: true
      t.references :language, null: false, foreign_key: true

      t.timestamps
    end
  end
end
