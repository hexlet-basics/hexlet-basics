class SetLangMemberFieldNotNull < ActiveRecord::Migration[7.0]
  def up
    change_column :language_lesson_members, :language_member_id, :integer, null: false
  end

  def down
    change_column :language_lesson_members, :language_member_id, :integer, null: true
  end
end
