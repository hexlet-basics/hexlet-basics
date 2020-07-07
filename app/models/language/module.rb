# frozen_string_literal: true

class Language::Module < ApplicationRecord
  belongs_to :language
  belongs_to :current_version, optional: true, class_name: 'Language::Module::Version'

  has_many :lessons, dependent: :destroy
  has_many :data, through: :versions, class_name: 'Language::Module::Version::Datum'
  has_many :current_data, through: :current_version, class_name: 'Language::Module::Version::Datum', source: :data
end
