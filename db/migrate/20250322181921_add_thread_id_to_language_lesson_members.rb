class AddThreadIdToLanguageLessonMembers < ActiveRecord::Migration[8.0]
  def change
    add_column :language_lesson_members, :openai_thread_id, :string
  end
end
