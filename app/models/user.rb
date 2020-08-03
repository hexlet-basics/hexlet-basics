# frozen_string_literal: true

class User < ApplicationRecord
  has_secure_password validations: false

  validates :email, presence: true, uniqueness: { case_sensitive: false }, format: { with: URI::MailTo::EMAIL_REGEXP }

  has_many :lesson_members, class_name: 'Language::Lesson::Member', dependent: :destroy
  has_many :language_members, class_name: 'Language::Member', dependent: :destroy
  has_many :accounts, dependent: :destroy

  def guest?
    false
  end

  def not_finished_lessons_for_language(language)
    language.current_lessons.left_join_lesson_member_and_user(self)
            .merge(Language::Lesson::Member.started_or_nil)
  end

  def finished_lessons_for_language(language)
    language.current_lessons.includes(:members)
            .merge(Language::Lesson::Member.finished)
            .where(language_lesson_members: { user_id: id })
  end

  def valid_password?(password)
    return false if password_digest.nil?

    authenticate(password)
  end

  def complete_language?(language)
    language_members.find_by(language: language)&.finished?
  end
end
