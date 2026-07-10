# typed: false
# frozen_string_literal: true

class AddPhoneToUsersAndUniqueAccountIndex < ActiveRecord::Migration[8.0]
  def change
    add_column :users, :phone, :string
    add_column :users, :phone_verified_at, :datetime

    # Indexes are created concurrently in the following migration
    # (AddPhoneAndUserAccountIndexesConcurrently) to avoid locking writes on
    # large tables.
  end
end
