# frozen_string_literal: true

class Language::Module::Lesson::Member < ApplicationRecord
  belongs_to :user
  belongs_to :lesson
  belongs_to :language
end
