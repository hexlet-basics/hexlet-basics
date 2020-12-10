# frozen_string_literal: true

class User < ApplicationRecord
  include UserRepository
  include AASM

  has_secure_password validations: false

  # FIXME AASM has not without_state method. We have to change string comparsion to method that will check state
  validates :email, presence: true, uniqueness: { case_sensitive: false }, 'valid_email_2/email': true, if: -> { state != :removed }

  has_many :lesson_members, class_name: 'Language::Lesson::Member', dependent: :destroy
  has_many :lessons, through: :lesson_members, class_name: 'Language::Lesson'
  has_many :language_members, class_name: 'Language::Member', dependent: :destroy
  has_many :accounts, dependent: :destroy

  aasm :state do
    state :active, initial: true
    state :initail
    state :waiting_confirmation
    state :removed

    event :activate do
      transitions from: %i[created archived], to: :active
    end

    event :mark_as_archived do
      transitions from: %i[created archived], to: :archived
    end
  end

  def guest?
    false
  end

  def valid_password?(password)
    return false if password_digest.nil?

    authenticate(password)
  end
end
