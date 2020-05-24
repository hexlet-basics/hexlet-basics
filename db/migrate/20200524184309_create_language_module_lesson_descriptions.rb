class CreateLanguageModuleLessonDescriptions < ActiveRecord::Migration[6.0]
  def change
    create_table :language_module_lesson_descriptions do |t|
      t.string :instructions
      t.string :locale
      t.string :name
      t.string :theory
      t.string :tips
      t.string :definitions
      t.references :language_module_lesson, null: false, foreign_key: true, index: { name: :index_language_module_lesson_descriptions_on_lesson_id }
      t.references :language, null: false, foreign_key: true

      t.timestamps
    end
  end
end
