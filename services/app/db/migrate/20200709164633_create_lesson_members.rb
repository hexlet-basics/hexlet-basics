class CreateLessonMembers < ActiveRecord::Migration[6.0]
  def change
    rename_table :user_finished_lessons, :language_lesson_members
    rename_column :language_lesson_members, :language_module_lesson_id, :lesson_id
    add_column :language_lesson_members, :state, :string
    add_column :language_lesson_members, :language_id, :bigint, references: { to_table: :languages }
    add_column :language_lesson_members, :created_at, :datetime, precision: 6, null: true

    Language::Lesson::Member.update_all("state = 'finished', created_at = inserted_at")
    change_column_null :language_lesson_members, :created_at, false
    remove_column :language_lesson_members, :inserted_at

  end
end
