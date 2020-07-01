class CreateLanguageLessons < ActiveRecord::Migration[6.0]
  def change
    create_table :language_lessons do |t|
      t.string :slug
      t.references :language, null: false, foreign_key: true
      t.references :language_module, null: false, foreign_key: true
      t.references :version

      t.timestamps
    end
  end
end
