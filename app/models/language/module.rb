# frozen_string_literal: true

class Language::Module < ApplicationRecord
  belongs_to :language

  has_many :lessons, dependent: :destroy
  has_many :infos, through: :versions, class_name: 'Language::Module::Version::Info'
end
