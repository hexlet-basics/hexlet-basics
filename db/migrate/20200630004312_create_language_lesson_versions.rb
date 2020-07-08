class CreateLanguageLessonVersions < ActiveRecord::Migration[6.0]
  def change
    create_table :language_lesson_versions do |t|
      t.string :order
      t.string :original_code
      t.string :prepared_code
      t.string :test_code
      t.string :path_to_code
      t.references :language_version, null: false, foreign_key: true
      t.references :language, null: false, foreign_key: true
      t.references :language_lesson, null: false, foreign_key: true
      t.references :language_module_version, null: false, foreign_key: true


      t.timestamps
    end
  end
end
