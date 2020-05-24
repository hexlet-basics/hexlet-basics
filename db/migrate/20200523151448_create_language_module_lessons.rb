class CreateLanguageModuleLessons < ActiveRecord::Migration[6.0]
  def change
    create_table :language_module_lessons do |t|
      t.string :slug
      t.integer :order, null: false
      t.string :path_to_code, null: false
      t.references :language, null: false, foreign_key: true
      t.references :language_module, null: false, foreign_key: true

      t.timestamps
    end
  end
end
