# frozen_string_literal: true

class BlogPost < ApplicationRecord
  extend Enumerize
  include AASM

  enumerize :locale, in: I18n.available_locales

  validates :name, presence: true
  validates :slug, presence: true, uniqueness: true, format: { with: /\A\w+\z/ }
  validates :locale, presence: true
  validates :body, presence: true

  belongs_to :language, optional: true
  belongs_to :creator, class_name: 'User'

  aasm column: :state do
    state :draft, initial: true
    state :published
    state :archived

    event :publish do
      transitions to: :published
    end

    event :archive do
      transitions from: :running, to: :cleaning
    end
  end

  def to_s
    name
  end
end
