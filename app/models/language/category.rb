class Language::Category < ApplicationRecord
  validates :name, presence: true

  has_many :languages, dependent: :nullify

  def name
    send :"name_#{I18n.locale}"
  end

  def to_s
    name
  end
end
