# frozen_string_literal: true

module Language::Lesson::Version::InfoRepository
  extend ActiveSupport::Concern

  included do
    scope :left_join_lesson_member_and_user, ->(user) { joins("LEFT OUTER JOIN language_lesson_members ON language_lesson_members.lesson_id = language_lesson_version_infos.language_lesson_id AND language_lesson_members.user_id = #{user.id}") }
    scope :with_locale, ->(locale = I18n.locale) { where(locale: locale) }
  end
end
