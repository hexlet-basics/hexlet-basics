# typed: false
# frozen_string_literal: true

class AddWebauthnIdIndexConcurrently < ActiveRecord::Migration[8.0]
  # Concurrent index creation cannot run inside a transaction and does not lock
  # writes on the table (unlike a plain CREATE INDEX, which blocks INSERT/UPDATE/
  # DELETE for the whole build).
  disable_ddl_transaction!

  def change
    # Partial unique index: many users may have NULL webauthn_id, but a present
    # handle is unique.
    add_index :users, :webauthn_id,
              unique: true,
              where: "webauthn_id IS NOT NULL",
              algorithm: :concurrently,
              if_not_exists: true
  end
end
