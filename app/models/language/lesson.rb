# frozen_string_literal: true

class Language::Lesson < ApplicationRecord
  belongs_to :language
  belongs_to :module
  belongs_to :current_version, optional: true, class_name: 'Language::Lesson::Version'

  has_many :versions, dependent: :destroy
  has_many :members, dependent: :destroy

  has_many :infos, through: :versions, class_name: 'Language::Lesson::Version::Info'
  has_many :current_infos, through: :current_version, class_name: 'Language::Lesson::Version::Info', source: :infos
end
