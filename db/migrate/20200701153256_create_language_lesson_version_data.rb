class CreateLanguageLessonVersionData < ActiveRecord::Migration[6.0]
  def change
    create_table :language_lesson_version_infos do |t|
      t.string :name
      t.string :description
      t.string :locale
      t.string :theory
      t.string :tips
      t.string :definitions
      t.string :instructions
      t.references :language, null: false, foreign_key: true
      t.references :language_version, null: false, foreign_key: true
      t.references :language_lesson_version, null: false, foreign_key: true, index: false


      t.timestamps
    end
  end
end
