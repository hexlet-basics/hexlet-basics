class AddUserToLanguageLessonMemberMessages < ActiveRecord::Migration[8.0]
  def change
    add_reference :language_lesson_member_messages, :user, null: true, foreign_key: true
  end
end
