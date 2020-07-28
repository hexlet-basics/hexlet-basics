# frozen_string_literal: true

class User < ApplicationRecord
  has_secure_password

  validates :email, presence: true, uniqueness: { case_sensitive: false }, format: { with: URI::MailTo::EMAIL_REGEXP }

  has_many :lesson_members, class_name: 'Language::Lesson::Member', dependent: :destroy

  def guest?
    false
  end

  def finished_members
    lesson_members.finished
  end
end
