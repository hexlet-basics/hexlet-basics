# frozen_string_literal: true

class LanguageModuleLessonDescription < ApplicationRecord
  belongs_to :lesson
  belongs_to :language
end
