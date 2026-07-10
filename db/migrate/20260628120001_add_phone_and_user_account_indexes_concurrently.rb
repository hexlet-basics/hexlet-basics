# typed: false
# frozen_string_literal: true

class AddPhoneAndUserAccountIndexesConcurrently < ActiveRecord::Migration[8.0]
  # Concurrent index creation cannot run inside a transaction and does not lock
  # writes on the table (unlike a plain CREATE INDEX, which blocks INSERT/UPDATE/
  # DELETE for the whole build).
  disable_ddl_transaction!

  def change
    # Partial unique index: many users may have NULL phone, but a present phone is unique.
    add_index :users, :phone,
              unique: true,
              where: "phone IS NOT NULL",
              algorithm: :concurrently,
              if_not_exists: true

    # A social identity (provider + uid) must map to a single user.
    # Uniqueness was previously enforced only at the application level
    # (find_or_initialize). If duplicate (provider, uid) rows already exist,
    # this will leave an INVALID index that must be dropped and rebuilt after
    # de-duplicating the data.
    add_index :user_accounts, %i[provider uid],
              unique: true,
              algorithm: :concurrently,
              if_not_exists: true
  end
end
