class Language::Version < ApplicationRecord
  has_many :modules, dependent: :destroy, class_name: 'Language::Module::Version'
  has_many :lessons, dependent: :destroy, class_name: 'Language::Module::Lesson::Version'
end
