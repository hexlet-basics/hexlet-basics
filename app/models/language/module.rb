# frozen_string_literal: true

class Language::Module < ApplicationRecord
  belongs_to :language
  belongs_to :version, optional: true, class_name: 'Language::Module::Version'
  has_many :descriptions, dependent: :destroy, class_name: 'Language::Module::Description'
  has_many :lessons, dependent: :destroy
end
