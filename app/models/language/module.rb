# frozen_string_literal: true

class Language::Module < ApplicationRecord
  include Language::ModuleRepository

  belongs_to :language
  has_many :descriptions, dependent: :destroy
  has_many :lessons, dependent: :destroy

  def directory
    "#{order}-#{slug}"
  end
end
