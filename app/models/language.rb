# frozen_string_literal: true

class Language < ApplicationRecord
  has_many :modules, dependent: :destroy
  has_many :lessons, class_name: 'Language::Module::Lesson', dependent: :destroy

  def to_s
    name
  end
end
