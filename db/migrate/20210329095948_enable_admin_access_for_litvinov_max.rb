# frozen_string_literal: true

class EnableAdminAccessForLitvinovMax < ActiveRecord::Migration[6.1]
  def change
    user = User.find_by email: 'litvinovmksm@gmail.com'
    return unless user

    user.admin = true
    user.save!
  end
end
