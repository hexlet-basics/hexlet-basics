# frozen_string_literal: true

class User < ApplicationRecord
  include UserRepository

  has_secure_password validations: false

  validates :email, presence: true, uniqueness: { case_sensitive: false }, 'valid_email_2/email': true

  has_many :lesson_members, class_name: 'Language::Lesson::Member', dependent: :destroy
  has_many :lessons, through: :lesson_members, class_name: 'Language::Lesson'
  has_many :language_members, class_name: 'Language::Member', dependent: :destroy
  has_many :accounts, dependent: :destroy

  def guest?
    false
  end

  def valid_password?(password)
    return false if password_digest.nil?

    authenticate(password)
  end
end
