# frozen_string_literal: true

class Language::Lesson::Version < ApplicationRecord
  belongs_to :language_version, class_name: 'Language::Version'
  belongs_to :lesson
  belongs_to :language
end
