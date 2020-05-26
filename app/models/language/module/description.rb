# frozen_string_literal: true

class Language::Module::Description < ApplicationRecord
  belongs_to :module
  belongs_to :language

  def to_s
    name
  end
end
