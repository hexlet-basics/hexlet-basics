class SetNotNull < ActiveRecord::Migration[7.0]
  def change
    change_column_null(:language_lesson_members, :language_member_id, false)
  end
end
