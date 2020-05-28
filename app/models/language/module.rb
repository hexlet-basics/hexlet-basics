# frozen_string_literal: true

class Language::Module < ApplicationRecord
  include Language::ModuleRepository

  belongs_to :language
  has_many :descriptions, dependent: :destroy, class_name: 'Language::Module::Description'
  has_many :lessons, dependent: :destroy

  def directory
    "#{order}-#{slug}"
  end
end
