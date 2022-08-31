class AddFinishedLessonsCountToLanguageMembers < ActiveRecord::Migration[7.0]
  def self.up
    add_column :language_members, :finished_lessons_count, :integer, null: false, default: 0
  end

  def self.down
    remove_column :language_members, :finished_lessons_count
  end
end
