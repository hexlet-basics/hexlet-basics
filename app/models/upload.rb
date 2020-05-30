# frozen_string_literal: true

class Upload < ApplicationRecord
  has_many :lessons, dependent: :destroy, class_name: 'Language::Module::Lesson'
  has_many :modules, dependent: :destroy, class_name: 'Language::Module'
  has_one :language, dependent: :destroy
end
