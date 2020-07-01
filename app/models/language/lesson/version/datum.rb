# frozen_string_literal: true

class Language::Lesson::Version::Datum < ApplicationRecord
  belongs_to :lesson
  belongs_to :language
  belongs_to :version
  belongs_to :language_version, class_name: 'Language::Version'

  def to_s
    name
  end
end
