# frozen_string_literal: true

class Language::Lesson < ApplicationRecord
  belongs_to :language
  belongs_to :module
  belongs_to :version, optional: true, class_name: 'Language::Lesson::Version'
  has_many :descriptions, dependent: :destroy
end
