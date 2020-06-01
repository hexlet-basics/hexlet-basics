class CreateLanguageModuleLessonVersions < ActiveRecord::Migration[6.0]
  def change
    create_table :language_module_lesson_versions do |t|
      t.references :language_version, null: false, foreign_key: true
      t.references :language_module_version, null: false, foreign_key: true, index: { name: :index_language_module_lesson_version_on_module_version_id }
      t.string :order
      t.string :original_code
      t.string :prepared_code
      t.string :test_code
      t.string :path_to_code

      t.timestamps
    end
  end
end
