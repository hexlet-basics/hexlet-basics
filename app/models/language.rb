# frozen_string_literal: true

class Language < ApplicationRecord
  has_many :modules, dependent: :destroy

  def to_s
    name
  end
end
