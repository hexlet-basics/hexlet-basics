# == Schema Information
#
# Table name: language_landing_pages
#
#  id                   :integer          not null, primary key
#  description          :string
#  header               :string
#  listed               :boolean
#  locale               :string
#  main                 :boolean
#  meta_description     :string
#  meta_title           :string
#  order                :string
#  slug                 :string
#  state                :string
#  created_at           :datetime         not null
#  updated_at           :datetime         not null
#  language_category_id :integer
#  language_id          :integer          not null
#
# Indexes
#
#  index_language_landing_pages_on_language_category_id  (language_category_id)
#  index_language_landing_pages_on_language_id           (language_id)
#
# Foreign Keys
#
#  language_category_id  (language_category_id => language_categories.id)
#  language_id           (language_id => languages.id)
#
class Language::LandingPage < ApplicationRecord
  include AASM
  include Language::LandingPageRepository

  belongs_to :language
  belongs_to :language_category, class_name: "Language::Category"
  validates :meta_title, presence: true
  validates :header, presence: true
  validates :slug, presence: true, uniqueness: { scope: :locale }
  validates :main, uniqueness: { scope: [ :locale, :language_id ] }
  # validates :description, presence: true
  validates :locale, presence: true # , inclusion: I18n.available_locales

  def self.ransackable_associations(auth_object = nil)
    [ "language" ]
  end

  def self.ransackable_attributes(_auth_object = nil)
    [ "created_at", "language_slug" ]
  end

  aasm column: :state do
    state :draft, initial: true
    state :archived
    state :published do
      validates :description, presence: true
      validates :language_category, presence: true
    end

    event :archive do
      transitions to: :archived
    end

    event :publish do
      transitions to: :published
    end
  end

  def to_s
    header
  end
end
