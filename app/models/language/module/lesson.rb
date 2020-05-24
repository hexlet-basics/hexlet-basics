# frozen_string_literal: true

class Language::Module::Lesson < ApplicationRecord
  belongs_to :language
  belongs_to :module
  has_many :descriptions, dependent: :destroy
end
