# frozen_string_literal: true

class Language::Version::Info < ApplicationRecord
  include Language::Version::InfoRepository

  validates :description, presence: true

  belongs_to :language
  belongs_to :language_version, class_name: 'Language::Version'
end
