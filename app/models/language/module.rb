# frozen_string_literal: true

class Language::Module < ApplicationRecord
  has_many :lessons, dependent: :destroy
  belongs_to :language

  validates :slug, uniqueness: { scope: :language, message: 'slug should be uniqueness on language' },
                   presence: true
  validates :order, presence: true

  def self.get_directory(language_module)
    "#{language_module.order}-#{language_module.slug}"
  end
end
