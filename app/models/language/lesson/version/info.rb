# frozen_string_literal: true

class Language::Lesson::Version::Info < ApplicationRecord
  include Language::Lesson::Version::InfoRepository

  belongs_to :language
  belongs_to :version
  belongs_to :language_version, class_name: 'Language::Version'

  def to_s
    name
  end
end
