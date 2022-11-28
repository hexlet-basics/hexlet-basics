# frozen_string_literal: true

class Review < ApplicationRecord
  extend Enumerize

  enumerize :locale, in: I18n.available_locales

  validates :locale, presence: true
  validates :body, presence: true

  belongs_to :language
  belongs_to :user

  def to_s
    "#{user} (#{locale})"
  end
end
