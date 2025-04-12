class AddMessagesCountToLangaugeLessonMember < ActiveRecord::Migration[8.0]
  def change
    add_column :language_lesson_members, :messages_count, :integer, default: 0
  end
end
