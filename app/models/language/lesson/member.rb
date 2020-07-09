# frozen_string_literal: true

class Language::Lesson::Member < ApplicationRecord
  belongs_to :user
  belongs_to :language
  belongs_to :lesson
  belongs_to :lesson_version, class_name: 'Language::Lesson::Version'
end
