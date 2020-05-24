# frozen_string_literal: true

class Language::Module < ApplicationRecord
  belongs_to :language
  has_many :descriptions, dependent: :destroy
  has_many :lessons, dependent: :destroy

  def directory
    "#{order}-#{slug}"
  end
end
