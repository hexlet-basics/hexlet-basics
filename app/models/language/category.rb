class Language::Category < ApplicationRecord
  validates :name, presence: true

  has_many :languages, dependent: :nullify
end
