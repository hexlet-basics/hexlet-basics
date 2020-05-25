# frozen_string_literal: true

class User < ApplicationRecord
  has_secure_password

  validates :email, presence: true, uniqueness: { case_sensitive: false }, format: { with: URI::MailTo::EMAIL_REGEXP }
  validates :password, presence: true, length: { minimum: 6 }

  def guest?
    false
  end

  def directory_for_code
    chunked = id
              .to_s
              .rjust(6, '0')
              .reverse
              .split('')
              .each_slice(3)
              .to_a
              .map(&:join)
    File.join(chunked)
  end
end
