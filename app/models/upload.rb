# frozen_string_literal: true

class Upload < ApplicationRecord
  extend Enumerize

  enumerize :language_name, in: %i[php javascript python java ruby]

  has_one :language_version, dependent: :destroy, class_name: 'Language::Version'
  has_many :module_versions, dependent: :destroy, class_name: 'Language::Module::Version'
  has_many :lesson_versions, dependent: :destroy, class_name: 'Language::Module::Lesson::Version'

end
