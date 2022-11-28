# frozen_string_literal: true

class Language::Category < ApplicationRecord
  validates :name, presence: true

  has_many :languages, ->(category) { where(category: category) }, dependent: :nullify, inverse_of: :category

  def name
    send :"name_#{I18n.locale}"
  end

  def to_s
    name
  end
end
