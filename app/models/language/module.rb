# frozen_string_literal: true

class Language::Module < ApplicationRecord
  has_many :lessons, dependent: :destroy
  belongs_to :language

  validates :slug, uniqueness: { scope: :language, message: 'slug should be uniqueness on language' },
                   presence: true
end
