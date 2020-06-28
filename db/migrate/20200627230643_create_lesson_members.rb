class CreateLessonMembers < ActiveRecord::Migration[6.0]
  def change
    create_table :lesson_members do |t|
      t.references :user, null: false, foreign_key: true
      t.references :language_module_lesson, null: false, foreign_key: true
      t.string :state
      t.references :language_module_lesson_version, null: false, foreign_key: true

      t.timestamps
    end
  end
end
