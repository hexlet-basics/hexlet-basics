# frozen_string_literal: true

class Language::Version::Info < ApplicationRecord
  validates :description, presence: true

  belongs_to :language
  belongs_to :language_version
end
