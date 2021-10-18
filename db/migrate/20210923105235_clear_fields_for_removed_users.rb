# frozen_string_literal: true

class ClearFieldsForRemovedUsers < ActiveRecord::Migration[6.1]
  def change
    User.removed.where.not(email: nil).update_all(
      first_name: nil,
      last_name: nil,
      nickname: nil,
      password_digest: nil,
      reset_password_token: nil,
      confirmation_token: nil,
      email: nil
    )
  end
end
