# frozen_string_literal: true

class Language::Lesson < ApplicationRecord
  belongs_to :language
  belongs_to :module
  belongs_to :version, optional: true, class_name: 'Language::Lesson::Version'

  has_many :versions, dependent: :destroy

  has_many :data, through: :versions, class_name: 'Language::Lesson::Version::Datum'
  has_many :current_data, through: :version, class_name: 'Language::Lesson::Version::Datum', source: :data
end
