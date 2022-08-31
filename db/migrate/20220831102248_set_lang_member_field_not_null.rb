class SetLangMemberFieldNotNull < ActiveRecord::Migration[7.0]
  def change
    change_column :language_lesson_members, :language_member_id, :integer, null: false
  end
end
