# frozen_string_literal: true

class Language::Module::Lesson::Exercise < ApplicationRecord
  belongs_to :lesson
  belongs_to :language
end
