class CreateLessonMembers < ActiveRecord::Migration[6.0]
  def change
    rename_table :user_finished_lessons, :language_lesson_members
    rename_column :language_lesson_members, :language_module_lesson_id, :lesson_id
    add_column :language_lesson_members, :state, :string
    add_column :language_lesson_members, :language_id, :bigint, references: { to_table: :languages }
    Language::Lesson::Member.update_all(state: :finished)
    Language::Lesson::Member.find_each do |member|
      member.language = member.lesson.language
      member.save!
    end

    change_column_null :language_lesson_members, :language_id, false
    # create_table :language_lesson_members do |t|
    #   t.references :language, null: false, foreign_key: true
    #   t.references :language_lesson, null: false, foreign_key: true
    #   t.references :user, null: false, foreign_key: true
    #   t.string :state

    #   t.timestamps
    # end
  end
end
