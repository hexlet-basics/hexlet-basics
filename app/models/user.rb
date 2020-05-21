# frozen_string_literal: true

class User < ApplicationRecord
  has_secure_password

  def guest?
    false
  end
end
