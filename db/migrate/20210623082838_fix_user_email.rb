# frozen_string_literal: true

class FixUserEmail < ActiveRecord::Migration[6.1]
  def change
    user = User.find_by email: 'vika.mozgova@mail.u'
    return unless user

    user.email = 'vika.mozgova@mail.ru'
    user.save!
  end
end
