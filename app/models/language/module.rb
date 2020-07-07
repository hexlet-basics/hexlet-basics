# frozen_string_literal: true

class Language::Module < ApplicationRecord
  belongs_to :language
  belongs_to :current_version, optional: true, class_name: 'Language::Module::Version'

  has_many :lessons, dependent: :destroy
  has_many :infos, through: :versions, class_name: 'Language::Module::Version::Info'
  has_many :current_infos, through: :current_version, class_name: 'Language::Module::Version::Info', source: :infos
end
