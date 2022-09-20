# frozen_string_literal: true

class AddLessonsCountToLanguages < ActiveRecord::Migration[7.0]
  def self.up
    add_column :languages, :lessons_count, :integer, null: false, default: 0
  end

  def self.down
    remove_column :languages, :lessons_count
  end
end
