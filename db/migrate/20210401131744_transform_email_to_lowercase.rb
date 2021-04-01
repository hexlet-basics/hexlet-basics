# frozen_string_literal: true

class TransformEmailToLowercase < ActiveRecord::Migration[6.1]
  def change
    User.update_all('email = lower(email)')
  end
end
