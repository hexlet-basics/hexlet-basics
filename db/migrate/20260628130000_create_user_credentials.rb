# typed: false
# frozen_string_literal: true

class CreateUserCredentials < ActiveRecord::Migration[8.0]
  def change
    # Stable per-user WebAuthn handle, generated on first passkey registration.
    add_column :users, :webauthn_id, :string
    add_index :users, :webauthn_id, unique: true, where: "webauthn_id IS NOT NULL"

    create_table :user_credentials do |t|
      t.references :user, null: false, foreign_key: true
      t.string :external_id, null: false
      t.string :public_key, null: false
      t.bigint :sign_count, null: false, default: 0
      t.string :nickname

      t.timestamps
    end

    add_index :user_credentials, :external_id, unique: true
  end
end
