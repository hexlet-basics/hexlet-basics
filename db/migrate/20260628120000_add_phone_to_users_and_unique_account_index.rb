# typed: false
# frozen_string_literal: true

class AddPhoneToUsersAndUniqueAccountIndex < ActiveRecord::Migration[8.0]
  def change
    add_column :users, :phone, :string
    add_column :users, :phone_verified_at, :datetime

    # Partial unique index: many users may have NULL phone, but a present phone is unique.
    add_index :users, :phone, unique: true, where: "phone IS NOT NULL"

    # A social identity (provider + uid) must map to a single user.
    # Uniqueness was previously enforced only at the application level (find_or_initialize).
    add_index :user_accounts, %i[provider uid], unique: true
  end
end
