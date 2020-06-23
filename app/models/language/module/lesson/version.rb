# frozen_string_literal: true

class Language::Module::Lesson::Version < ApplicationRecord
  belongs_to :language_version, class_name: 'Language::Version'
  belongs_to :module_version, class_name: 'Language::Module::Version'
  belongs_to :lesson
  belongs_to :upload
end
