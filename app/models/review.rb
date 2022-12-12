# frozen_string_literal: true

# == Schema Information
#
# Table name: reviews
#
#  id          :bigint           not null, primary key
#  language_id :bigint           not null
#  user_id     :bigint           not null
#  state       :string
#  body        :text
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  locale      :string
#
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
