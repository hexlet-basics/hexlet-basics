# frozen_string_literal: true

class Language::Module::Lesson::Description < ApplicationRecord
  belongs_to :lesson
  belongs_to :language

  def to_s
    name
  end
end
