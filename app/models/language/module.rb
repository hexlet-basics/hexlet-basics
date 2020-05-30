# frozen_string_literal: true

class Language::Module < ApplicationRecord
  include Language::ModuleRepository

  has_paper_trail

  belongs_to :language
  belongs_to :upload
  has_many :descriptions, dependent: :destroy, class_name: 'Language::Module::Description'
  has_many :lessons, dependent: :destroy
end
