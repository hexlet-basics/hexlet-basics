# frozen_string_literal: true

class Language::Module::Description < ApplicationRecord
  belongs_to :language
  belongs_to :module

  def to_s
    name
  end
end
