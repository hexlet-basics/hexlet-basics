# frozen_string_literal: true

class Language < ApplicationRecord
  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy

  def to_s
    name
  end
end
