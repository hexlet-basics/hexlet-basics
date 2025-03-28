class CreateLanguageLessonMemberMessages < ActiveRecord::Migration[8.0]
  def change
    create_table :language_lesson_member_messages do |t|
      t.references :language, null: false, foreign_key: true
      t.references :language_lesson, null: false, foreign_key: true
      t.references :language_lesson_member, null: false, foreign_key: true
      t.string :role
      t.text :body

      t.timestamps
    end
  end
end
