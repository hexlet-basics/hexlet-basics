# frozen_string_literal: true

class AddMembersCountToLanguages < ActiveRecord::Migration[7.0]
  def self.up
    add_column :languages, :members_count, :integer, null: false, default: 0
  end

  def self.down
    remove_column :languages, :members_count
  end
end
