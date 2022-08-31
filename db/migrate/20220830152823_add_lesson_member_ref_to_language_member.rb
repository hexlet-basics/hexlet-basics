class AddLessonMemberRefToLanguageMember < ActiveRecord::Migration[7.0]
  def change
    add_reference :language_lesson_members, :language_member, null: true, foreign_key: true
  end
end
