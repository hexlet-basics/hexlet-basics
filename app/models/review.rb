# frozen_string_literal: true

class Review < ApplicationRecord
  extend Enumerize
  include AASM

  enumerize :locale, in: I18n.available_locales

  validates :locale, presence: true
  validates :body, presence: true

  belongs_to :language
  belongs_to :user

  aasm column: :state do
    state :draft, initial: true
    state :published
    state :archived

    event :publish do
      transitions to: :published
    end

    event :archive do
      transitions to: :archive
    end
  end

  def to_s
    "#{user} (#{locale})"
  end
end
