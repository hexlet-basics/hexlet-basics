# frozen_string_literal: true

# == Schema Information
#
# Table name: reviews
#
#  id          :bigint           not null, primary key
#  body        :text
#  first_name  :string
#  last_name   :string
#  locale      :string
#  pinned      :boolean
#  state       :string
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  language_id :bigint           not null
#  user_id     :bigint           not null
#
# Indexes
#
#  index_reviews_on_language_id  (language_id)
#  index_reviews_on_user_id      (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (language_id => languages.id)
#  fk_rails_...  (user_id => users.id)
#
class Review < ApplicationRecord
  extend Enumerize
  include AASM
  include ReviewRepository

  # TODO: switch to enum
  enumerize :locale, in: I18n.available_locales

  def self.ransackable_attributes(_auth_object = nil)
    [ "created_at" ]
  end

  def self.ransackable_associations(_auth_object = nil)
    []
  end

  validates :locale, presence: true
  validates :body, presence: true
  validates :first_name, presence: true

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
      transitions to: :archived
    end
  end

  def to_s
    "#{first_name} #{last_name}"
  end
end
