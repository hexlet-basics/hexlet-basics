class CreateLanguageModuleLessons < ActiveRecord::Migration[6.0]
  def change
    create_table :language_module_lessons do |t|
      t.integer :order
      t.integer :natural_order
      t.string :slug
      t.string :original_code
      t.string :prepared_code
      t.string :test_code
      t.string :path_to_code
      t.references :language, null: false, foreign_key: true
      t.references :language_module, null: false, foreign_key: true

      t.timestamps
    end
  end
end
